package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MerchantService 商家服务关联表
type MerchantService struct {
	ID         string    `gorm:"type:uuid;primaryKey" json:"id"`
	MerchantID string    `gorm:"type:uuid;not null;uniqueIndex:idx_ms_merchant_service,priority:1;index:idx_merchant_services_merchant" json:"merchant_id"`
	ServiceID  string    `gorm:"type:uuid;not null;uniqueIndex:idx_ms_merchant_service,priority:2;index:idx_merchant_services_service" json:"service_id"`
	Price      *float64  `gorm:"type:decimal(10,2)" json:"price"`
	CreatedAt  time.Time `gorm:"not null;default:now()" json:"created_at"`
}

// TableName 指定表名
func (MerchantService) TableName() string {
	return "merchant_services"
}

// BeforeCreate 创建前生成 UUID
func (m *MerchantService) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
