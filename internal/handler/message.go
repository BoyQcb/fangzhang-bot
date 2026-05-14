package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// DeleteMessage 删除消息
func DeleteMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	// 检查权限
	if !isAdmin(ctx, b, update) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ 你没有权限执行此操作",
		})
		return
	}

	// 解析命令参数：/delete <message_id>
	args := strings.Fields(update.Message.Text)
	if len(args) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "用法: /delete <消息ID>",
		})
		return
	}

	messageID := parseInt(args[1])
	if messageID == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ 无效的消息ID",
		})
		return
	}

	// 删除消息
	_, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: messageID,
	})

	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("❌ 删除失败: %v", err),
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "✅ 消息已删除",
	})
}

// EditMessage 编辑消息
func EditMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
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

	// 用法: /edit <message_id> <新内容>
	args := strings.Fields(update.Message.Text)
	if len(args) < 3 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "用法: /edit <消息ID> <新内容>",
		})
		return
	}

	messageID := parseInt(args[1])
	if messageID == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ 无效的消息ID",
		})
		return
	}

	newText := strings.Join(args[2:], " ")

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: messageID,
		Text:      newText,
	})

	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("❌ 编辑失败: %v", err),
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "✅ 消息已编辑",
	})
}

// ForwardMessage 转发消息
func ForwardMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	// 用法: /forward <from_chat_id> <message_id>
	args := strings.Fields(update.Message.Text)
	if len(args) < 3 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "用法: /forward <来源聊天ID> <消息ID>",
		})
		return
	}

	fromChatID := args[1]
	messageID := parseInt(args[2])

	if messageID == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ 无效的消息ID",
		})
		return
	}

	_, err := b.ForwardMessage(ctx, &bot.ForwardMessageParams{
		ChatID:     update.Message.Chat.ID,
		FromChatID: fromChatID,
		MessageID:  messageID,
	})

	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("❌ 转发失败: %v", err),
		})
		return
	}
}

// DefaultHandler 默认消息处理器
func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// 处理普通消息，可用于记录统计信息等
	if update.Message != nil {
		// 这里可以添加消息统计逻辑
		fmt.Printf("收到消息: %s\n", update.Message.Text)
	}
}

// 辅助函数
func isAdmin(ctx context.Context, b *bot.Bot, update *models.Update) bool {
	// 简化版本：实际应该查询Telegram API获取聊天成员状态
	// 这里暂时返回true，实际使用时需要完善
	return true
}

func parseInt(s string) int {
	// 简单的字符串转整数
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		} else {
			return 0
		}
	}
	return result
}
