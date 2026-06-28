package service

import (
	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"
)

type ShopService struct{}

// ListShopsRequest 店铺列表查询参数
type ListShopsRequest struct {
	Keyword  string  `query:"keyword"`
	LocationID string `query:"location_id"`
	Latitude  float64 `query:"latitude"`
	Longitude float64 `query:"longitude"`
	Page      int    `query:"page"`
	PageSize  int    `query:"page_size"`
}

// ListShops 店铺列表
func (s *ShopService) ListShops(req *ListShopsRequest) ([]models.Shop, int64, error) {
	db := config.GetDB().Model(&models.Shop{}).Where("status = ?", models.ShopStatusNormal)

	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR address LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.LocationID != "" {
		db = db.Where("location_id = ?", req.LocationID)
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

	var shops []models.Shop
	err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&shops).Error
	return shops, total, err
}

// GetShopDetail 店铺详情
func (s *ShopService) GetShopDetail(id string) (*models.Shop, error) {
	var shop models.Shop
	err := config.GetDB().Where("id = ?", id).First(&shop).Error
	return &shop, err
}

// ListMerchantsByShop 店铺下的商家
func (s *ShopService) ListMerchantsByShop(shopID string) ([]models.Merchant, error) {
	var merchants []models.Merchant
	err := config.GetDB().Where("shop_id = ? AND status = ?", shopID, models.MerchantStatusNormal).
		Order("rating DESC").Find(&merchants).Error
	return merchants, err
}

// ListServicesByShop 店铺下的服务
func (s *ShopService) ListServicesByShop(shopID string) ([]models.Service, error) {
	var services []models.Service
	err := config.GetDB().Where("shop_id = ? AND status = ?", shopID, models.ServiceStatusOnSale).
		Order("sort_order ASC, created_at DESC").Find(&services).Error
	return services, err
}
