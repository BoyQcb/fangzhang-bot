package handler

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// LotterySession 抽奖会话
type LotterySession struct {
	ID          int
	ChatID      int64
	MessageID   int
	Prize       string
	Participants map[int64]bool
	IsActive    bool
	CreatorID   int64
	EndTime     time.Time
}

var (
	lotterySessions []LotterySession
	lotteryMutex   sync.RWMutex
	lotteryIDCounter int
)

// StartLottery 开始抽奖
func StartLottery(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	// 用法: /lottery <奖品名称> [持续时间(秒)]
	args := strings.Fields(update.Message.Text)
	if len(args) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "用法: /lottery <奖品名称> [持续时间(秒)]\n例如: /lottery \"iPhone 15\" 300",
		})
		return
	}

	prize := args[1]
	duration := 300 // 默认5分钟

	if len(args) >= 3 {
		duration = parseInt(args[2])
		if duration <= 0 {
			duration = 300
		}
	}

	// 创建抽奖
	lotteryMutex.Lock()
	lotteryIDCounter++
	session := LotterySession{
		ID:          lotteryIDCounter,
		ChatID:      update.Message.Chat.ID,
		MessageID:   update.Message.ID,
		Prize:       prize,
		Participants: make(map[int64]bool),
		IsActive:    true,
		CreatorID:   update.Message.From.ID,
		EndTime:     time.Now().Add(time.Duration(duration) * time.Second),
	}
	lotterySessions = append(lotterySessions, session)
	lotteryMutex.Unlock()

	// 发送抽奖消息
	msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf(
			"🎉 抽奖开始！\n\n"+
				"奖品: %s\n"+
				"持续时间: %d 秒\n"+
				"参与方式: 回复此消息并输入 /join\n"+
				"截止时间: %s\n\n"+
				"发起人: %s",
			prize,
			duration,
			session.EndTime.Format("15:04:05"),
			update.Message.From.FirstName,
		),
	})

	// 更新MessageID
	lotteryMutex.Lock()
	for i := range lotterySessions {
		if lotterySessions[i].ID == session.ID {
			lotterySessions[i].MessageID = msg.ID
			break
		}
	}
	lotteryMutex.Unlock()

	// 定时开奖
	time.AfterFunc(time.Duration(duration)*time.Second, func() {
		DrawLotteryByID(ctx, b, session.ID)
	})
}

// JoinLottery 参与抽奖
func JoinLottery(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.ReplyToMessage == nil {
		return
	}

	// 查找对应的抽奖
	lotteryMutex.Lock()
	defer lotteryMutex.Unlock()

	for i := range lotterySessions {
		if lotterySessions[i].ChatID == update.Message.Chat.ID &&
			lotterySessions[i].MessageID == update.Message.ReplyToMessage.ID &&
			lotterySessions[i].IsActive {
			// 检查是否已参与
			if lotterySessions[i].Participants[update.Message.From.ID] {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "⚠️ 你已经参与过这个抽奖了",
				})
				return
			}

			// 添加参与者
			lotterySessions[i].Participants[update.Message.From.ID] = true

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("✅ %s 已参与抽奖！", update.Message.From.FirstName),
			})
			return
		}
	}
}

// DrawLottery 手动开奖
func DrawLottery(ctx context.Context, b *bot.Bot, update *models.Update) {
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
			Text:   "用法: /draw <抽奖ID>",
		})
		return
	}

	lotteryID := parseInt(args[1])
	DrawLotteryByID(ctx, b, lotteryID)
}

// DrawLotteryByID 根据ID开奖
func DrawLotteryByID(ctx context.Context, b *bot.Bot, lotteryID int) {
	lotteryMutex.Lock()
	defer lotteryMutex.Unlock()

	for i := range lotterySessions {
		if lotterySessions[i].ID == lotteryID && lotterySessions[i].IsActive {
			if len(lotterySessions[i].Participants) == 0 {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: lotterySessions[i].ChatID,
					Text:   "😢 没有人参与抽奖",
				})
				lotterySessions[i].IsActive = false
				return
			}

			// 随机选择中奖者
			participantIDs := []int64{}
			for uid := range lotterySessions[i].Participants {
				participantIDs = append(participantIDs, uid)
			}

			rand.Seed(time.Now().UnixNano())
			winnerID := participantIDs[rand.Intn(len(participantIDs))]

			// 发送开奖消息
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: lotterySessions[i].ChatID,
				Text: fmt.Sprintf(
					"🎊 开奖啦！\n\n"+
						"奖品: %s\n"+
						"中奖者ID: %d\n"+
						"参与人数: %d\n\n"+
						"恭喜中奖者！",
					lotterySessions[i].Prize,
					winnerID,
					len(participantIDs),
				),
			})

			lotterySessions[i].IsActive = false
			return
		}
	}
}
