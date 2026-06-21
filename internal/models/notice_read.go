package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NoticeRead 通知已读记录表
type NoticeRead struct {
	ID       string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID   string    `gorm:"type:uuid;not null;uniqueIndex:idx_notice_reads_user_notice" json:"user_id"`
	NoticeID string    `gorm:"type:uuid;not null;uniqueIndex:idx_notice_reads_user_notice" json:"notice_id"`
	ReadAt   time.Time `gorm:"not null;default:now()" json:"read_at"`
}

// TableName 指定表名
func (NoticeRead) TableName() string {
	return "notice_reads"
}

// BeforeCreate 创建前生成 UUID
func (n *NoticeRead) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}