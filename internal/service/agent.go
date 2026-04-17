package service

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/Notailab/Notailab/internal/dao"
	"github.com/Notailab/Notailab/internal/model"
	memory "github.com/Notailab/Notailab/pkg/agent_memory"
	notailab_tools "github.com/Notailab/Notailab/pkg/tools"

	agent "github.com/Notailab/go-agent/agent/agent"
	agent_core "github.com/Notailab/go-agent/agent/core"
	agent_tools "github.com/Notailab/go-agent/agent/tools"
)

type AgentService struct {
	agentDAO        *dao.AgentDAO
	projectDAO      *dao.ProjectDAO
	fileDAO         *dao.FileDAO
	settingDAO      *dao.UserSettingDAO
	agentMu         sync.Mutex
	agentCache      map[string]*list.Element
	agentOrder      *list.List
	agentCacheLimit int
}

type cachedAgentEntry struct {
	key   string
	agent *cachedAgent
}

const defaultAgentCacheLimit = 64

type cachedAgent struct {
	client *agent.ReactAgent
	mu     sync.Mutex
}

type AgentChatRequest struct {
	ProjectID uint   `json:"project_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

type AgentChatResult struct {
	Conversation *model.AgentConversation
	Reply        string
	History      []model.AgentMessage
}

func NewAgentService(agentDAO *dao.AgentDAO, projectDAO *dao.ProjectDAO, fileDAO *dao.FileDAO, settingDAO *dao.UserSettingDAO) *AgentService {
	return &AgentService{
		agentDAO:        agentDAO,
		projectDAO:      projectDAO,
		fileDAO:         fileDAO,
		settingDAO:      settingDAO,
		agentCache:      make(map[string]*list.Element),
		agentOrder:      list.New(),
		agentCacheLimit: defaultAgentCacheLimit,
	}
}

func (s *AgentService) getOrCreateConversation(userID, projectID uint, projectTitle string) (*model.AgentConversation, error) {
	conversation, err := s.agentDAO.GetConversationByUserAndProject(userID, projectID)
	if err == nil {
		return conversation, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	conversation = &model.AgentConversation{
		UserID:        userID,
		ProjectID:     projectID,
		Title:         projectTitle,
		LastMessageAt: time.Now(),
	}
	if err := s.agentDAO.CreateConversation(conversation); err != nil {
		return nil, err
	}
	return conversation, nil
}

func (s *AgentService) agentCacheKey(userID, projectID, conversationID uint) string {
	return fmt.Sprintf("%d:%d:%d", userID, projectID, conversationID)
}

func (s *AgentService) getCachedAgentLocked(cacheKey string) *cachedAgent {
	if s.agentCache == nil || s.agentOrder == nil {
		return nil
	}
	if element := s.agentCache[cacheKey]; element != nil {
		s.agentOrder.MoveToFront(element)
		if entry, ok := element.Value.(*cachedAgentEntry); ok && entry != nil {
			return entry.agent
		}
	}
	return nil
}

func (s *AgentService) putCachedAgentLocked(cacheKey string, agent *cachedAgent) {
	if s.agentCache == nil {
		s.agentCache = make(map[string]*list.Element)
	}
	if s.agentOrder == nil {
		s.agentOrder = list.New()
	}
	if s.agentCacheLimit <= 0 {
		s.agentCacheLimit = defaultAgentCacheLimit
	}

	if element := s.agentCache[cacheKey]; element != nil {
		element.Value = &cachedAgentEntry{key: cacheKey, agent: agent}
		s.agentOrder.MoveToFront(element)
		return
	}

	element := s.agentOrder.PushFront(&cachedAgentEntry{key: cacheKey, agent: agent})
	s.agentCache[cacheKey] = element

	for s.agentOrder.Len() > s.agentCacheLimit {
		oldest := s.agentOrder.Back()
		if oldest == nil {
			break
		}
		entry, _ := oldest.Value.(*cachedAgentEntry)
		s.agentOrder.Remove(oldest)
		if entry != nil {
			delete(s.agentCache, entry.key)
		}
	}
}

func (s *AgentService) buildAgent(userID, conversationID uint, project *model.Project) (*cachedAgent, error) {
	setting, err := s.settingDAO.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	baseurl := setting.LLMBaseURL
	apikey := setting.LLMAPIKey
	llmModel := setting.LLMModel
	temperature := setting.Temperature
	if temperature <= 0 {
		temperature = 0.7
	}

	cacheKey := s.agentCacheKey(userID, project.ProjectID, conversationID)

	s.agentMu.Lock()
	defer s.agentMu.Unlock()
	if cached := s.getCachedAgentLocked(cacheKey); cached != nil {
		return cached, nil
	}

	chatStore, err := memory.NewChatStore(s.agentDAO, conversationID)
	if err != nil {
		return nil, err
	}
	projectLongStore, err := memory.NewProjectLongStore(s.fileDAO, project.ProjectID)
	if err != nil {
		return nil, err
	}
	memory := agent_core.NewMemory(chatStore, projectLongStore)
	projectTools, err := notailab_tools.NewFileTools(s.fileDAO, project.ProjectID)
	if err != nil {
		return nil, err
	}

	client := agent.NewReactAgent(
		agent.WithLLM(baseurl, llmModel, apikey),
		agent.WithStaticSystemPrompt(buildNotailabSystemPrompt(project)),
		agent.WithMemory(memory),
		agent.WithTools(append(projectTools, agent_tools.NewLongMemoryTool(memory))...),
		agent.WithReporter(agent.NoopReporter{}),
		agent.WithTemperature(temperature),
	)

	if client == nil {
		return nil, fmt.Errorf("failed to create agent client")
	}

	entry := &cachedAgent{client: client}
	s.putCachedAgentLocked(cacheKey, entry)
	return entry, nil
}

func (s *AgentService) Chat(ctx context.Context, userID uint, req AgentChatRequest) (string, error) {
	project, err := s.projectDAO.GetProjectByIDAndUserID(req.ProjectID, userID)
	if err != nil {
		return "", err
	}

	// TODO: consider multiple conversations per project in the future
	conversation, err := s.getOrCreateConversation(userID, req.ProjectID, project.Title)
	if err != nil {
		return "", err
	}

	agentEntry, err := s.buildAgent(userID, conversation.ConversationID, project)
	if err != nil {
		return "", err
	}

	agentEntry.mu.Lock()
	defer agentEntry.mu.Unlock()

	return agentEntry.client.Run(ctx, req.Content)
}
