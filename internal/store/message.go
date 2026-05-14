package store

import (
	"time"

	"github.com/xxx/fangzhang-bot/internal/model"
)

// CreateMessage 创建消息记录
func CreateMessage(msg *model.Message) error {
	return DB.Create(msg).Error
}

// GetMessagesByUser 获取用户的消息
func GetMessagesByUser(userID int64, offset, limit int) ([]model.Message, error) {
	var messages []model.Message
	err := DB.Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}

// GetMessagesByChat 获取群组的消息
func GetMessagesByChat(chatID int64, offset, limit int) ([]model.Message, error) {
	var messages []model.Message
	err := DB.Where("chat_id = ?", chatID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}

// CountMessagesByUser 统计用户消息数
func CountMessagesByUser(userID int64, since time.Time) (int64, error) {
	var count int64
	query := DB.Model(&model.Message{}).Where("user_id = ?", userID)
	if !since.IsZero() {
		query = query.Where("created_at >= ?", since)
	}
	err := query.Count(&count).Error
	return count, err
}

// CountMessagesByChat 统计群组消息数
func CountMessagesByChat(chatID int64, since time.Time) (int64, error) {
	var count int64
	query := DB.Model(&model.Message{}).Where("chat_id = ?", chatID)
	if !since.IsZero() {
		query = query.Where("created_at >= ?", since)
	}
	err := query.Count(&count).Error
	return count, err
}

// GetTopUsers 获取活跃用户排行
func GetTopUsers(chatID int64, since time.Time, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := DB.Model(&model.Message{}).
		Select("user_id, COUNT(*) as message_count").
		Where("chat_id = ?", chatID)

	if !since.IsZero() {
		query = query.Where("created_at >= ?", since)
	}

	err := query.Group("user_id").
		Order("message_count DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// DeleteOldMessages 删除旧消息
func DeleteOldMessages(before time.Time) error {
	return DB.Where("created_at < ?", before).Delete(&model.Message{}).Error
}
