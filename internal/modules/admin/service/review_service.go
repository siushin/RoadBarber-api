package service

import (
	"errors"

	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"
)

type AdminReviewService struct{}

// AdminReviewItem 管理员端评价（含顾客名/商家title/服务名）
type AdminReviewItem struct {
	models.Review
	CustomerName  string `json:"customer_name"`
	MerchantTitle string `json:"merchant_title"`
	ServiceName   string `json:"service_name"`
}

// AdminListReviewsRequest 管理员评价列表查询参数
type AdminListReviewsRequest struct {
	MerchantID string `query:"merchant_id"`
	Status     int8   `query:"status"`
	Page       int    `query:"page"`
	PageSize   int    `query:"page_size"`
}

// ListReviews 管理员评价列表（含关联信息）
func (s *AdminReviewService) ListReviews(req *AdminListReviewsRequest) ([]AdminReviewItem, int64, error) {
	db := config.GetDB().Table("reviews").
		Select(`reviews.*,
			users.nickname AS customer_name,
			merchants.title AS merchant_title,
			services.name AS service_name`).
		Joins("LEFT JOIN users ON users.id = reviews.customer_id").
		Joins("LEFT JOIN merchants ON merchants.id = reviews.merchant_id").
		Joins("LEFT JOIN services ON services.id = reviews.service_id")

	if req.MerchantID != "" {
		db = db.Where("reviews.merchant_id = ?", req.MerchantID)
	}
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
		MerchantTitle string
		ServiceName   string
	}
	var rows []Result
	if err := db.Offset(offset).Limit(pageSize).Order("reviews.created_at DESC").Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]AdminReviewItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, AdminReviewItem{
			Review:        r.Review,
			CustomerName:  r.CustomerName,
			MerchantTitle: r.MerchantTitle,
			ServiceName:   r.ServiceName,
		})
	}
	return items, total, nil
}

// DeleteReview 软删除评价（status=2）
func (s *AdminReviewService) DeleteReview(id string) error {
	var review models.Review
	if err := config.GetDB().Where("id = ?", id).First(&review).Error; err != nil {
		return errors.New("评价不存在")
	}
	return config.GetDB().Model(&review).Update("status", models.ReviewStatusHidden).Error
}