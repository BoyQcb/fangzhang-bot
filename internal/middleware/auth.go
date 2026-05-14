package middleware

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Auth 权限验证中间件
func Auth(superUsers []int64) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			// 如果没有消息，跳过验证
			if update.Message == nil {
				next(ctx, b, update)
				return
			}

			// 检查用户是否在超级管理员列表中
			userID := update.Message.From.ID
			isSuperUser := false
			for _, id := range superUsers {
				if id == userID {
					isSuperUser = true
					break
				}
			}

			// 如果是超级管理员，直接通过
			if isSuperUser {
				next(ctx, b, update)
				return
			}

			// 普通用户只能使用部分命令
			// 这里可以根据需要设置更细粒度的权限控制
			command := update.Message.Text
			if strings.HasPrefix(command, "/start") ||
				strings.HasPrefix(command, "/help") ||
				strings.HasPrefix(command, "/lottery") ||
				strings.HasPrefix(command, "/join") {
				next(ctx, b, update)
				return
			}

			// 其他命令需要管理员权限
			// 这里应该查询数据库或使用Telegram API检查用户是否为管理员
			// 简化处理：直接拒绝
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ 你没有权限执行此操作",
			})
		}
	}
}

// isAdmin 检查用户是否为管理员（简化版本）
func isAdmin(ctx context.Context, b *bot.Bot, userID int64, chatID int64) bool {
	// 实际应该调用 Telegram API: GetChatMember
	// 这里返回 true 作为示例
	return true
}
