package model

import "time"

type UserSetting struct {
	UserID      uint      `gorm:"primaryKey;column:user_id" json:"user_id"`
	LLMProvider string    `gorm:"type:varchar(50);column:llm_provider" json:"llm_provider"`
	LLMBaseURL  string    `gorm:"type:text;column:llm_base_url" json:"llm_base_url"`
	LLMAPIKey   string    `gorm:"type:text;column:llm_api_key" json:"llm_api_key"`
	LLMModel    string    `gorm:"type:varchar(255);column:llm_model" json:"llm_model"`
	Temperature float64   `gorm:"default:0.7;column:temperature" json:"temperature"`
	MaxTokens   *int      `gorm:"column:max_tokens" json:"max_tokens,omitempty"`
	CreatedAt   time.Time `gorm:"default:current_timestamp;column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp;column:updated_at" json:"updated_at"`
}

func (s *UserSetting) TableName() string {
	return "user_settings"
}
