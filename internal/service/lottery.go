package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/xxx/fangzhang-bot/internal/model"
	"github.com/xxx/fangzhang-bot/internal/store"
)

// LotteryService 抽奖服务
type LotteryService struct {
	db *gorm.DB
}

// NewLotteryService 创建抽奖服务
func NewLotteryService(db *gorm.DB) *LotteryService {
	return &LotteryService{db: db}
}

// CreateLottery 创建抽奖
func (s *LotteryService) CreateLottery(chatID int64, messageID int, prize string, createdBy int64, duration int) (*model.Lottery, error) {
	lottery := &model.Lottery{
		ChatID:    chatID,
		MessageID: messageID,
		Prize:     prize,
		IsActive:  true,
		CreatedBy: createdBy,
		EndTime:   time.Now().Add(time.Duration(duration) * time.Second),
	}

	err := store.CreateLottery(lottery)
	if err != nil {
		return nil, fmt.Errorf("创建抽奖失败: %w", err)
	}

	return lottery, nil
}

// GetLottery 获取抽奖
func (s *LotteryService) GetLottery(id uint) (*model.Lottery, error) {
	return store.GetLotteryByID(id)
}

// UpdateLottery 更新抽奖
func (s *LotteryService) UpdateLottery(lottery *model.Lottery) error {
	return store.UpdateLottery(lottery)
}

// DeleteLottery 删除抽奖
func (s *LotteryService) DeleteLottery(id uint) error {
	return store.DeleteLottery(id)
}

// ListLotteries 列出抽奖
func (s *LotteryService) ListLotteries(chatID int64, offset, limit int) ([]model.Lottery, error) {
	return store.ListLotteries(chatID, offset, limit)
}

// ListActiveLotteries 列出活跃的抽奖
func (s *LotteryService) ListActiveLotteries(chatID int64) ([]model.Lottery, error) {
	return store.ListActiveLotteries(chatID)
}

// JoinLottery 参与抽奖
func (s *LotteryService) JoinLottery(lotteryID uint, userID int64) error {
	participant := &model.Participant{
		LotteryID: lotteryID,
		UserID:    userID,
	}

	err := store.AddParticipant(participant)
	if err != nil {
		return fmt.Errorf("参与抽奖失败: %w", err)
	}

	// 更新参与者数量
	lottery, err := store.GetLotteryByID(lotteryID)
	if err != nil {
		return err
	}

	count, err := store.CountParticipants(lotteryID)
	if err != nil {
		return err
	}

	lottery.ParticipantCount = int(count)
	return store.UpdateLottery(lottery)
}

// DrawLottery 抽奖
func (s *LotteryService) DrawLottery(lotteryID uint) (*model.Lottery, error) {
	// 获取抽奖
	lottery, err := store.GetLotteryByID(lotteryID)
	if err != nil {
		return nil, fmt.Errorf("获取抽奖失败: %w", err)
	}

	if !lottery.IsActive {
		return nil, fmt.Errorf("抽奖已结束")
	}

	// 获取参与者
	participants, err := store.GetParticipants(lotteryID)
	if err != nil {
		return nil, fmt.Errorf("获取参与者失败: %w", err)
	}

	if len(participants) == 0 {
		return nil, fmt.Errorf("没有参与者")
	}

	// 随机选择中奖者
	rand.Seed(time.Now().UnixNano())
	winnerIndex := rand.Intn(len(participants))
	winnerID := participants[winnerIndex].UserID

	// 更新抽奖结果
	lottery.WinnerID = winnerID
	lottery.IsActive = false
	err = store.UpdateLottery(lottery)
	if err != nil {
		return nil, fmt.Errorf("更新抽奖结果失败: %w", err)
	}

	return lottery, nil
}

// GetParticipants 获取参与者列表
func (s *LotteryService) GetParticipants(lotteryID uint) ([]model.Participant, error) {
	return store.GetParticipants(lotteryID)
}

// CountParticipants 统计参与者数量
func (s *LotteryService) CountParticipants(lotteryID uint) (int, error) {
	count, err := store.CountParticipants(lotteryID)
	return int(count), err
}
