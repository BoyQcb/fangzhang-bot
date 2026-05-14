package service

import (
	"context"
	"fmt"
	"time"

	"github.com/xxx/fangzhang-bot/internal/model"
	"github.com/xxx/fangzhang-bot/internal/store"
	"gorm.io/gorm"
)

// SchedulerService 定时任务服务
type SchedulerService struct {
	db *gorm.DB
}

// NewSchedulerService 创建定时任务服务
func NewSchedulerService(db *gorm.DB) *SchedulerService {
	return &SchedulerService{db: db}
}

// CreateSchedule 创建定时任务
func (s *SchedulerService) CreateSchedule(chatID int64, message, cronSpec string, createdBy int64) (*model.Schedule, error) {
	schedule := &model.Schedule{
		ChatID:    chatID,
		Message:   message,
		CronSpec:  cronSpec,
		IsActive:  true,
		CreatedBy: createdBy,
	}

	err := store.CreateSchedule(schedule)
	if err != nil {
		return nil, fmt.Errorf("创建定时任务失败: %w", err)
	}

	return schedule, nil
}

// GetSchedule 获取定时任务
func (s *SchedulerService) GetSchedule(id uint) (*model.Schedule, error) {
	return store.GetScheduleByID(id)
}

// UpdateSchedule 更新定时任务
func (s *SchedulerService) UpdateSchedule(schedule *model.Schedule) error {
	return store.UpdateSchedule(schedule)
}

// DeleteSchedule 删除定时任务
func (s *SchedulerService) DeleteSchedule(id uint) error {
	return store.DeleteSchedule(id)
}

// ListSchedules 列出定时任务
func (s *SchedulerService) ListSchedules(chatID int64, offset, limit int) ([]model.Schedule, error) {
	return store.ListSchedules(chatID, offset, limit)
}

// ListActiveSchedules 列出活跃的定时任务
func (s *SchedulerService) ListActiveSchedules() ([]model.Schedule, error) {
	return store.ListActiveSchedules()
}

// ExecuteSchedule 执行定时任务
func (s *SchedulerService) ExecuteSchedule(schedule *model.Schedule) error {
	// 这里应该调用 Telegram Bot API 发送消息
	// 实际实现需要在 bot 上下文中调用
	fmt.Printf("执行定时任务: %s - %s\n", schedule.CronSpec, schedule.Message)
	return nil
}

// CalculateNextRun 计算下次执行时间
func (s *SchedulerService) CalculateNextRun(cronSpec string) (time.Time, error) {
	// 使用 robfig/cron 库解析 cron 表达式
	// 这里返回当前时间作为占位符
	return time.Now().Add(1 * time.Hour), nil
}
