package service

import (
	"roadbarber/backend/internal/config"
	"roadbarber/backend/internal/models"
)

// isAdminUser 判断 userID 是否对应超级管理员（role=3）
// 用于让管理员可访问商家端接口但避开商家专属数据
func isAdminUser(userID string) bool {
	var u models.User
	if err := config.GetDB().Where("id = ?", userID).First(&u).Error; err != nil {
		return false
	}
	return u.Role == models.RoleAdmin
}