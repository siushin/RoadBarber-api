package service

import (
	"errors"

	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"
)

type AdminShopService struct{}

// AdminListShopsRequest 管理员店铺列表查询参数
type AdminListShopsRequest struct {
	Keyword  string `query:"keyword"`
	Status   int8   `query:"status"`
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
}

// AdminShopUpsertRequest 店铺创建/更新请求
type AdminShopUpsertRequest struct {
	Name          string  `json:"name"`
	LocationID    *string `json:"location_id"`
	Address       string  `json:"address"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	Phone         string  `json:"phone"`
	BusinessHours string  `json:"business_hours"`
	Images        string  `json:"images"`
	Description   string  `json:"description"`
	Status        int8    `json:"status"`
}

// ListShops 管理员店铺列表
func (s *AdminShopService) ListShops(req *AdminListShopsRequest) ([]models.Shop, int64, error) {
	db := config.GetDB().Model(&models.Shop{})
	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR address LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status > 0 {
		db = db.Where("status = ?", req.Status)
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

// CreateShop 创建店铺
func (s *AdminShopService) CreateShop(req *AdminShopUpsertRequest, creatorID string) (*models.Shop, error) {
	status := req.Status
	if status == 0 {
		status = models.ShopStatusNormal
	}

	shop := models.Shop{
		Name:          req.Name,
		LocationID:    req.LocationID,
		Address:       req.Address,
		Longitude:     req.Longitude,
		Latitude:      req.Latitude,
		Phone:         req.Phone,
		BusinessHours: req.BusinessHours,
		Images:        req.Images,
		Description:   req.Description,
		Status:        status,
		CreatorID:     creatorID,
	}
	if err := config.GetDB().Create(&shop).Error; err != nil {
		return nil, err
	}
	return &shop, nil
}

// UpdateShop 更新店铺
func (s *AdminShopService) UpdateShop(id string, req *AdminShopUpsertRequest) (*models.Shop, error) {
	var shop models.Shop
	if err := config.GetDB().Where("id = ?", id).First(&shop).Error; err != nil {
		return nil, errors.New("店铺不存在")
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.LocationID != nil {
		updates["location_id"] = req.LocationID
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Longitude != 0 {
		updates["longitude"] = req.Longitude
	}
	if req.Latitude != 0 {
		updates["latitude"] = req.Latitude
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.BusinessHours != "" {
		updates["business_hours"] = req.BusinessHours
	}
	if req.Images != "" {
		updates["images"] = req.Images
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}

	if err := config.GetDB().Model(&shop).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &shop, nil
}