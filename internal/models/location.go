package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Location 地区表：支持省市区街道小区楼栋多级下钻
type Location struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	ParentID  *string   `gorm:"type:uuid;index:idx_locations_parent" json:"parent_id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Code      string    `gorm:"type:varchar(20);not null;uniqueIndex:idx_locations_code" json:"code"`
	Level     int8      `gorm:"type:smallint;not null;index:idx_locations_level" json:"level"`
	SortOrder int       `gorm:"type:int;default:0" json:"sort_order"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName 指定表名
func (Location) TableName() string {
	return "locations"
}

// BeforeCreate 创建前生成 UUID
func (l *Location) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}

// 地区级别常量
const (
	LocationLevelProvince   = 1 // 省
	LocationLevelCity       = 2 // 市
	LocationLevelDistrict   = 3 // 区/县
	LocationLevelStreet     = 4 // 街道/乡镇
	LocationLevelVillage    = 5 // 村/居委会
	LocationLevelCommunity  = 6 // 小区
	LocationLevelBuilding   = 7 // 楼栋
)
