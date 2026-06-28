package service

import (
	"errors"

	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"

	"gorm.io/gorm"
)

type ReviewService struct{}

// CreateReviewRequest 提交评价请求
type CreateReviewRequest struct {
	BookingID   string   `json:"booking_id" validate:"required"`
	Rating      int8     `json:"rating" validate:"required,min=1,max=5"`
	Content     string   `json:"content"`
	Images      []string `json:"images"`
	IsAnonymous bool     `json:"is_anonymous"`
}

// CreateReview 提交评价
func (s *ReviewService) CreateReview(customerID string, req *CreateReviewRequest) (*models.Review, error) {
	// 校验预约
	var booking models.Booking
	err := config.GetDB().Where("id = ?", req.BookingID).First(&booking).Error
	if err != nil {
		return nil, errors.New("预约不存在")
	}
	if booking.CustomerID != customerID {
		return nil, errors.New("无权评价此预约")
	}
	if booking.Status != models.BookingStatusCompleted {
		return nil, errors.New("仅已完成的预约可以评价")
	}

	// 校验是否已评价
	var count int64
	config.GetDB().Model(&models.Review{}).Where("booking_id = ?", req.BookingID).Count(&count)
	if count > 0 {
		return nil, errors.New("该预约已评价")
	}

	imagesJSON := ""
	if len(req.Images) > 0 {
		// 简化处理：实际应序列化为 JSON
		imagesJSON = joinStrings(req.Images, ",")
	}

	review := models.Review{
		BookingID:   req.BookingID,
		CustomerID:  customerID,
		MerchantID:  booking.MerchantID,
		ShopID:      booking.ShopID,
		ServiceID:   &booking.ServiceID,
		Rating:      req.Rating,
		Content:     req.Content,
		Images:      imagesJSON,
		IsAnonymous: req.IsAnonymous,
		Status:      models.ReviewStatusVisible,
	}

	err = config.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&review).Error; err != nil {
			return err
		}
		// 更新商家评分
		return updateMerchantRating(tx, booking.MerchantID)
	})

	if err != nil {
		return nil, err
	}

	return &review, nil
}

// updateMerchantRating 更新商家评分
func updateMerchantRating(tx *gorm.DB, merchantID string) error {
	var avgRating float64
	var count int64
	if err := tx.Model(&models.Review{}).
		Where("merchant_id = ? AND status = ?", merchantID, models.ReviewStatusVisible).
		Select("AVG(rating), COUNT(*)").Row().Scan(&avgRating, &count); err != nil {
		return err
	}
	return tx.Model(&models.Merchant{}).Where("id = ?", merchantID).Updates(map[string]interface{}{
		"rating":       avgRating,
		"review_count": count,
	}).Error
}

func joinStrings(arr []string, sep string) string {
	result := ""
	for i, s := range arr {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

// ListByMerchant 商家评价列表
func (s *ReviewService) ListByMerchant(merchantID string, page, pageSize int) ([]models.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	var total int64
	if err := config.GetDB().Model(&models.Review{}).
		Where("merchant_id = ? AND status = ?", merchantID, models.ReviewStatusVisible).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var reviews []models.Review
	err := config.GetDB().Where("merchant_id = ? AND status = ?", merchantID, models.ReviewStatusVisible).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&reviews).Error
	return reviews, total, err
}
