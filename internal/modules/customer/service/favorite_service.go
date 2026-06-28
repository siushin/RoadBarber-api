package service

import (
	"errors"

	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"

	"gorm.io/gorm"
)

type FavoriteService struct{}

// AddFavorite 收藏商家
func (s *FavoriteService) AddFavorite(userID, merchantID string) error {
	var existing models.Favorite
	err := config.GetDB().Where("user_id = ? AND merchant_id = ?", userID, merchantID).First(&existing).Error
	if err == nil {
		return errors.New("已收藏")
	}

	favorite := models.Favorite{
		UserID:     userID,
		MerchantID: merchantID,
	}
	return config.GetDB().Create(&favorite).Error
}

// RemoveFavorite 取消收藏
func (s *FavoriteService) RemoveFavorite(userID, merchantID string) error {
	return config.GetDB().Where("user_id = ? AND merchant_id = ?", userID, merchantID).
		Delete(&models.Favorite{}).Error
}

// ListMyFavorites 我的收藏
func (s *FavoriteService) ListMyFavorites(userID string, page, pageSize int) ([]models.Merchant, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	var total int64
	if err := config.GetDB().Model(&models.Favorite{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var merchants []models.Merchant
	err := config.GetDB().Table("merchants").
		Joins("JOIN favorites ON favorites.merchant_id = merchants.id").
		Where("favorites.user_id = ?", userID).
		Order("favorites.created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&merchants).Error
	return merchants, total, err
}

var _ = gorm.ErrRecordNotFound
