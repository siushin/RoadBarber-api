package service

import (
	"roadbarber/backend/internal/config"
	"roadbarber/backend/internal/models"
)

type ServiceCatalogService struct{}

// ListByShop 店铺下的服务
func (s *ServiceCatalogService) ListByShop(shopID string) ([]models.Service, error) {
	var services []models.Service
	err := config.GetDB().Where("shop_id = ? AND status = ?", shopID, models.ServiceStatusOnSale).
		Order("sort_order ASC, created_at DESC").Find(&services).Error
	return services, err
}
