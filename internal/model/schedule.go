package model

import "time"

// Schedule 定时任务模型
type Schedule struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ChatID    int64     `gorm:"index" json:"chat_id"`
	Message   string    `gorm:"type:text" json:"message"`
	CronSpec  string    `gorm:"size:255" json:"cron_spec"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Schedule) TableName() string {
	return "schedules"
}
