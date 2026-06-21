package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notice 首页滚动公告表
// icon 存 lucide 图标名（gift / truck / bell / star 等），前端作为 wd-notice-bar 的 prefix 显示；
// text_color 给前端独立设置文字颜色（背景统一透明）。
type Notice struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Content   string    `gorm:"type:varchar(500);not null" json:"content"`
	Icon      *string   `gorm:"type:varchar(50)" json:"icon"`
	TextColor *string   `gorm:"type:varchar(20)" json:"text_color"`
	SortOrder int       `gorm:"type:int;not null;default:0" json:"sort_order"`
	Status    int8      `gorm:"type:smallint;not null;default:1" json:"status"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName 指定表名
func (Notice) TableName() string {
	return "notices"
}

// BeforeCreate 创建前生成 UUID
func (n *Notice) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}

// Notice 状态常量
const (
	NoticeStatusActive  = 1 // 启用
	NoticeStatusDisable = 2 // 禁用
)