package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MerchantApply 商家入驻申请表：商家提交入驻申请
type MerchantApply struct {
	ID              string     `gorm:"type:uuid;primaryKey" json:"id"`
	ApplicantName   string     `gorm:"type:varchar(100);not null" json:"applicant_name"`
	ApplicantPhone  string     `gorm:"type:varchar(20);not null;index:idx_merchant_applies_phone" json:"applicant_phone"`
	ApplicantType   int8       `gorm:"type:smallint;not null" json:"applicant_type"`
	IDCard          string     `gorm:"type:varchar(20)" json:"id_card"`
	CompanyName     string     `gorm:"type:varchar(200)" json:"company_name"`
	BusinessLicense string     `gorm:"type:varchar(500)" json:"business_license"`
	LocationID      *string    `gorm:"type:uuid" json:"location_id"`
	Address         string     `gorm:"type:varchar(255)" json:"address"`
	Longitude       float64    `gorm:"type:decimal(10,7)" json:"longitude"`
	Latitude        float64    `gorm:"type:decimal(10,7)" json:"latitude"`
	Status          int8       `gorm:"type:smallint;not null;default:2;index:idx_merchant_applies_status" json:"status"`
	RejectReason    string     `gorm:"type:varchar(500)" json:"reject_reason"`
	AuditTime       *time.Time `json:"audit_time"`
	AuditorID       string     `gorm:"type:uuid" json:"auditor_id"`
	CreatedAt       time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName 指定表名
func (MerchantApply) TableName() string {
	return "merchant_applies"
}

// BeforeCreate 创建前生成 UUID
func (m *MerchantApply) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
