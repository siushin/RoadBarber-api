package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service 服务项目表：店铺提供的服务
type Service struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	ShopID      string    `gorm:"type:uuid;not null;index:idx_services_shop" json:"shop_id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Duration    int       `gorm:"type:int;not null;default:60" json:"duration"`
	Price       float64   `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	Category    string    `gorm:"type:varchar(50);index:idx_services_category" json:"category"`
	Images      string    `gorm:"type:jsonb" json:"images"`
	Status      int8      `gorm:"type:smallint;not null;default:1;index:idx_services_status" json:"status"`
	SortOrder   int       `gorm:"type:int;default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName 指定表名
func (Service) TableName() string {
	return "services"
}

// BeforeCreate 创建前生成 UUID
func (s *Service) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// 服务状态常量
const (
	ServiceStatusOnSale  = 1 // 上架
	ServiceStatusOffSale = 2 // 下架
)
