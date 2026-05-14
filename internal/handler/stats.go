package handler

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// MessageRecord 消息记录
type MessageRecord struct {
	UserID    int64
	Username  string
	ChatID    int64
	MessageID int
	Timestamp time.Time
}

var (
	messageStore []MessageRecord
	storeMutex   sync.RWMutex
)

// RecordMessage 记录消息（用于统计）
func RecordMessage(userID int64, username string, chatID int64, messageID int) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	messageStore = append(messageStore, MessageRecord{
		UserID:    userID,
		Username:  username,
		ChatID:    chatID,
		MessageID: messageID,
		Timestamp: time.Now(),
	})

	// 只保留最近10000条记录
	if len(messageStore) > 10000 {
		messageStore = messageStore[1000:]
	}
}

// ShowStats 显示统计信息
func ShowStats(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	storeMutex.RLock()
	defer storeMutex.RUnlock()

	// 统计总消息数
	totalMessages := len(messageStore)

	// 统计今日消息数
	today := time.Now().Truncate(24 * time.Hour)
	todayCount := 0
	for _, record := range messageStore {
		if record.Timestamp.After(today) {
			todayCount++
		}
	}

	statsText := fmt.Sprintf(
		"📊 消息统计\n\n"+
			"总消息数: %d\n"+
			"今日消息: %d\n",
		totalMessages,
		todayCount,
	)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   statsText,
	})
}

// TopUsers 显示活跃用户排行
func TopUsers(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	storeMutex.RLock()
	defer storeMutex.RUnlock()

	// 统计每个用户的消息数
	userCounts := make(map[int64]int)
	userNames := make(map[int64]string)

	// 只统计最近7天
	weekAgo := time.Now().Add(-7 * 24 * time.Hour)
	for _, record := range messageStore {
		if record.Timestamp.After(weekAgo) {
			userCounts[record.UserID]++
			if record.Username != "" {
				userNames[record.UserID] = record.Username
			}
		}
	}

	// 排序
	type userCount struct {
		UserID int64
		Count  int
	}

	var sorted []userCount
	for userID, count := range userCounts {
		sorted = append(sorted, userCount{UserID: userID, Count: count})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})

	// 取前10名
	if len(sorted) > 10 {
		sorted = sorted[:10]
	}

	// 生成文本
	var sb strings.Builder
	sb.WriteString("🏆 活跃用户排行 (最近7天)\n\n")

	for i, uc := range sorted {
		username := userNames[uc.UserID]
		if username == "" {
			username = fmt.Sprintf("用户%d", uc.UserID)
		}
		sb.WriteString(fmt.Sprintf("%d. %s - %d 条消息\n", i+1, username, uc.Count))
	}

	if len(sorted) == 0 {
		sb.WriteString("暂无数据")
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   sb.String(),
	})
}
