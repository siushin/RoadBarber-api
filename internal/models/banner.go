package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Banner 首页轮播图表
type Banner struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Image     string    `gorm:"type:varchar(500);not null" json:"image"`
	Title     string    `gorm:"type:varchar(50)" json:"title"`
	Subtitle  string    `gorm:"type:varchar(100)" json:"subtitle"`
	Text      string    `gorm:"type:varchar(200)" json:"text"`
	Align     string    `gorm:"type:varchar(20);not null;default:'flex-start'" json:"align"`
	LinkURL   string    `gorm:"type:varchar(500)" json:"link_url"`
	SortOrder int       `gorm:"type:int;not null;default:0" json:"sort_order"`
	Status    int8      `gorm:"type:smallint;not null;default:1" json:"status"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName 指定表名
func (Banner) TableName() string {
	return "banners"
}

// BeforeCreate 创建前生成 UUID
func (b *Banner) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// Banner 状态常量
const (
	BannerStatusActive  = 1 // 启用
	BannerStatusDisable = 2 // 禁用
)