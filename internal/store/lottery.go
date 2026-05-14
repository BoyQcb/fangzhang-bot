package store

import (
	"fmt"
	"time"

	"github.com/xxx/fangzhang-bot/internal/model"
)

// CreateLottery 创建抽奖
func CreateLottery(lottery *model.Lottery) error {
	return DB.Create(lottery).Error
}

// GetLotteryByID 根据ID获取抽奖
func GetLotteryByID(id uint) (*model.Lottery, error) {
	var lottery model.Lottery
	err := DB.First(&lottery, id).Error
	return &lottery, err
}

// UpdateLottery 更新抽奖
func UpdateLottery(lottery *model.Lottery) error {
	return DB.Save(lottery).Error
}

// DeleteLottery 删除抽奖
func DeleteLottery(id uint) error {
	// 先删除参与者
	DB.Where("lottery_id = ?", id).Delete(&model.Participant{})
	// 再删除抽奖
	return DB.Delete(&model.Lottery{}, id).Error
}

// ListLotteries 列出抽奖
func ListLotteries(chatID int64, offset, limit int) ([]model.Lottery, error) {
	var lotteries []model.Lottery
	query := DB.Where("chat_id = ?", chatID)
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Find(&lotteries).Error
	return lotteries, err
}

// ListActiveLotteries 列出活跃的抽奖
func ListActiveLotteries(chatID int64) ([]model.Lottery, error) {
	var lotteries []model.Lottery
	err := DB.Where("chat_id = ? AND is_active = ?", chatID, true).Find(&lotteries).Error
	return lotteries, err
}

// AddParticipant 添加参与者
func AddParticipant(participant *model.Participant) error {
	// 检查是否已参与
	var count int64
	DB.Model(&model.Participant{}).
		Where("lottery_id = ? AND user_id = ?", participant.LotteryID, participant.UserID).
		Count(&count)

	if count > 0 {
		return fmt.Errorf("用户已参与此抽奖")
	}

	return DB.Create(participant).Error
}

// GetParticipants 获取参与者列表
func GetParticipants(lotteryID uint) ([]model.Participant, error) {
	var participants []model.Participant
	err := DB.Where("lottery_id = ?", lotteryID).Find(&participants).Error
	return participants, err
}

// CountParticipants 统计参与者数量
func CountParticipants(lotteryID uint) (int64, error) {
	var count int64
	err := DB.Model(&model.Participant{}).Where("lottery_id = ?", lotteryID).Count(&count).Error
	return count, err
}

// DeleteOldLotteries 删除旧抽奖
func DeleteOldLotteries(before time.Time) error {
	// 先删除参与者
	var oldLotteries []model.Lottery
	DB.Where("created_at < ?", before).Find(&oldLotteries)

	for _, lottery := range oldLotteries {
		DB.Where("lottery_id = ?", lottery.ID).Delete(&model.Participant{})
	}

	// 再删除抽奖
	return DB.Where("created_at < ?", before).Delete(&model.Lottery{}).Error
}
