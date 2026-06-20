package service

import (
	"roadbarber/backend/internal/config"
	"roadbarber/backend/internal/models"
)

type MerchantService struct{}

// ListMerchantsRequest 商家列表查询参数
type ListMerchantsRequest struct {
	ShopID    string `query:"shop_id"`
	Category  string `query:"category"`
	Keyword   string `query:"keyword"`
	SortBy    string `query:"sort_by"`     // rating, distance, service_count
	Latitude  float64 `query:"latitude"`
	Longitude float64 `query:"longitude"`
	Page      int    `query:"page"`
	PageSize  int    `query:"page_size"`
}

// ListMerchants 商家列表
func (s *MerchantService) ListMerchants(req *ListMerchantsRequest) ([]models.Merchant, int64, error) {
	db := config.GetDB().Model(&models.Merchant{}).Where("status = ?", models.MerchantStatusNormal)

	if req.ShopID != "" {
		db = db.Where("shop_id = ?", req.ShopID)
	}
	if req.Keyword != "" {
		db = db.Where("title LIKE ? OR introduction LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 排序
	switch req.SortBy {
	case "service_count":
		db = db.Order("service_count DESC")
	case "rating":
		db = db.Order("rating DESC, review_count DESC")
	default:
		db = db.Order("is_top DESC, sort_order DESC, rating DESC")
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

	var merchants []models.Merchant
	err := db.Offset(offset).Limit(pageSize).Find(&merchants).Error
	return merchants, total, err
}

// GetMerchantDetail 商家详情
func (s *MerchantService) GetMerchantDetail(id string) (*models.Merchant, error) {
	var merchant models.Merchant
	err := config.GetDB().Where("id = ?", id).First(&merchant).Error
	return &merchant, err
}
