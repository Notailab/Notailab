package service

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Notailab/Notailab/internal/dao"
	"github.com/Notailab/Notailab/internal/model"
	agentcore "github.com/Notailab/go-agent/agent/core"
	"gorm.io/gorm"
)

type AgentService struct {
	agentDAO   *dao.AgentDAO
	projectDAO *dao.ProjectDAO
	fileDAO    *dao.FileDAO
	settingDAO *dao.UserSettingDAO
	config     *agentcore.Config
}

type AgentChatRequest struct {
	ProjectID     uint   `json:"project_id" binding:"required"`
	FileID        *uint  `json:"file_id"`
	Message       string `json:"message" binding:"required"`
	EditorContent string `json:"editor_content"`
	SelectedText  string `json:"selected_text"`
}

type AgentChatResult struct {
	Conversation *model.AgentConversation
	Reply        string
	History      []model.AgentMessage
}

func NewAgentService(agentDAO *dao.AgentDAO, projectDAO *dao.ProjectDAO, fileDAO *dao.FileDAO, settingDAO *dao.UserSettingDAO) *AgentService {
	return &AgentService{
		agentDAO:   agentDAO,
		projectDAO: projectDAO,
		fileDAO:    fileDAO,
		settingDAO: settingDAO,
		config:     agentcore.NewConfig(),
	}
}

func (s *AgentService) resolveLLMConfig(userID uint) (string, string, string, float64, *int, error) {
	baseURL := os.Getenv("NOTAILAB_LLM_BASE_URL")
	apiKey := os.Getenv("NOTAILAB_LLM_API_KEY")
	llmModel := os.Getenv("NOTAILAB_LLM_MODEL")
	temperature := 0.7
	var maxTokens *int

	if s.settingDAO != nil {
		setting, err := s.settingDAO.GetByUserID(userID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", "", 0, nil, err
		}
		if err == nil && setting != nil {
			if setting.LLMBaseURL != "" {
				baseURL = setting.LLMBaseURL
			}
			if setting.LLMAPIKey != "" {
				apiKey = setting.LLMAPIKey
			}
			if setting.LLMModel != "" {
				llmModel = setting.LLMModel
			}
			if setting.Temperature > 0 {
				temperature = setting.Temperature
			}
			if setting.MaxTokens != nil {
				maxTokens = setting.MaxTokens
			}
		}
	}

	if baseURL == "" || apiKey == "" || llmModel == "" {
		return "", "", "", 0, nil, errors.New("缺少大模型配置，请在用户设置或环境变量中配置 NOTAILAB_LLM_BASE_URL、NOTAILAB_LLM_API_KEY 和 NOTAILAB_LLM_MODEL")
	}

	return baseURL, apiKey, llmModel, temperature, maxTokens, nil
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

func (s *AgentService) buildSystemPrompt(project *model.Project, file *model.File) string {
	var builder strings.Builder
	builder.WriteString("你是 Notailab 的项目学习记录助手。你的任务是围绕项目给出可执行建议、总结、规划和文档优化建议。请始终使用中文，回答要具体、简洁、可落地。\n\n")
	builder.WriteString(fmt.Sprintf("项目名称：%s\n", project.Title))
	if project.Description != "" {
		builder.WriteString(fmt.Sprintf("项目描述：%s\n", project.Description))
	}
	if project.Status != "" {
		builder.WriteString(fmt.Sprintf("项目状态：%s\n", project.Status))
	}
	builder.WriteString(fmt.Sprintf("项目进度：%d%%\n", project.Progress))

	if file != nil {
		builder.WriteString(fmt.Sprintf("当前文件：%s\n", file.Name))
		if file.Content != "" {
			builder.WriteString("当前文件内容：\n")
			builder.WriteString(truncateText(file.Content, 6000))
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func truncateText(text string, max int) string {
	if max <= 0 || len(text) <= max {
		return text
	}
	return text[:max] + "\n..."
}

func (s *AgentService) Chat(userID uint, req AgentChatRequest) (*AgentChatResult, error) {
	baseURL, apiKey, llmModel, temperature, maxTokens, err := s.resolveLLMConfig(userID)
	
	if err != nil {
		return nil, err
	}
	
	project, err := s.projectDAO.GetProjectByIDAndUserID(req.ProjectID, userID)
	if err != nil {
		return nil, err
	}

	var currentFile *model.File
	if req.FileID != nil {
		currentFile, err = s.fileDAO.GetFileByIDAndProjectID(*req.FileID, req.ProjectID)
		if err != nil {
			return nil, err
		}
	}
	
	conversation, err := s.getOrCreateConversation(userID, req.ProjectID, project.Title)
	if err != nil {
		return nil, err
	}
	
	historyModels, err := s.agentDAO.ListMessages(conversation.ConversationID, 20)
	if err != nil {
		return nil, err
	}
	
	llm := agentcore.NewLLM(baseURL, apiKey, llmModel, temperature, maxTokens)
	agent := agentcore.NewAgent(
		project.Title,
		llm,
		s.buildSystemPrompt(project, currentFile),
		s.config,
	)

	for _, item := range historyModels {
		agent.AddMessage(agentcore.Message{Role: item.Role, Content: item.Content})
	}

	userContent := req.Message
	if req.SelectedText != "" {
		userContent += "\n\n当前选中文本：\n" + req.SelectedText
	}
	if req.EditorContent != "" {
		userContent += "\n\n当前编辑内容：\n" + truncateText(req.EditorContent, 6000)
	}
	
	reply, err := agent.Run(userContent, nil)
	if err != nil {
		return nil, err
	}

	if err := s.agentDAO.CreateMessage(&model.AgentMessage{
		ConversationID: conversation.ConversationID,
		Role:           "user",
		Content:        userContent,
	}); err != nil {
		return nil, err
	}

	if err := s.agentDAO.CreateMessage(&model.AgentMessage{
		ConversationID: conversation.ConversationID,
		Role:           "assistant",
		Content:        reply,
	}); err != nil {
		return nil, err
	}

	conversation.LastMessageAt = time.Now()
	if err := s.agentDAO.UpdateConversation(conversation); err != nil {
		return nil, err
	}

	historyModels, err = s.agentDAO.ListMessages(conversation.ConversationID, 20)
	if err != nil {
		return nil, err
	}

	return &AgentChatResult{
		Conversation: conversation,
		Reply:        reply,
		History:      historyModels,
	}, nil
}
