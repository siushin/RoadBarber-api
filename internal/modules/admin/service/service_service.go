package service

import (
	"errors"

	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"
)

type AdminServiceItemService struct{}

// AdminListServicesRequest 服务列表查询参数
type AdminListServicesRequest struct {
	ShopID   string `query:"shop_id"`
	Category string `query:"category"`
	Status   int8   `query:"status"`
	Keyword  string `query:"keyword"`
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
}

// AdminServiceUpsertRequest 服务创建/更新请求
type AdminServiceUpsertRequest struct {
	ShopID      string  `json:"shop_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Duration    int     `json:"duration"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Images      string  `json:"images"`
	Status      int8    `json:"status"`
	SortOrder   int     `json:"sort_order"`
}

// ListServices 服务列表（管理员端）
func (s *AdminServiceItemService) ListServices(req *AdminListServicesRequest) ([]models.Service, int64, error) {
	db := config.GetDB().Model(&models.Service{})
	if req.ShopID != "" {
		db = db.Where("shop_id = ?", req.ShopID)
	}
	if req.Category != "" {
		db = db.Where("category = ?", req.Category)
	}
	if req.Status > 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
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

	var items []models.Service
	err := db.Offset(offset).Limit(pageSize).Order("sort_order ASC, created_at DESC").Find(&items).Error
	return items, total, err
}

// CreateService 创建服务
func (s *AdminServiceItemService) CreateService(req *AdminServiceUpsertRequest) (*models.Service, error) {
	status := req.Status
	if status == 0 {
		status = models.ServiceStatusOnSale
	}
	duration := req.Duration
	if duration == 0 {
		duration = 60
	}

	item := models.Service{
		ShopID:      req.ShopID,
		Name:        req.Name,
		Description: req.Description,
		Duration:    duration,
		Price:       req.Price,
		Category:    req.Category,
		Images:      req.Images,
		Status:      status,
		SortOrder:   req.SortOrder,
	}
	if err := config.GetDB().Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateService 更新服务
func (s *AdminServiceItemService) UpdateService(id string, req *AdminServiceUpsertRequest) (*models.Service, error) {
	var item models.Service
	if err := config.GetDB().Where("id = ?", id).First(&item).Error; err != nil {
		return nil, errors.New("服务不存在")
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Duration > 0 {
		updates["duration"] = req.Duration
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Images != "" {
		updates["images"] = req.Images
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}
	if req.SortOrder != 0 {
		updates["sort_order"] = req.SortOrder
	}

	if err := config.GetDB().Model(&item).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &item, nil
}