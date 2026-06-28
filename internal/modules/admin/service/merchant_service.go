package service

import (
	"errors"
	"time"

	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"

	"gorm.io/gorm"
)

type AdminService struct{}

// TrendPoint 趋势图数据点
type TrendPoint struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// TopMerchant Top 商家
type TopMerchant struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Avatar      string  `json:"avatar"`
	Rating      float64 `json:"rating"`
	ReviewCount int     `json:"review_count"`
}

// StatusCount 状态分布
type StatusCount struct {
	Status int   `json:"status"`
	Count  int64 `json:"count"`
}

// DashboardStats 仪表盘统计
type DashboardStats struct {
	TotalCustomers   int64         `json:"total_customers"`
	TotalMerchants   int64         `json:"total_merchants"`
	TotalShops       int64         `json:"total_shops"`
	TodayBookings    int64         `json:"today_bookings"`
	PendingApplies   int64         `json:"pending_applies"`
	PendingMerchants int64         `json:"pending_merchants"`
	TotalBookings    int64         `json:"total_bookings"`
	TotalReviews     int64         `json:"total_reviews"`
	MonthRevenue     float64       `json:"month_revenue"`
	BookingTrend     []TrendPoint  `json:"booking_trend"`
	TopMerchants     []TopMerchant `json:"top_merchants"`
	StatusBreakdown  []StatusCount `json:"status_breakdown"`
}

// GetDashboardStats 仪表盘统计
func (s *AdminService) GetDashboardStats() (*DashboardStats, error) {
	var stats DashboardStats
	db := config.GetDB()

	// 顾客数
	db.Model(&models.User{}).Where("role = ?", models.RoleCustomer).Count(&stats.TotalCustomers)
	// 商家数
	db.Model(&models.Merchant{}).Where("status = ?", models.MerchantStatusNormal).Count(&stats.TotalMerchants)
	// 店铺数
	db.Model(&models.Shop{}).Where("status = ?", models.ShopStatusNormal).Count(&stats.TotalShops)
	// 今日预约
	db.Model(&models.Booking{}).Where("DATE(appointment_date) = ?", time.Now().Format("2006-01-02")).Count(&stats.TodayBookings)
	// 待审核入驻申请
	db.Model(&models.MerchantApply{}).Where("status = ?", models.AuditStatusPending).Count(&stats.PendingApplies)
	// 待审核商家
	db.Model(&models.MerchantProfile{}).Where("audit_status = ?", models.AuditStatusPending).Count(&stats.PendingMerchants)
	// 总预约数
	db.Model(&models.Booking{}).Count(&stats.TotalBookings)
	// 总评价数
	db.Model(&models.Review{}).Count(&stats.TotalReviews)

	// 本月营收
	monthRevenue, _ := MonthRevenue()
	stats.MonthRevenue = monthRevenue

	// 近 30 日预约趋势
	stats.BookingTrend = buildBookingTrend(db)

	// Top 5 商家（按 review_count 排）
	stats.TopMerchants = buildTopMerchants(db)

	// 预约状态分布
	stats.StatusBreakdown = buildStatusBreakdown(db)

	return &stats, nil
}

func buildBookingTrend(db *gorm.DB) []TrendPoint {
	var points []TrendPoint
	for i := 29; i >= 0; i-- {
		day := time.Now().AddDate(0, 0, -i)
		var count int64
		db.Model(&models.Booking{}).
			Where("DATE(created_at) = ?", day.Format("2006-01-02")).
			Count(&count)
		points = append(points, TrendPoint{
			Date:  day.Format("01-02"),
			Count: count,
		})
	}
	return points
}

func buildTopMerchants(db *gorm.DB) []TopMerchant {
	type row struct {
		ID          string
		Title       string
		Avatar      string
		Rating      float64
		ReviewCount int
	}
	var rows []row
	db.Table("merchants").
		Select("merchants.id, merchants.title, merchants.avatar, merchants.rating, merchants.review_count").
		Where("merchants.status = ?", models.MerchantStatusNormal).
		Order("merchants.review_count DESC, merchants.rating DESC").
		Limit(5).
		Scan(&rows)

	out := make([]TopMerchant, 0, len(rows))
	for _, r := range rows {
		out = append(out, TopMerchant{
			ID:          r.ID,
			Title:       r.Title,
			Avatar:      r.Avatar,
			Rating:      r.Rating,
			ReviewCount: r.ReviewCount,
		})
	}
	return out
}

func buildStatusBreakdown(db *gorm.DB) []StatusCount {
	type row struct {
		Status int8
		Count  int64
	}
	var rows []row
	db.Model(&models.Booking{}).
		Select("status, COUNT(*) AS count").
		Group("status").
		Scan(&rows)
	out := make([]StatusCount, 0, len(rows))
	for _, r := range rows {
		out = append(out, StatusCount{Status: int(r.Status), Count: r.Count})
	}
	return out
}

// ListMerchantsRequest 商家列表
type ListMerchantsRequest struct {
	AuditStatus int8  `query:"audit_status"`
	Keyword     string `query:"keyword"`
	Page        int   `query:"page"`
	PageSize    int   `query:"page_size"`
}

// MerchantWithProfile 商家+用户+资质
type MerchantWithProfile struct {
	models.Merchant
	Phone     string                  `json:"phone"`
	Nickname  string                  `json:"nickname"`
	Profile   *models.MerchantProfile `json:"profile,omitempty"`
}

// ListMerchants 商家列表（管理员端）
func (s *AdminService) ListMerchants(req *ListMerchantsRequest) ([]MerchantWithProfile, int64, error) {
	db := config.GetDB().Table("merchants").
		Select("merchants.*, users.phone, users.nickname").
		Joins("LEFT JOIN users ON users.id = merchants.user_id")

	if req.AuditStatus > 0 {
		db = db.Joins("LEFT JOIN merchant_profiles ON merchant_profiles.user_id = merchants.user_id").
			Where("merchant_profiles.audit_status = ?", req.AuditStatus)
	}
	if req.Keyword != "" {
		db = db.Where("users.phone LIKE ? OR users.nickname LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%")
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
		models.Merchant
		Phone    string
		Nickname string
	}
	var results []Result
	err := db.Offset(offset).Limit(pageSize).
		Order("merchants.created_at DESC").
		Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	// 加载资质信息
	merchants := make([]MerchantWithProfile, 0, len(results))
	for _, r := range results {
		var profile models.MerchantProfile
		if err := config.GetDB().Where("user_id = ?", r.UserID).First(&profile).Error; err == nil {
			merchants = append(merchants, MerchantWithProfile{
				Merchant: r.Merchant,
				Phone:    r.Phone,
				Nickname: r.Nickname,
				Profile:  &profile,
			})
		} else {
			merchants = append(merchants, MerchantWithProfile{
				Merchant: r.Merchant,
				Phone:    r.Phone,
				Nickname: r.Nickname,
			})
		}
	}

	return merchants, total, nil
}

// VerifyMerchant 审核商家资质
func (s *AdminService) VerifyMerchant(merchantID, auditorID string, auditStatus int8, remark string) error {
	return config.GetDB().Transaction(func(tx *gorm.DB) error {
		// 找到商家对应的 user_id
		var merchant models.Merchant
		if err := tx.Where("id = ?", merchantID).First(&merchant).Error; err != nil {
			return errors.New("商家不存在")
		}

		// 更新资质审核状态
		if err := tx.Model(&models.MerchantProfile{}).Where("user_id = ?", merchant.UserID).
			Updates(map[string]interface{}{
				"audit_status": auditStatus,
				"audit_remark": remark,
				"audit_time":   time.Now(),
				"auditor_id":   auditorID,
			}).Error; err != nil {
			return err
		}

		// 同步更新商家的 is_verified
		isVerified := auditStatus == models.AuditStatusApproved
		return tx.Model(&merchant).Update("is_verified", isVerified).Error
	})
}

// ListAppliesRequest 申请列表
type ListAppliesRequest struct {
	Status   int8  `query:"status"`
	Page     int   `query:"page"`
	PageSize int   `query:"page_size"`
}

// ListApplies 申请列表（管理员端）
func (s *AdminService) ListApplies(req *ListAppliesRequest) ([]models.MerchantApply, int64, error) {
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

// ApproveApply 审核通过申请
func (s *AdminService) ApproveApply(id, auditorID string) error {
	// TODO: 完整流程需要创建 user + merchant + merchant_profile
	return config.GetDB().Model(&models.MerchantApply{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     models.AuditStatusApproved,
		"auditor_id": auditorID,
		"audit_time": time.Now(),
	}).Error
}

// RejectApply 审核拒绝申请
func (s *AdminService) RejectApply(id, auditorID, reason string) error {
	return config.GetDB().Model(&models.MerchantApply{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":         models.AuditStatusRejected,
		"auditor_id":     auditorID,
		"reject_reason":  reason,
		"audit_time":     time.Now(),
	}).Error
}

var _ = gorm.ErrRecordNotFound
