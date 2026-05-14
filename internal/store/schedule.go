package store

import (
	"github.com/xxx/fangzhang-bot/internal/model"
)

// ListAllSchedules 列出所有定时任务（不分聊天）
func ListAllSchedules() ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := DB.Order("created_at DESC").Find(&schedules).Error
	return schedules, err
}
