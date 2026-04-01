package dao

import (
	"gorm.io/gorm"

	"github.com/Notailab/Notailab/internal/model"
)

type AgentDAO struct {
	db *gorm.DB
}

func NewAgentDAO(db *gorm.DB) *AgentDAO {
	return &AgentDAO{db: db}
}

func (dao *AgentDAO) CreateConversation(conversation *model.AgentConversation) error {
	return dao.db.Create(conversation).Error
}

func (dao *AgentDAO) GetConversationByUserAndProject(userID, projectID uint) (*model.AgentConversation, error) {
	var conversation model.AgentConversation
	if err := dao.db.Where("user_id = ? AND project_id = ?", userID, projectID).First(&conversation).Error; err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (dao *AgentDAO) UpdateConversation(conversation *model.AgentConversation) error {
	return dao.db.Save(conversation).Error
}

func (dao *AgentDAO) CreateMessage(message *model.AgentMessage) error {
	return dao.db.Create(message).Error
}

func (dao *AgentDAO) ListMessages(conversationID uint, limit int) ([]model.AgentMessage, error) {
	var messages []model.AgentMessage
	query := dao.db.Where("conversation_id = ?", conversationID).Order("message_id desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	for left, right := 0, len(messages)-1; left < right; left, right = left+1, right-1 {
		messages[left], messages[right] = messages[right], messages[left]
	}

	return messages, nil
}
