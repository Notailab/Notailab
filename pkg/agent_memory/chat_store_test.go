package agent_memory

import (
	"errors"
	"testing"

	"github.com/Notailab/Notailab/internal/model"
	agent_core "github.com/Notailab/go-agent/agent/core"
)

type fakeChatMessageStore struct {
	records             []model.AgentMessage
	created             []model.AgentMessage
	updated             []model.AgentMessage
	deleted             []uint
	deletedConversation []uint
	replaced            []model.AgentMessage
}

func (f *fakeChatMessageStore) ListMessages(conversationID uint, limit int) ([]model.AgentMessage, error) {
	return append([]model.AgentMessage(nil), f.records...), nil
}

func (f *fakeChatMessageStore) CreateMessage(message *model.AgentMessage) error {
	if message == nil {
		return errors.New("nil message")
	}
	message.MessageID = uint(len(f.records) + 1)
	f.created = append(f.created, *message)
	f.records = append(f.records, *message)
	return nil
}

func (f *fakeChatMessageStore) UpdateMessage(message *model.AgentMessage) error {
	if message == nil {
		return errors.New("nil message")
	}
	f.updated = append(f.updated, *message)
	for idx := range f.records {
		if f.records[idx].ConversationID == message.ConversationID && f.records[idx].MessageID == message.MessageID {
			f.records[idx] = *message
			return nil
		}
	}
	return errors.New("message not found")
}

func (f *fakeChatMessageStore) DeleteMessage(conversationID, messageID uint) error {
	f.deleted = append(f.deleted, messageID)
	for idx := range f.records {
		if f.records[idx].ConversationID == conversationID && f.records[idx].MessageID == messageID {
			f.records = append(f.records[:idx], f.records[idx+1:]...)
			return nil
		}
	}
	return errors.New("message not found")
}

func (f *fakeChatMessageStore) DeleteMessagesByConversation(conversationID uint) error {
	f.deletedConversation = append(f.deletedConversation, conversationID)
	f.records = nil
	return nil
}

func (f *fakeChatMessageStore) ReplaceMessages(conversationID uint, messages []model.AgentMessage) error {
	f.replaced = append([]model.AgentMessage(nil), messages...)
	f.records = append([]model.AgentMessage(nil), messages...)
	return nil
}

func newTestChatStore(store *fakeChatMessageStore) *ChatStore {
	return &ChatStore{
		messageStore:   store,
		ConversationID: 7,
	}
}

func TestChatStoreAppendUpdateDeleteUseIncrementalPersistence(t *testing.T) {
	store := &fakeChatMessageStore{}
	chatStore := newTestChatStore(store)

	if err := chatStore.Append(agent_core.ChatMessage{Role: agent_core.RoleUser, Content: "hello"}); err != nil {
		t.Fatalf("append failed: %v", err)
	}
	if len(store.created) != 1 {
		t.Fatalf("expected 1 create call, got %d", len(store.created))
	}
	if len(store.replaced) != 0 {
		t.Fatalf("append should not replace all messages")
	}

	if err := chatStore.Update(0, agent_core.ChatMessage{Role: agent_core.RoleAssistant, Content: "hi"}); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if len(store.updated) != 1 {
		t.Fatalf("expected 1 update call, got %d", len(store.updated))
	}
	if len(store.replaced) != 0 {
		t.Fatalf("update should not replace all messages")
	}

	if err := chatStore.Delete(0); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if len(store.deleted) != 1 {
		t.Fatalf("expected 1 delete call, got %d", len(store.deleted))
	}
	if len(store.replaced) != 0 {
		t.Fatalf("delete should not replace all messages")
	}
	if got := len(chatStore.messages); got != 0 {
		t.Fatalf("expected in-memory messages to be empty, got %d", got)
	}
}

func TestChatStoreClearUsesConversationDelete(t *testing.T) {
	store := &fakeChatMessageStore{}
	chatStore := newTestChatStore(store)

	if err := chatStore.Append(agent_core.ChatMessage{Role: agent_core.RoleUser, Content: "hello"}); err != nil {
		t.Fatalf("append failed: %v", err)
	}
	if err := chatStore.Clear(); err != nil {
		t.Fatalf("clear failed: %v", err)
	}
	if len(store.deletedConversation) != 1 {
		t.Fatalf("expected 1 conversation delete call, got %d", len(store.deletedConversation))
	}
	if got := len(chatStore.messages); got != 0 {
		t.Fatalf("expected in-memory messages to be empty, got %d", got)
	}
}
