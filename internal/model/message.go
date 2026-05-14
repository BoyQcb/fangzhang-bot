package model

import "time"

// Message 消息模型（用于统计）
type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"index" json:"user_id"`
	ChatID    int64     `gorm:"index" json:"chat_id"`
	MessageID int       `json:"message_id"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}
