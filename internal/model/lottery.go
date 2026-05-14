package model

import "time"

// Lottery 抽奖模型
type Lottery struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ChatID      int64     `gorm:"index" json:"chat_id"`
	MessageID   int       `json:"message_id"`
	Prize       string    `gorm:"size:255" json:"prize"`
	WinnerID    int64     `json:"winner_id,omitempty"`
	ParticipantCount int   `json:"participant_count"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedBy   int64     `json:"created_by"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Lottery) TableName() string {
	return "lotteries"
}

// Participant 参与者模型
type Participant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	LotteryID uint      `gorm:"index" json:"lottery_id"`
	UserID    int64     `gorm:"index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (Participant) TableName() string {
	return "participants"
}
