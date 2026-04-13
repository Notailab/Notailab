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

func (dao *AgentDAO) UpdateMessage(message *model.AgentMessage) error {
	if message == nil {
		return gorm.ErrInvalidData
	}
	return dao.db.Model(&model.AgentMessage{}).
		Where("conversation_id = ? AND message_id = ?", message.ConversationID, message.MessageID).
		Update("message_json", message.MessageJSON).Error
}

func (dao *AgentDAO) DeleteMessage(conversationID, messageID uint) error {
	return dao.db.Where("conversation_id = ? AND message_id = ?", conversationID, messageID).Delete(&model.AgentMessage{}).Error
}

func (dao *AgentDAO) DeleteMessagesByConversation(conversationID uint) error {
	return dao.db.Where("conversation_id = ?", conversationID).Delete(&model.AgentMessage{}).Error
}

func (dao *AgentDAO) ReplaceMessages(conversationID uint, messages []model.AgentMessage) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("conversation_id = ?", conversationID).Delete(&model.AgentMessage{}).Error; err != nil {
			return err
		}
		if len(messages) == 0 {
			return nil
		}
		return tx.Create(&messages).Error
	})
}

func (dao *AgentDAO) ListMessages(conversationID uint, limit int) ([]model.AgentMessage, error) {
	var messages []model.AgentMessage
	query := dao.db.Where("conversation_id = ?", conversationID).Order("message_id asc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
