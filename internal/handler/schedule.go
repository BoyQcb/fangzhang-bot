package handler

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/robfig/cron/v3"
)

var (
	cronScheduler *cron.Cron
	schedules     []ScheduleRecord
	scheduleMutex sync.RWMutex
)

// ScheduleRecord 定时任务记录
type ScheduleRecord struct {
	ID        int
	ChatID    int64
	Message   string
	CronSpec  string
	IsActive  bool
	EntryID   cron.EntryID
}

// InitScheduler 初始化定时任务调度器
func InitScheduler() {
	cronScheduler = cron.New(cron.WithSeconds())
	go cronScheduler.Start()
}

// AddSchedule 添加定时任务
func AddSchedule(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	if !isAdmin(ctx, b, update) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ 你没有权限执行此操作",
		})
		return
	}

	// 用法: /addschedule <cron表达式> <消息内容>
	// 例如: /addschedule "0 9 * * *" "早上好！"
	args := strings.Fields(update.Message.Text)
	if len(args) < 3 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "用法: /addschedule <cron表达式> <消息内容>\n例如: /addschedule \"0 9 * * *\" \"早上好！\"",
		})
		return
	}

	cronSpec := args[1]
	message := strings.Join(args[2:], " ")

	// 添加定时任务
	entryID, err := cronScheduler.AddFunc(cronSpec, func() {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   message,
		})
	})

	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("❌ 添加定时任务失败: %v", err),
		})
		return
	}

	// 记录任务
	scheduleMutex.Lock()
	schedules = append(schedules, ScheduleRecord{
		ID:       len(schedules) + 1,
		ChatID:   update.Message.Chat.ID,
		Message:  message,
		CronSpec: cronSpec,
		IsActive: true,
		EntryID:  entryID,
	})
	scheduleMutex.Unlock()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("✅ 定时任务已添加\nCron: %s\n消息: %s", cronSpec, message),
	})
}

// DeleteSchedule 删除定时任务
func DeleteSchedule(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	if !isAdmin(ctx, b, update) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ 你没有权限执行此操作",
		})
		return
	}

	args := strings.Fields(update.Message.Text)
	if len(args) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "用法: /delschedule <任务ID>",
		})
		return
	}

	scheduleID := parseInt(args[1])

	scheduleMutex.Lock()
	defer scheduleMutex.Unlock()

	found := false
	for i, sched := range schedules {
		if sched.ID == scheduleID {
			// 移除定时任务
			cronScheduler.Remove(sched.EntryID)
			// 从切片中删除
			schedules = append(schedules[:i], schedules[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ 未找到该任务",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("✅ 定时任务 %d 已删除", scheduleID),
	})
}

// ListSchedules 列出所有定时任务
func ListSchedules(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	scheduleMutex.RLock()
	defer scheduleMutex.RUnlock()

	if len(schedules) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "当前没有定时任务",
		})
		return
	}

	var sb strings.Builder
	sb.WriteString("📋 定时任务列表:\n\n")

	for _, sched := range schedules {
		if sched.IsActive {
			sb.WriteString(fmt.Sprintf("ID: %d\n", sched.ID))
			sb.WriteString(fmt.Sprintf("Cron: %s\n", sched.CronSpec))
			sb.WriteString(fmt.Sprintf("消息: %s\n", sched.Message))
			sb.WriteString(fmt.Sprintf("下次执行: %s\n\n", getNextRunTime(sched.CronSpec)))
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   sb.String(),
	})
}

// 辅助函数
func getNextRunTime(cronSpec string) string {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronSpec)
	if err != nil {
		return "未知"
	}

	next := schedule.Next(time.Now())
	return next.Format("2006-01-02 15:04:05")
}
