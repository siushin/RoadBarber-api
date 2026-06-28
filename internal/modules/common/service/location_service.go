package service

import (
	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"
)

type LocationService struct{}

// ListProvinces 获取省份列表（parent_id = NULL）
func (s *LocationService) ListProvinces() ([]models.Location, error) {
	var locations []models.Location
	err := config.GetDB().Where("parent_id IS NULL").Order("sort_order ASC, name ASC").Find(&locations).Error
	return locations, err
}

// ListByParentID 获取指定地区的下级地区
func (s *LocationService) ListByParentID(parentID string) ([]models.Location, error) {
	var locations []models.Location
	err := config.GetDB().Where("parent_id = ?", parentID).Order("sort_order ASC, name ASC").Find(&locations).Error
	return locations, err
}

// GetTree 获取完整地区树
func (s *LocationService) GetTree() ([]models.Location, error) {
	var locations []models.Location
	err := config.GetDB().Order("level ASC, sort_order ASC, name ASC").Find(&locations).Error
	return locations, err
}
