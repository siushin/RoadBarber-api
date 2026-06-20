package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MerchantProfile 商家扩展表：存储商家资质信息
type MerchantProfile struct {
	ID                 string     `gorm:"type:uuid;primaryKey" json:"id"`
	UserID             string     `gorm:"type:uuid;not null;uniqueIndex:idx_merchant_profiles_user" json:"user_id"`
	MerchantType       int8       `gorm:"type:smallint;not null" json:"merchant_type"`
	IDCard             string     `gorm:"type:varchar(20)" json:"id_card"`
	IDCardFront        string     `gorm:"type:varchar(500)" json:"id_card_front"`
	IDCardBack         string     `gorm:"type:varchar(500)" json:"id_card_back"`
	BusinessLicense    string     `gorm:"type:varchar(500)" json:"business_license"`
	CompanyName        string     `gorm:"type:varchar(200)" json:"company_name"`
	TaxNumber          string     `gorm:"type:varchar(50)" json:"tax_number"`
	QualificationDocs  string     `gorm:"type:jsonb" json:"qualification_docs"`
	AuditStatus        int8       `gorm:"type:smallint;not null;default:2;index:idx_merchant_profiles_audit" json:"audit_status"`
	AuditRemark        string     `gorm:"type:varchar(500)" json:"audit_remark"`
	AuditTime          *time.Time `json:"audit_time"`
	AuditorID          string     `gorm:"type:uuid" json:"auditor_id"`
	CreatedAt          time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName 指定表名
func (MerchantProfile) TableName() string {
	return "merchant_profiles"
}

// BeforeCreate 创建前生成 UUID
func (m *MerchantProfile) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

// 商家类型常量
const (
	MerchantTypeIndividual = 1 // 个人
	MerchantTypeSelfEmployed = 2 // 个体户
	MerchantTypeCompany    = 3 // 公司
)

// 商家审核状态常量
const (
	AuditStatusApproved = 1 // 通过
	AuditStatusPending  = 2 // 待审核
	AuditStatusRejected = 3 // 拒绝
)
