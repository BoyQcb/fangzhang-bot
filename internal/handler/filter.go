package handler

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var (
	sensitiveWords []string
	wordsMutex     sync.RWMutex
	wordsFile      = "data/sensitive_words.txt"
)

// LoadSensitiveWords 加载敏感词
func LoadSensitiveWords() {
	wordsMutex.Lock()
	defer wordsMutex.Unlock()

	data, err := os.ReadFile(wordsFile)
	if err != nil {
		// 文件不存在，创建空文件
		os.WriteFile(wordsFile, []byte(""), 0644)
		sensitiveWords = []string{}
		return
	}

	lines := strings.Split(string(data), "\n")
	sensitiveWords = []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			sensitiveWords = append(sensitiveWords, line)
		}
	}
}

// AddSensitiveWord 添加敏感词
func AddSensitiveWord(ctx context.Context, b *bot.Bot, update *models.Update) {
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
			Text:   "用法: /addword <敏感词>",
		})
		return
	}

	word := strings.Join(args[1:], " ")

	wordsMutex.Lock()
	defer wordsMutex.Unlock()

	// 检查是否已存在
	for _, w := range sensitiveWords {
		if w == word {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "⚠️ 该敏感词已存在",
			})
			return
		}
	}

	sensitiveWords = append(sensitiveWords, word)

	// 保存到文件
	saveSensitiveWords()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("✅ 已添加敏感词: %s", word),
	})
}

// DeleteSensitiveWord 删除敏感词
func DeleteSensitiveWord(ctx context.Context, b *bot.Bot, update *models.Update) {
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
			Text:   "用法: /delword <敏感词>",
		})
		return
	}

	word := strings.Join(args[1:], " ")

	wordsMutex.Lock()
	defer wordsMutex.Unlock()

	// 查找并删除
	found := false
	newWords := []string{}
	for _, w := range sensitiveWords {
		if w == word {
			found = true
			continue
		}
		newWords = append(newWords, w)
	}

	if !found {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "⚠️ 未找到该敏感词",
		})
		return
	}

	sensitiveWords = newWords
	saveSensitiveWords()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("✅ 已删除敏感词: %s", word),
	})
}

// ListSensitiveWords 列出所有敏感词
func ListSensitiveWords(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	wordsMutex.RLock()
	defer wordsMutex.RUnlock()

	if len(sensitiveWords) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "当前没有敏感词",
		})
		return
	}

	var sb strings.Builder
	sb.WriteString("📋 敏感词列表:\n\n")
	for i, word := range sensitiveWords {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, word))
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   sb.String(),
	})
}

// FilterMessage 过滤消息（中间件调用）
func FilterMessage(ctx context.Context, b *bot.Bot, update *models.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	wordsMutex.RLock()
	defer wordsMutex.RUnlock()

	// 检查敏感词
	for _, word := range sensitiveWords {
		if strings.Contains(update.Message.Text, word) {
			// 删除消息
			_, _ = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.ID,
			})

			// 警告用户
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("⚠️ 用户 %s 发送了包含敏感词的消息", update.Message.From.FirstName),
			})

			return true // 消息被过滤
		}
	}

	return false // 消息未被过滤
}

// 辅助函数
func saveSensitiveWords() {
	var sb strings.Builder
	for _, word := range sensitiveWords {
		sb.WriteString(word)
		sb.WriteString("\n")
	}

	os.WriteFile(wordsFile, []byte(sb.String()), 0644)
}
