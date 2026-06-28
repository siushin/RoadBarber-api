package service

import (
	"errors"
	"time"

	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"
)

type MerchantDashboardService struct{}

// MerchantDashboardStats 商家仪表盘统计
type MerchantDashboardStats struct {
	TodayBookings    int64   `json:"today_bookings"`
	PendingBookings  int64   `json:"pending_bookings"`
	ConfirmedBookings int64  `json:"confirmed_bookings"`
	CompletedBookings int64  `json:"completed_bookings"`
	MonthRevenue     float64 `json:"month_revenue"`
	AvgRating        float64 `json:"avg_rating"`
	ReviewCount      int     `json:"review_count"`
	ServiceCount     int     `json:"service_count"`
}

// GetStats 商家仪表盘
// 管理员（role=3）无商家数据，返回全 0 的统计
func (s *MerchantDashboardService) GetStats(userID string) (*MerchantDashboardStats, error) {
	if isAdminUser(userID) {
		return &MerchantDashboardStats{}, nil
	}

	var merchant models.Merchant
	if err := config.GetDB().Where("user_id = ?", userID).First(&merchant).Error; err != nil {
		return nil, errors.New("商家信息不存在")
	}

	db := config.GetDB()
	today := time.Now().Format("2006-01-02")
	firstOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())

	var stats MerchantDashboardStats
	stats.AvgRating = merchant.Rating
	stats.ReviewCount = merchant.ReviewCount
	stats.ServiceCount = merchant.ServiceCount

	db.Model(&models.Booking{}).
		Where("merchant_id = ? AND DATE(appointment_date) = ?", merchant.ID, today).
		Count(&stats.TodayBookings)
	db.Model(&models.Booking{}).
		Where("merchant_id = ? AND status = ?", merchant.ID, models.BookingStatusPending).
		Count(&stats.PendingBookings)
	db.Model(&models.Booking{}).
		Where("merchant_id = ? AND status = ?", merchant.ID, models.BookingStatusConfirmed).
		Count(&stats.ConfirmedBookings)
	db.Model(&models.Booking{}).
		Where("merchant_id = ? AND status = ?", merchant.ID, models.BookingStatusCompleted).
		Count(&stats.CompletedBookings)

	db.Model(&models.Booking{}).
		Where("merchant_id = ? AND status = ? AND finish_time >= ?", merchant.ID, models.BookingStatusCompleted, firstOfMonth).
		Select("COALESCE(SUM(price), 0)").
		Scan(&stats.MonthRevenue)

	return &stats, nil
}