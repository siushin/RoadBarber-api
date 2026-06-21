package service

import (
	"errors"
	"time"

	"roadbarber/backend/internal/config"
	"roadbarber/backend/internal/models"
)

type MerchantReviewService struct{}

// MerchantReviewItem 商家端评价（含顾客信息）
type MerchantReviewItem struct {
	models.Review
	CustomerName  string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
	ServiceName   string `json:"service_name"`
}

// MerchantListReviewsRequest 商家评价列表查询参数
type MerchantListReviewsRequest struct {
	Status   int8  `query:"status"`
	Page     int   `query:"page"`
	PageSize int   `query:"page_size"`
}

// ListReviews 商家端评价列表
func (s *MerchantReviewService) ListReviews(merchantUserID string, req *MerchantListReviewsRequest) ([]MerchantReviewItem, int64, error) {
	var merchant models.Merchant
	if err := config.GetDB().Where("user_id = ?", merchantUserID).First(&merchant).Error; err != nil {
		return nil, 0, errors.New("商家信息不存在")
	}

	db := config.GetDB().Table("reviews").
		Select(`reviews.*,
			users.nickname AS customer_name,
			users.phone AS customer_phone,
			services.name AS service_name`).
		Joins("LEFT JOIN users ON users.id = reviews.customer_id").
		Joins("LEFT JOIN services ON services.id = reviews.service_id").
		Where("reviews.merchant_id = ?", merchant.ID)

	if req.Status > 0 {
		db = db.Where("reviews.status = ?", req.Status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	type Result struct {
		models.Review
		CustomerName  string
		CustomerPhone string
		ServiceName   string
	}
	var rows []Result
	if err := db.Offset(offset).Limit(pageSize).Order("reviews.created_at DESC").Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]MerchantReviewItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, MerchantReviewItem{
			Review:        r.Review,
			CustomerName:  r.CustomerName,
			CustomerPhone: r.CustomerPhone,
			ServiceName:   r.ServiceName,
		})
	}
	return items, total, nil
}

// ReplyReview 商家回复评价
func (s *MerchantReviewService) ReplyReview(merchantUserID, reviewID, content string) error {
	var merchant models.Merchant
	if err := config.GetDB().Where("user_id = ?", merchantUserID).First(&merchant).Error; err != nil {
		return errors.New("商家信息不存在")
	}

	var review models.Review
	if err := config.GetDB().Where("id = ? AND merchant_id = ?", reviewID, merchant.ID).First(&review).Error; err != nil {
		return errors.New("评价不存在或无权回复")
	}

	return config.GetDB().Model(&review).Updates(map[string]interface{}{
		"reply_content": content,
		"reply_time":    time.Now(),
	}).Error
}