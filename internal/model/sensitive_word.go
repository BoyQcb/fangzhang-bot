package model

import "time"

// SensitiveWord 敏感词模型
type SensitiveWord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Word      string    `gorm:"size:255;uniqueIndex" json:"word"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (SensitiveWord) TableName() string {
	return "sensitive_words"
}
