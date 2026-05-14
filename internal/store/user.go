package store

import (
	"github.com/xxx/fangzhang-bot/internal/model"
)

// CreateUser 创建用户
func CreateUser(user *model.User) error {
	return DB.Create(user).Error
}

// GetUserByID 根据ID获取用户
func GetUserByID(userID int64) (*model.User, error) {
	var user model.User
	err := DB.Where("user_id = ?", userID).First(&user).Error
	return &user, err
}

// UpdateUser 更新用户
func UpdateUser(user *model.User) error {
	return DB.Save(user).Error
}

// DeleteUser 删除用户
func DeleteUser(userID int64) error {
	return DB.Where("user_id = ?", userID).Delete(&model.User{}).Error
}

// ListUsers 列出所有用户
func ListUsers(offset, limit int) ([]model.User, error) {
	var users []model.User
	err := DB.Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// CountUsers 统计用户数量
func CountUsers() (int64, error) {
	var count int64
	err := DB.Model(&model.User{}).Count(&count).Error
	return count, err
}
