package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Schedule 排班时段表：商家发布的可用时段
type Schedule struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	MerchantID  string    `gorm:"type:uuid;not null;index:idx_schedules_merchant" json:"merchant_id"`
	WorkDate    time.Time `gorm:"type:date;not null;index:idx_schedules_date" json:"work_date"`
	StartTime   string    `gorm:"type:varchar(8);not null" json:"start_time"`
	EndTime     string    `gorm:"type:varchar(8);not null" json:"end_time"`
	IsAvailable bool      `gorm:"type:boolean;default:true;index:idx_schedules_available" json:"is_available"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName 指定表名
func (Schedule) TableName() string {
	return "schedules"
}

// BeforeCreate 创建前生成 UUID
func (s *Schedule) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}
