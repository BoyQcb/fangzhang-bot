package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// KickUser 踢人
func KickUser(ctx context.Context, b *bot.Bot, update *models.Update) {
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

	// 用法: /kick <user_id> 或回复某条消息
	var userID int64
	if update.Message.ReplyToMessage != nil {
		userID = update.Message.ReplyToMessage.From.ID
	} else {
		args := strings.Fields(update.Message.Text)
		if len(args) < 2 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "用法: /kick <用户ID> 或回复某条消息",
			})
			return
		}
		userID = int64(parseInt(args[1]))
	}

	if userID == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ 无效的用户ID",
		})
		return
	}

	// 踢出用户（Ban）
	_, err := b.BanChatMember(ctx, &bot.BanChatMemberParams{
		ChatID: update.Message.Chat.ID,
		UserID: userID,
	})

	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("❌ 踢人失败: %v", err),
		})
		return
	}

	// 立即解封，这样用户可以重新加入
	_, _ = b.UnbanChatMember(ctx, &bot.UnbanChatMemberParams{
		ChatID: update.Message.Chat.ID,
		UserID: userID,
	})

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("✅ 用户 %d 已被踢出", userID),
	})
}

// MuteUser 禁言
func MuteUser(ctx context.Context, b *bot.Bot, update *models.Update) {
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

	// 用法: /mute <user_id> <duration_seconds> 或回复消息
	var userID int64
	var duration int = 60 // 默认禁言60秒

	args := strings.Fields(update.Message.Text)
	if update.Message.ReplyToMessage != nil {
		userID = update.Message.ReplyToMessage.From.ID
		if len(args) >= 2 {
			duration = parseInt(args[1])
		}
	} else {
		if len(args) < 2 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "用法: /mute <用户ID> [时长(秒)] 或回复某条消息",
			})
			return
		}
		userID = int64(parseInt(args[1]))
		if len(args) >= 3 {
			duration = parseInt(args[2])
		}
	}

	if userID == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ 无效的用户ID",
		})
		return
	}

	// 计算禁言结束时间
	untilDate := time.Now().Add(time.Duration(duration) * time.Second).Unix()

	// 限制用户权限（禁言）
	_, err := b.RestrictChatMember(ctx, &bot.RestrictChatMemberParams{
		ChatID: update.Message.Chat.ID,
		UserID: userID,
		Permissions: &models.ChatPermissions{
			CanSendMessages:      false,
			CanSendPolls:         false,
			CanSendOtherMessages: false,
			CanAddWebPagePreviews: false,
			CanChangeInfo:        false,
			CanInviteUsers:       false,
			CanPinMessages:       false,
		},
		UntilDate: untilDate,
	})

	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("❌ 禁言失败: %v", err),
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("✅ 用户 %d 已被禁言 %d 秒", userID, duration),
	})
}

// UnmuteUser 解除禁言
func UnmuteUser(ctx context.Context, b *bot.Bot, update *models.Update) {
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

	var userID int64
	if update.Message.ReplyToMessage != nil {
		userID = update.Message.ReplyToMessage.From.ID
	} else {
		args := strings.Fields(update.Message.Text)
		if len(args) < 2 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "用法: /unmute <用户ID> 或回复某条消息",
			})
			return
		}
		userID = int64(parseInt(args[1]))
	}

	if userID == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ 无效的用户ID",
		})
		return
	}

	// 恢复用户权限
	_, err := b.RestrictChatMember(ctx, &bot.RestrictChatMemberParams{
		ChatID: update.Message.Chat.ID,
		UserID: userID,
		Permissions: &models.ChatPermissions{
			CanSendMessages:       true,
			CanSendPolls:          true,
			CanSendOtherMessages: true,
			CanAddWebPagePreviews: true,
			CanChangeInfo:         true,
			CanInviteUsers:        true,
			CanPinMessages:        true,
		},
	})

	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("❌ 解除禁言失败: %v", err),
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("✅ 用户 %d 已解除禁言", userID),
	})
}

// PromoteAdmin 设置管理员
func PromoteAdmin(ctx context.Context, b *bot.Bot, update *models.Update) {
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

	var userID int64
	if update.Message.ReplyToMessage != nil {
		userID = update.Message.ReplyToMessage.From.ID
	} else {
		args := strings.Fields(update.Message.Text)
		if len(args) < 2 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "用法: /promote <用户ID> 或回复某条消息",
			})
			return
		}
		userID = int64(parseInt(args[1]))
	}

	if userID == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ 无效的用户ID",
		})
		return
	}

	// 提升为管理员
	_, err := b.PromoteChatMember(ctx, &bot.PromoteChatMemberParams{
		ChatID: update.Message.Chat.ID,
		UserID: userID,
		CanChangeInfo:        true,
		CanPostMessages:      true,
		CanEditMessages:     true,
		CanDeleteMessages:   true,
		CanInviteUsers:      true,
		CanRestrictMembers:   true,
		CanPinMessages:       true,
		CanPromoteMembers:     false, // 不能再次提升管理员
	})

	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("❌ 设置管理员失败: %v", err),
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("✅ 用户 %d 已提升为管理员", userID),
	})
}

// 辅助函数
func boolPtr(b bool) *bool {
	return &b
}

func getTime() *time.Time {
	t := time.Now()
	return &t
}
