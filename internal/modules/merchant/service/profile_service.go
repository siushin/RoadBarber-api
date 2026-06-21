package service

import (
	"errors"

	"roadbarber/backend/internal/config"
	"roadbarber/backend/internal/models"
)

type ProfileService struct{}

// ProfileResponse 商家 profile（含 user 基础信息）
type ProfileResponse struct {
	User        models.User          `json:"user"`
	Merchant    models.Merchant      `json:"merchant"`
	Profile     *models.MerchantProfile `json:"profile,omitempty"`
}

// UpdateProfileRequest 更新 profile 请求
type UpdateProfileRequest struct {
	Title           *string `json:"title"`
	Specialties     *string `json:"specialties"`
	ExperienceYears *int    `json:"experience_years"`
	Introduction    *string `json:"introduction"`
	Avatar          *string `json:"avatar"`
}

// GetProfile 获取商家 profile（含 user/merchant/profile）
func (s *ProfileService) GetProfile(userID string) (*ProfileResponse, error) {
	var user models.User
	if err := config.GetDB().Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	var merchant models.Merchant
	if err := config.GetDB().Where("user_id = ?", userID).First(&merchant).Error; err != nil {
		return nil, errors.New("商家信息不存在")
	}

	var profile models.MerchantProfile
	if err := config.GetDB().Where("user_id = ?", userID).First(&profile).Error; err == nil {
		return &ProfileResponse{User: user, Merchant: merchant, Profile: &profile}, nil
	}
	return &ProfileResponse{User: user, Merchant: merchant}, nil
}

// UpdateProfile 更新商家基础信息
func (s *ProfileService) UpdateProfile(userID string, req *UpdateProfileRequest) error {
	var merchant models.Merchant
	if err := config.GetDB().Where("user_id = ?", userID).First(&merchant).Error; err != nil {
		return errors.New("商家信息不存在")
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Specialties != nil {
		updates["specialties"] = *req.Specialties
	}
	if req.ExperienceYears != nil {
		updates["experience_years"] = *req.ExperienceYears
	}
	if req.Introduction != nil {
		updates["introduction"] = *req.Introduction
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}

	if len(updates) == 0 {
		return nil
	}
	return config.GetDB().Model(&merchant).Updates(updates).Error
}