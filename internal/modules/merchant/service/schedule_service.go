package service

import (
	"errors"
	"time"

	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"
)

type ScheduleService struct{}

// CreateScheduleRequest 发布排班
type CreateScheduleRequest struct {
	WorkDate  string `json:"work_date" validate:"required"`  // 2025-01-01
	StartTime string `json:"start_time" validate:"required"` // 09:00
	EndTime   string `json:"end_time" validate:"required"`   // 18:00
}

// BatchCreateRequest 批量发布排班
type BatchCreateRequest struct {
	Schedules []CreateScheduleRequest `json:"schedules" validate:"required"`
}

// CreateSchedule 发布排班
func (s *ScheduleService) CreateSchedule(merchantID string, req *CreateScheduleRequest) (*models.Schedule, error) {
	workDate, err := time.Parse("2006-01-02", req.WorkDate)
	if err != nil {
		return nil, errors.New("日期格式错误，应为 YYYY-MM-DD")
	}

	schedule := models.Schedule{
		MerchantID:  merchantID,
		WorkDate:    workDate,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsAvailable: true,
	}

	if err := config.GetDB().Create(&schedule).Error; err != nil {
		return nil, err
	}

	return &schedule, nil
}

// BatchCreate 批量发布
func (s *ScheduleService) BatchCreate(merchantID string, req *BatchCreateRequest) error {
	schedules := make([]models.Schedule, 0, len(req.Schedules))
	for _, item := range req.Schedules {
		workDate, err := time.Parse("2006-01-02", item.WorkDate)
		if err != nil {
			return errors.New("日期格式错误，应为 YYYY-MM-DD")
		}
		schedules = append(schedules, models.Schedule{
			MerchantID:  merchantID,
			WorkDate:    workDate,
			StartTime:   item.StartTime,
			EndTime:     item.EndTime,
			IsAvailable: true,
		})
	}
	return config.GetDB().Create(&schedules).Error
}

// ListMySchedules 我的排班（商家端）
func (s *ScheduleService) ListMySchedules(merchantID, startDate, endDate string) ([]models.Schedule, error) {
	db := config.GetDB().Where("merchant_id = ?", merchantID)

	if startDate != "" {
		db = db.Where("work_date >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("work_date <= ?", endDate)
	}

	var schedules []models.Schedule
	err := db.Order("work_date ASC, start_time ASC").Find(&schedules).Error
	return schedules, err
}

// DeleteSchedule 删除排班
func (s *ScheduleService) DeleteSchedule(id, merchantID string) error {
	var schedule models.Schedule
	err := config.GetDB().Where("id = ? AND merchant_id = ?", id, merchantID).First(&schedule).Error
	if err != nil {
		return errors.New("排班不存在")
	}
	if !schedule.IsAvailable {
		return errors.New("该时段已被预约，无法删除")
	}
	return config.GetDB().Delete(&schedule).Error
}
