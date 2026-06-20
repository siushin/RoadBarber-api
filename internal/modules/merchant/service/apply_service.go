package service

import (
	"errors"
	"fmt"

	"roadbarber/backend/internal/config"
	"roadbarber/backend/internal/models"
)

type MerchantApplyService struct{}

// CreateApplyRequest 提交入驻申请
type CreateApplyRequest struct {
	ApplicantName   string  `json:"applicant_name" validate:"required"`
	ApplicantPhone  string  `json:"applicant_phone" validate:"required"`
	ApplicantType   int8    `json:"applicant_type" validate:"required"`
	IDCard          string  `json:"id_card"`
	CompanyName     string  `json:"company_name"`
	BusinessLicense string  `json:"business_license"`
	LocationID      string  `json:"location_id"`
	Address         string  `json:"address"`
	Longitude       float64 `json:"longitude"`
	Latitude        float64 `json:"latitude"`
}

// SubmitApply 提交入驻申请
func (s *MerchantApplyService) SubmitApply(req *CreateApplyRequest) (*models.MerchantApply, error) {
	if req.ApplicantName == "" || req.ApplicantPhone == "" {
		return nil, errors.New("姓名和手机号必填")
	}

	// 检查是否已有待审核申请
	var count int64
	config.GetDB().Model(&models.MerchantApply{}).
		Where("applicant_phone = ? AND status = ?", req.ApplicantPhone, models.AuditStatusPending).
		Count(&count)
	if count > 0 {
		return nil, errors.New("您有待审核的申请，请勿重复提交")
	}

	apply := models.MerchantApply{
		ApplicantName:   req.ApplicantName,
		ApplicantPhone:  req.ApplicantPhone,
		ApplicantType:   req.ApplicantType,
		IDCard:          req.IDCard,
		CompanyName:     req.CompanyName,
		BusinessLicense: req.BusinessLicense,
		Address:         req.Address,
		Longitude:       req.Longitude,
		Latitude:        req.Latitude,
		Status:          models.AuditStatusPending,
	}
	if req.LocationID != "" {
		apply.LocationID = &req.LocationID
	}

	if err := config.GetDB().Create(&apply).Error; err != nil {
		return nil, fmt.Errorf("提交申请失败: %w", err)
	}

	return &apply, nil
}

// ListMyApplies 我提交的申请记录
func (s *MerchantApplyService) ListMyApplies(phone string) ([]models.MerchantApply, error) {
	var applies []models.MerchantApply
	err := config.GetDB().Where("applicant_phone = ?", phone).
		Order("created_at DESC").Find(&applies).Error
	return applies, err
}

// ListAppliesRequest 申请列表查询（管理员端）
type ListAppliesRequest struct {
	Status   int8  `query:"status"`
	Page     int   `query:"page"`
	PageSize int   `query:"page_size"`
}

// ListApplies 申请列表（管理员端）
func (s *MerchantApplyService) ListApplies(req *ListAppliesRequest) ([]models.MerchantApply, int64, error) {
	db := config.GetDB().Model(&models.MerchantApply{})
	if req.Status > 0 {
		db = db.Where("status = ?", req.Status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var applies []models.MerchantApply
	err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&applies).Error
	return applies, total, err
}

// GetApplyDetail 申请详情
func (s *MerchantApplyService) GetApplyDetail(id string) (*models.MerchantApply, error) {
	var apply models.MerchantApply
	err := config.GetDB().Where("id = ?", id).First(&apply).Error
	if err != nil {
		return nil, errors.New("申请不存在")
	}
	return &apply, nil
}

// ApproveApply 审核通过
func (s *MerchantApplyService) ApproveApply(id, auditorID string) error {
	now := models.MerchantApply{}
	_ = now
	return config.GetDB().Model(&models.MerchantApply{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     models.AuditStatusApproved,
		"audit_time": now.CreatedAt, // 简化处理
		"auditor_id": auditorID,
	}).Error
}

// RejectApply 审核拒绝
func (s *MerchantApplyService) RejectApply(id, auditorID, reason string) error {
	return config.GetDB().Model(&models.MerchantApply{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":         models.AuditStatusRejected,
		"auditor_id":     auditorID,
		"reject_reason":  reason,
	}).Error
}
