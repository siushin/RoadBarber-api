package service

import (
	"time"

	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"
)

type AdminBookingService struct{}

// AdminBookingItem 管理员端预约（含关联快照）
type AdminBookingItem struct {
	models.Booking
	CustomerPhone string `json:"customer_phone"`
	CustomerName  string `json:"customer_name"`
	MerchantTitle string `json:"merchant_title"`
	ServiceName   string `json:"service_name"`
	ShopName      string `json:"shop_name"`
}

// AdminListBookingsRequest 管理员预约列表查询参数
type AdminListBookingsRequest struct {
	Status   int8  `query:"status"`
	Keyword  string `query:"keyword"`
	Page     int   `query:"page"`
	PageSize int   `query:"page_size"`
}

// ListBookings 管理员预约列表（含关联信息）
func (s *AdminBookingService) ListBookings(req *AdminListBookingsRequest) ([]AdminBookingItem, int64, error) {
	db := config.GetDB().Table("bookings").
		Select(`bookings.*,
			users.phone AS customer_phone,
			users.nickname AS customer_name,
			merchants.title AS merchant_title,
			services.name AS service_name,
			shops.name AS shop_name`).
		Joins("LEFT JOIN users ON users.id = bookings.customer_id").
		Joins("LEFT JOIN merchants ON merchants.id = bookings.merchant_id").
		Joins("LEFT JOIN services ON services.id = bookings.service_id").
		Joins("LEFT JOIN shops ON shops.id = bookings.shop_id")

	if req.Status > 0 {
		db = db.Where("bookings.status = ?", req.Status)
	}
	if req.Keyword != "" {
		db = db.Where("users.phone LIKE ? OR users.nickname LIKE ? OR services.name LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
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

	type Result struct {
		models.Booking
		CustomerPhone string
		CustomerName  string
		MerchantTitle string
		ServiceName   string
		ShopName      string
	}
	var rows []Result
	if err := db.Offset(offset).Limit(pageSize).Order("bookings.created_at DESC").Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]AdminBookingItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, AdminBookingItem{
			Booking:       r.Booking,
			CustomerPhone: r.CustomerPhone,
			CustomerName:  r.CustomerName,
			MerchantTitle: r.MerchantTitle,
			ServiceName:   r.ServiceName,
			ShopName:      r.ShopName,
		})
	}
	return items, total, nil
}

// MonthRevenue 本月已完成预约营收
func MonthRevenue() (float64, error) {
	var total float64
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	err := config.GetDB().Model(&models.Booking{}).
		Where("status = ? AND finish_time >= ?", models.BookingStatusCompleted, firstOfMonth).
		Select("COALESCE(SUM(price), 0)").
		Scan(&total).Error
	return total, err
}