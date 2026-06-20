package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Shop 店铺表：商家入驻的店铺
type Shop struct {
	ID            string    `gorm:"type:uuid;primaryKey" json:"id"`
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	LocationID    *string   `gorm:"type:uuid;index:idx_shops_location" json:"location_id"`
	Address       string    `gorm:"type:varchar(255);not null" json:"address"`
	Longitude     float64   `gorm:"type:decimal(10,7)" json:"longitude"`
	Latitude      float64   `gorm:"type:decimal(10,7)" json:"latitude"`
	Phone         string    `gorm:"type:varchar(20)" json:"phone"`
	BusinessHours string    `gorm:"type:jsonb" json:"business_hours"`
	Images        string    `gorm:"type:jsonb" json:"images"`
	Description   string    `gorm:"type:text" json:"description"`
	Status        int8      `gorm:"type:smallint;not null;default:1;index:idx_shops_status" json:"status"`
	CreatorID     string    `gorm:"type:uuid;index:idx_shops_creator" json:"creator_id"`
	CreatedAt     time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName 指定表名
func (Shop) TableName() string {
	return "shops"
}

// BeforeCreate 创建前生成 UUID
func (s *Shop) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// 店铺状态常量
const (
	ShopStatusNormal    = 1 // 正常
	ShopStatusClosed    = 2 // 歇业
	ShopStatusDisabled  = 3 // 停用
)
