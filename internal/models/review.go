package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Review 评价表：顾客对服务的评价
type Review struct {
	ID           string     `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID    string     `gorm:"type:uuid;not null;uniqueIndex:idx_reviews_booking" json:"booking_id"`
	CustomerID   string     `gorm:"type:uuid;not null;index:idx_reviews_customer" json:"customer_id"`
	MerchantID   string     `gorm:"type:uuid;not null;index:idx_reviews_merchant" json:"merchant_id"`
	ShopID       *string    `gorm:"type:uuid" json:"shop_id"`
	ServiceID    *string    `gorm:"type:uuid" json:"service_id"`
	Rating       int8       `gorm:"type:smallint;not null" json:"rating"`
	Content      string     `gorm:"type:text" json:"content"`
	Images       string     `gorm:"type:jsonb" json:"images"`
	IsAnonymous  bool       `gorm:"type:boolean;default:false" json:"is_anonymous"`
	ReplyContent string     `gorm:"type:text" json:"reply_content"`
	ReplyTime    *time.Time `json:"reply_time"`
	Status       int8       `gorm:"type:smallint;default:1" json:"status"`
	CreatedAt    time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName 指定表名
func (Review) TableName() string {
	return "reviews"
}

// BeforeCreate 创建前生成 UUID
func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

// 评价状态常量
const (
	ReviewStatusVisible   = 1 // 显示
	ReviewStatusHidden    = 2 // 隐藏
)
