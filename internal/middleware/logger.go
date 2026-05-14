package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Logger 日志中间件
func Logger(ctx context.Context, b *bot.Bot, update *models.Update) {
	// 记录消息
	if update.Message != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		username := update.Message.From.FirstName
		if update.Message.From.Username != "" {
			username = "@" + update.Message.From.Username
		}

		fmt.Printf("[%s] %s (ID: %d): %s\n",
			timestamp,
			username,
			update.Message.From.ID,
			update.Message.Text,
		)

		// 记录消息到统计（可选）
		// 这里可以调用 handler.RecordMessage()
	}
}
