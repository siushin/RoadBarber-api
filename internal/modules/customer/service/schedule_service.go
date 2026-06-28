package service

import (
	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"
)

type ScheduleService struct{}

// ListAvailableSchedules 查询商家可用排班
func (s *ScheduleService) ListAvailableSchedules(merchantID, date string) ([]models.Schedule, error) {
	var schedules []models.Schedule
	db := config.GetDB().Where("merchant_id = ? AND is_available = ?", merchantID, true)

	if date != "" {
		db = db.Where("work_date = ?", date)
	}

	err := db.Order("work_date ASC, start_time ASC").Find(&schedules).Error
	return schedules, err
}
