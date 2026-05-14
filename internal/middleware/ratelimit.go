package middleware

import (
	"context"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// RateLimiter 频率限制器
type RateLimiter struct {
	mu       sync.Mutex
	visits   map[int64]int
	lastSeen map[int64]time.Time
	limit    int
	window   time.Duration
}

// NewRateLimiter 创建频率限制器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visits:   make(map[int64]int),
		lastSeen: make(map[int64]time.Time),
		limit:    limit,
		window:   window,
	}

	// 定期清理过期记录
	go func() {
		for {
			time.Sleep(window)
			rl.mu.Lock()
			now := time.Now()
			for userID, last := range rl.lastSeen {
				if now.Sub(last) > window {
					delete(rl.visits, userID)
					delete(rl.lastSeen, userID)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(userID int64) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	last, exists := rl.lastSeen[userID]

	// 如果超过时间窗口，重置计数
	if !exists || now.Sub(last) > rl.window {
		rl.visits[userID] = 1
		rl.lastSeen[userID] = now
		return true
	}

	// 检查是否超过限制
	if rl.visits[userID] >= rl.limit {
		return false
	}

	// 增加计数
	rl.visits[userID]++
	rl.lastSeen[userID] = now
	return true
}

// RateLimit 频率限制中间件
func RateLimit(limit int, window time.Duration) bot.Middleware {
	limiter := NewRateLimiter(limit, window)

	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if update.Message == nil {
				next(ctx, b, update)
				return
			}

			userID := update.Message.From.ID

			if !limiter.Allow(userID) {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "⚠️ 你发送消息太频繁，请稍后再试",
				})
				return
			}

			next(ctx, b, update)
		}
	}
}
