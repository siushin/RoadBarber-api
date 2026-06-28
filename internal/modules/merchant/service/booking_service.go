package service

import (
	"errors"
	"time"

	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"

	"gorm.io/gorm"
)

type BookingService struct{}

// ListMerchantBookingsRequest 商家预约列表
type ListMerchantBookingsRequest struct {
	MerchantID string `query:"-"`
	Status     int8   `query:"status"`
	Page       int    `query:"page"`
	PageSize   int    `query:"page_size"`
}

// ListMerchantBookings 商家预约列表
func (s *BookingService) ListMerchantBookings(req *ListMerchantBookingsRequest) ([]models.Booking, int64, error) {
	db := config.GetDB().Model(&models.Booking{}).Where("merchant_id = ?", req.MerchantID)
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

	var bookings []models.Booking
	err := db.Offset(offset).Limit(pageSize).Order("appointment_date DESC, appointment_time DESC").Find(&bookings).Error
	return bookings, total, err
}

// ConfirmBooking 商家确认预约
func (s *BookingService) ConfirmBooking(id, merchantID string) error {
	var booking models.Booking
	err := config.GetDB().Where("id = ?", id).First(&booking).Error
	if err != nil {
		return errors.New("预约不存在")
	}
	if booking.MerchantID != merchantID {
		return errors.New("无权操作")
	}
	if booking.Status != models.BookingStatusPending {
		return errors.New("该预约已处理")
	}

	now := time.Now()
	return config.GetDB().Model(&booking).Updates(map[string]interface{}{
		"status":       models.BookingStatusConfirmed,
		"confirm_time": now,
	}).Error
}

// RejectBooking 商家拒绝预约
func (s *BookingService) RejectBooking(id, merchantID, reason string) error {
	var booking models.Booking
	err := config.GetDB().Where("id = ?", id).First(&booking).Error
	if err != nil {
		return errors.New("预约不存在")
	}
	if booking.MerchantID != merchantID {
		return errors.New("无权操作")
	}
	if booking.Status != models.BookingStatusPending {
		return errors.New("该预约已处理")
	}

	now := time.Now()
	return config.GetDB().Transaction(func(tx *gorm.DB) error {
		// 恢复排班可用
		if err := tx.Model(&models.Schedule{}).Where("id = ?", booking.ScheduleID).
			Update("is_available", true).Error; err != nil {
			return err
		}
		// 更新预约状态
		return tx.Model(&booking).Updates(map[string]interface{}{
			"status":        models.BookingStatusRejected,
			"cancel_reason": reason,
			"cancel_time":   now,
		}).Error
	})
}

// StartService 开始服务
func (s *BookingService) StartService(id, merchantID string) error {
	var booking models.Booking
	err := config.GetDB().Where("id = ?", id).First(&booking).Error
	if err != nil {
		return errors.New("预约不存在")
	}
	if booking.MerchantID != merchantID {
		return errors.New("无权操作")
	}
	if booking.Status != models.BookingStatusConfirmed {
		return errors.New("该预约未确认，无法开始服务")
	}

	now := time.Now()
	return config.GetDB().Model(&booking).Updates(map[string]interface{}{
		"status":     models.BookingStatusServing,
		"start_time": now,
	}).Error
}

// FinishService 完成服务
func (s *BookingService) FinishService(id, merchantID string) error {
	var booking models.Booking
	err := config.GetDB().Where("id = ?", id).First(&booking).Error
	if err != nil {
		return errors.New("预约不存在")
	}
	if booking.MerchantID != merchantID {
		return errors.New("无权操作")
	}
	if booking.Status != models.BookingStatusServing {
		return errors.New("该预约未开始服务")
	}

	now := time.Now()
	return config.GetDB().Transaction(func(tx *gorm.DB) error {
		// 更新预约状态
		if err := tx.Model(&booking).Updates(map[string]interface{}{
			"status":      models.BookingStatusCompleted,
			"finish_time": now,
		}).Error; err != nil {
			return err
		}
		// 增加服务次数
		return tx.Model(&models.Merchant{}).Where("id = ?", booking.MerchantID).
			UpdateColumn("service_count", tx.Raw("service_count + 1")).Error
	})
}
