package service

import (
	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"
)

type CustomerService struct{}

// MerchantCustomerItem 商家端顾客（来自预约去重）
type MerchantCustomerItem struct {
	models.User
	LastBookingAt *string `json:"last_booking_at,omitempty"`
	BookingCount  int     `json:"booking_count"`
}

// ListCustomersByMerchantUser 列出该商家用户的所有顾客（按 user_id）
func (s *CustomerService) ListCustomersByMerchantUser(merchantUserID string) ([]MerchantCustomerItem, error) {
	// 找到该 user 对应的 merchant.id
	var merchant models.Merchant
	if err := config.GetDB().Where("user_id = ?", merchantUserID).First(&merchant).Error; err != nil {
		return []MerchantCustomerItem{}, nil
	}

	type row struct {
		models.User
		LastBookingAt *string
		BookingCount  int
	}
	var rows []row
	err := config.GetDB().Table("users").
		Select(`users.*,
			MAX(bookings.created_at) AS last_booking_at,
			COUNT(bookings.id) AS booking_count`).
		Joins("INNER JOIN bookings ON bookings.customer_id = users.id").
		Where("bookings.merchant_id = ?", merchant.ID).
		Group("users.id").
		Order("last_booking_at DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]MerchantCustomerItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, MerchantCustomerItem{
			User:          r.User,
			LastBookingAt: r.LastBookingAt,
			BookingCount:  r.BookingCount,
		})
	}
	return out, nil
}