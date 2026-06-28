package service

import (
	"errors"
	"fmt"
	"time"

	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingService struct{}

// CreateBookingRequest 创建预约请求
type CreateBookingRequest struct {
	MerchantID string `json:"merchant_id" validate:"required"`
	ServiceID  string `json:"service_id" validate:"required"`
	ScheduleID string `json:"schedule_id" validate:"required"`
	Remark     string `json:"remark"`
}

// CreateBooking 创建预约
func (s *BookingService) CreateBooking(customerID string, req *CreateBookingRequest) (*models.Booking, error) {
	// 校验排班是否可预约
	var schedule models.Schedule
	err := config.GetDB().Where("id = ? AND is_available = ?", req.ScheduleID, true).First(&schedule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("该时段不可预约")
		}
		return nil, fmt.Errorf("查询排班失败: %w", err)
	}

	// 查询服务
	var service models.Service
	err = config.GetDB().Where("id = ?", req.ServiceID).First(&service).Error
	if err != nil {
		return nil, errors.New("服务项目不存在")
	}

	// 查询商家个性化价格
	var ms models.MerchantService
	customPrice := service.Price
	err = config.GetDB().Where("merchant_id = ? AND service_id = ?", req.MerchantID, req.ServiceID).First(&ms).Error
	if err == nil && ms.Price != nil {
		customPrice = *ms.Price
	}

	// 生成订单号
	orderNo := fmt.Sprintf("BK%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8])

	// 事务：标记排班为不可用 + 创建预约
	var booking models.Booking
	err = config.GetDB().Transaction(func(tx *gorm.DB) error {
		// 锁定排班
		var lockedSchedule models.Schedule
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ? AND is_available = ?", req.ScheduleID, true).First(&lockedSchedule).Error; err != nil {
			return errors.New("该时段已被预约或不可用")
		}

		// 查询商家所属店铺
		var merchant models.Merchant
		if err := tx.Where("id = ?", req.MerchantID).First(&merchant).Error; err != nil {
			return errors.New("商家不存在")
		}

		// 标记排班为不可用
		if err := tx.Model(&lockedSchedule).Update("is_available", false).Error; err != nil {
			return fmt.Errorf("更新排班失败: %w", err)
		}

		// 创建预约
		booking = models.Booking{
			OrderNo:         orderNo,
			CustomerID:      customerID,
			MerchantID:      req.MerchantID,
			ShopID:          merchant.ShopID,
			ServiceID:       req.ServiceID,
			ScheduleID:      req.ScheduleID,
			AppointmentDate: schedule.WorkDate,
			AppointmentTime: schedule.StartTime,
			Duration:        service.Duration,
			Price:           customPrice,
			Status:          models.BookingStatusPending,
			Remark:          req.Remark,
		}

		if err := tx.Create(&booking).Error; err != nil {
			return fmt.Errorf("创建预约失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &booking, nil
}

// ListMyBookingsRequest 顾客预约列表
type ListMyBookingsRequest struct {
	CustomerID string `query:"-"`
	Status     int8   `query:"status"`
	Page       int    `query:"page"`
	PageSize   int    `query:"page_size"`
}

// ListMyBookings 我的预约
func (s *BookingService) ListMyBookings(req *ListMyBookingsRequest) ([]models.Booking, int64, error) {
	db := config.GetDB().Model(&models.Booking{}).Where("customer_id = ?", req.CustomerID)
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
	err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&bookings).Error
	return bookings, total, err
}

// GetBookingDetail 预约详情
func (s *BookingService) GetBookingDetail(id, userID string) (*models.Booking, error) {
	var booking models.Booking
	err := config.GetDB().Where("id = ?", id).First(&booking).Error
	if err != nil {
		return nil, errors.New("预约不存在")
	}
	if booking.CustomerID != userID {
		return nil, errors.New("无权查看")
	}
	return &booking, nil
}

// CancelBooking 取消预约
func (s *BookingService) CancelBooking(id, userID, reason string) error {
	var booking models.Booking
	err := config.GetDB().Where("id = ?", id).First(&booking).Error
	if err != nil {
		return errors.New("预约不存在")
	}
	if booking.CustomerID != userID {
		return errors.New("无权操作")
	}
	if booking.Status != models.BookingStatusPending && booking.Status != models.BookingStatusConfirmed {
		return errors.New("该预约不可取消")
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
			"status":        models.BookingStatusCancelled,
			"cancel_reason": reason,
			"cancel_time":   now,
		}).Error
	})
}
