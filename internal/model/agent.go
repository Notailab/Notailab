package model

import "time"

type AgentConversation struct {
	ConversationID uint      `gorm:"primaryKey;column:conversation_id" json:"conversation_id"`
	UserID         uint      `gorm:"not null;column:user_id;index:idx_agent_conversation_user_project,unique" json:"user_id"`
	ProjectID      uint      `gorm:"not null;column:project_id;index:idx_agent_conversation_user_project,unique" json:"project_id"`
	Title          string    `gorm:"type:varchar(255);not null;column:title" json:"title"`
	Summary        string    `gorm:"type:text;column:summary" json:"summary"`
	LastMessageAt  time.Time `gorm:"column:last_message_at" json:"last_message_at"`
	CreatedAt      time.Time `gorm:"default:current_timestamp;column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:current_timestamp;column:updated_at" json:"updated_at"`
}

func (c *AgentConversation) TableName() string {
	return "agent_conversations"
}

type AgentMessage struct {
	MessageID      uint      `gorm:"primaryKey;column:message_id" json:"message_id"`
	ConversationID uint      `gorm:"not null;column:conversation_id;index" json:"conversation_id"`
	Role           string    `gorm:"type:varchar(20);not null;column:role" json:"role"`
	Content        string    `gorm:"type:text;not null;column:content" json:"content"`
	CreatedAt      time.Time `gorm:"default:current_timestamp;column:created_at" json:"created_at"`
}

func (m *AgentMessage) TableName() string {
	return "agent_messages"
}
