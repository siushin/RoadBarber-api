package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Merchant 商家表：商家基本信息
type Merchant struct {
	ID              string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID          string    `gorm:"type:uuid;not null;uniqueIndex:idx_merchants_user" json:"user_id"`
	ShopID          *string   `gorm:"type:uuid;index:idx_merchants_shop" json:"shop_id"`
	Title           string    `gorm:"type:varchar(100)" json:"title"`
	Specialties     string    `gorm:"type:jsonb" json:"specialties"`
	ExperienceYears int       `gorm:"type:int;default:0" json:"experience_years"`
	Introduction    string    `gorm:"type:text" json:"introduction"`
	Rating          float64   `gorm:"type:decimal(2,1);default:5.0;index:idx_merchants_rating,sort:desc" json:"rating"`
	ReviewCount     int       `gorm:"type:int;default:0" json:"review_count"`
	ServiceCount    int       `gorm:"type:int;default:0" json:"service_count"`
	Avatar          string    `gorm:"type:varchar(500)" json:"avatar"`
	Status          int8      `gorm:"type:smallint;not null;default:1;index:idx_merchants_status" json:"status"`
	IsVerified      bool      `gorm:"type:boolean;default:false" json:"is_verified"`
	IsTop           bool      `gorm:"type:boolean;default:false" json:"is_top"`
	SortOrder       int       `gorm:"type:int;default:0" json:"sort_order"`
	Latitude        *float64  `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude       *float64  `gorm:"type:decimal(10,7)" json:"longitude"`
	StartPrice      float64   `gorm:"type:decimal(10,2);not null;default:0" json:"start_price"`
	BusinessHours   *string   `gorm:"type:varchar(50)" json:"business_hours"`
	Distance        *float64  `gorm:"type:decimal(10,2)" json:"distance"`
	AvailableSlots  int       `gorm:"type:int;not null;default:0" json:"available_slots"`
	CreatedAt       time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt       time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName 指定表名
func (Merchant) TableName() string {
	return "merchants"
}

// BeforeCreate 创建前生成 UUID
func (m *Merchant) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

// 商家状态常量
const (
	MerchantStatusNormal  = 1 // 正常
	MerchantStatusResting = 2 // 休息
	MerchantStatusLeft    = 3 // 离职
)
