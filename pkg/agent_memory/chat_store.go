package agent_memory

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Notailab/Notailab/internal/model"
	agent_core "github.com/Notailab/go-agent/agent/core"
	agent_storage "github.com/Notailab/go-agent/agent/storage"
)

type chatMessageStore interface {
	ListMessages(conversationID uint, limit int) ([]model.AgentMessage, error)
	CreateMessage(message *model.AgentMessage) error
	UpdateMessage(message *model.AgentMessage) error
	DeleteMessage(conversationID, messageID uint) error
	DeleteMessagesByConversation(conversationID uint) error
	ReplaceMessages(conversationID uint, messages []model.AgentMessage) error
}

type storedChatMessage struct {
	recordID uint
	message  agent_core.ChatMessage
}

type ChatStore struct {
	messageStore   chatMessageStore
	ConversationID uint
	mu             sync.RWMutex
	messages       []storedChatMessage
}

func NewChatStore(messageStore chatMessageStore, conversationID uint) (*ChatStore, error) {
	store := &ChatStore{messageStore: messageStore, ConversationID: conversationID}
	if messageStore == nil {
		return nil, fmt.Errorf("message store is nil")
	}
	if err := store.Load(); err != nil {
		return nil, err
	}
	return store, nil
}

func DecodeMessage(record model.AgentMessage) (agent_core.ChatMessage, error) {
	if len(record.MessageJSON) == 0 {
		return agent_core.ChatMessage{}, fmt.Errorf("message json is empty")
	}

	var message agent_core.ChatMessage
	if err := json.Unmarshal(record.MessageJSON, &message); err != nil {
		return agent_core.ChatMessage{}, fmt.Errorf("decode message json: %w", err)
	}
	return message, nil
}

func EncodeMessage(conversationID uint, message agent_core.ChatMessage) model.AgentMessage {
	payload, err := json.Marshal(message)
	if err != nil {
		payload = []byte(`{}`)
	}

	return model.AgentMessage{
		ConversationID: conversationID,
		MessageJSON:    payload,
	}
}

func (s *ChatStore) Load() error {
	if s == nil || s.messageStore == nil {
		return fmt.Errorf("store is nil")
	}
	records, err := s.messageStore.ListMessages(s.ConversationID, 0)
	if err != nil {
		return err
	}
	messages := make([]storedChatMessage, 0, len(records))
	for _, record := range records {
		message, err := DecodeMessage(record)
		if err != nil {
			return err
		}
		messages = append(messages, storedChatMessage{recordID: record.MessageID, message: message})
	}
	s.mu.Lock()
	s.messages = messages
	s.mu.Unlock()
	return nil
}

func (s *ChatStore) snapshotLocked() []agent_core.ChatMessage {
	result := make([]agent_core.ChatMessage, 0, len(s.messages))
	for _, item := range s.messages {
		result = append(result, item.message)
	}
	return result
}

func (s *ChatStore) reloadLocked() error {
	records, err := s.messageStore.ListMessages(s.ConversationID, 0)
	if err != nil {
		return err
	}
	messages := make([]storedChatMessage, 0, len(records))
	for _, record := range records {
		message, err := DecodeMessage(record)
		if err != nil {
			return err
		}
		messages = append(messages, storedChatMessage{recordID: record.MessageID, message: message})
	}
	s.messages = messages
	return nil
}

func (s *ChatStore) persistAllLocked(messages []agent_core.ChatMessage) error {
	records := make([]model.AgentMessage, 0, len(messages))
	for _, item := range messages {
		records = append(records, EncodeMessage(s.ConversationID, item))
	}
	if err := s.messageStore.ReplaceMessages(s.ConversationID, records); err != nil {
		return err
	}
	return s.reloadLocked()
}

func (s *ChatStore) Get(index int) (agent_core.ChatMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	messages := s.snapshotLocked()
	if index < 0 || index >= len(messages) {
		return agent_core.ChatMessage{}, fmt.Errorf("index out of bounds")
	}
	return messages[index], nil
}

func (s *ChatStore) Append(message agent_core.ChatMessage) error {
	if s == nil || s.messageStore == nil {
		return fmt.Errorf("store is nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	record := EncodeMessage(s.ConversationID, message)
	if err := s.messageStore.CreateMessage(&record); err != nil {
		return err
	}
	s.messages = append(s.messages, storedChatMessage{recordID: record.MessageID, message: message})
	return nil
}

func (s *ChatStore) Update(index int, message agent_core.ChatMessage) error {
	if s == nil || s.messageStore == nil {
		return fmt.Errorf("store is nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.messages) {
		return fmt.Errorf("index out of bounds")
	}
	record := EncodeMessage(s.ConversationID, message)
	record.MessageID = s.messages[index].recordID
	if err := s.messageStore.UpdateMessage(&record); err != nil {
		return err
	}
	s.messages[index] = storedChatMessage{recordID: record.MessageID, message: message}
	return nil
}

func (s *ChatStore) Replace(start, end int, messages []agent_core.ChatMessage) error {
	if s == nil || s.messageStore == nil {
		return fmt.Errorf("store is nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current := s.snapshotLocked()
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}
	if start > len(current) {
		start = len(current)
	}
	if end > len(current) {
		end = len(current)
	}
	updated := make([]agent_core.ChatMessage, 0, len(current)-(end-start)+len(messages))
	updated = append(updated, current[:start]...)
	updated = append(updated, messages...)
	updated = append(updated, current[end:]...)
	if err := s.persistAllLocked(updated); err != nil {
		return err
	}
	return nil
}

func (s *ChatStore) Delete(index int) error {
	if s == nil || s.messageStore == nil {
		return fmt.Errorf("store is nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.messages) {
		return fmt.Errorf("index out of bounds")
	}
	entry := s.messages[index]
	if err := s.messageStore.DeleteMessage(s.ConversationID, entry.recordID); err != nil {
		return err
	}
	s.messages = append(s.messages[:index], s.messages[index+1:]...)
	return nil
}

func (s *ChatStore) List() ([]agent_core.ChatMessage, error) {
	if s == nil || s.messageStore == nil {
		return []agent_core.ChatMessage{}, nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.snapshotLocked(), nil
}

func (s *ChatStore) Count() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.messages), nil
}

func (s *ChatStore) Clear() error {
	if s == nil || s.messageStore == nil {
		return fmt.Errorf("store is nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.messageStore.DeleteMessagesByConversation(s.ConversationID); err != nil {
		return err
	}
	s.messages = nil
	return nil
}

func (s *ChatStore) Clone() agent_core.ChatMemoryStore {
	s.mu.RLock()
	messages := s.snapshotLocked()
	s.mu.RUnlock()
	store := agent_storage.NewInMemoryChatStore()
	for _, message := range messages {
		_ = store.Append(message)
	}
	return store
}
