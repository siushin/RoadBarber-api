package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Favorite 收藏表：顾客收藏商家
type Favorite struct {
	ID         string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     string    `gorm:"type:uuid;not null;uniqueIndex:idx_favorites_user_merchant,priority:1;index:idx_favorites_user" json:"user_id"`
	MerchantID string    `gorm:"type:uuid;not null;uniqueIndex:idx_favorites_user_merchant,priority:2;index:idx_favorites_merchant" json:"merchant_id"`
	CreatedAt  time.Time `gorm:"not null;default:now()" json:"created_at"`
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorites"
}

// BeforeCreate 创建前生成 UUID
func (f *Favorite) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return nil
}
