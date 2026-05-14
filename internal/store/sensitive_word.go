package store

import (
	"github.com/xxx/fangzhang-bot/internal/model"
)

// CreateSensitiveWord 创建敏感词
func CreateSensitiveWord(word *model.SensitiveWord) error {
	return DB.Create(word).Error
}

// GetSensitiveWordByID 根据ID获取敏感词
func GetSensitiveWordByID(id uint) (*model.SensitiveWord, error) {
	var word model.SensitiveWord
	err := DB.First(&word, id).Error
	return &word, err
}

// GetSensitiveWordByWord 根据词语获取敏感词
func GetSensitiveWordByWord(word string) (*model.SensitiveWord, error) {
	var sensitiveWord model.SensitiveWord
	err := DB.Where("word = ?", word).First(&sensitiveWord).Error
	return &sensitiveWord, err
}

// UpdateSensitiveWord 更新敏感词
func UpdateSensitiveWord(word *model.SensitiveWord) error {
	return DB.Save(word).Error
}

// DeleteSensitiveWord 删除敏感词
func DeleteSensitiveWord(id uint) error {
	return DB.Delete(&model.SensitiveWord{}, id).Error
}

// ListSensitiveWords 列出所有敏感词
func ListSensitiveWords(offset, limit int) ([]model.SensitiveWord, error) {
	var words []model.SensitiveWord
	query := DB.Order("created_at DESC")
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Find(&words).Error
	return words, err
}

// CountSensitiveWords 统计敏感词数量
func CountSensitiveWords() (int64, error) {
	var count int64
	err := DB.Model(&model.SensitiveWord{}).Count(&count).Error
	return count, err
}

// GetAllSensitiveWords 获取所有敏感词（用于过滤）
func GetAllSensitiveWords() ([]model.SensitiveWord, error) {
	var words []model.SensitiveWord
	err := DB.Find(&words).Error
	return words, err
}
