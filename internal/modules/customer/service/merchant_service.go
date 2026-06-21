package service

import (
	"strings"
	"time"

	"roadbarber/backend/internal/config"
	"roadbarber/backend/internal/models"
)

type MerchantService struct{}

// ListMerchantsRequest 商家列表查询参数
type ListMerchantsRequest struct {
	ShopID    string  `query:"shop_id"`
	Category  string  `query:"category"`
	Keyword   string  `query:"keyword"`
	SortBy    string  `query:"sort_by"` // rating, distance, service_count
	Latitude  float64 `query:"latitude"`
	Longitude float64 `query:"longitude"`
	Date      string  `query:"date"` // YYYY-MM-DD，默认今天
	Page      int     `query:"page"`
	PageSize  int     `query:"page_size"`
}

// ListMerchants 商家列表（首页对接：含 start_price / business_hours / distance / available_slots）
func (s *MerchantService) ListMerchants(req *ListMerchantsRequest) ([]models.Merchant, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	workDate := req.Date
	if workDate == "" {
		workDate = time.Now().Format("2006-01-02")
	}

	// WHERE 子句
	whereParts := []string{"m.status = ?"}
	whereArgs := []interface{}{models.MerchantStatusNormal}

	// 过滤掉当日无排班的商家（首页需求：没有排班的不展示）
	whereParts = append(whereParts, `(SELECT COUNT(*) FROM schedules
		WHERE merchant_id = m.id AND work_date = ? AND is_available = true) > 0`)
	whereArgs = append(whereArgs, workDate)

	if req.ShopID != "" {
		whereParts = append(whereParts, "m.shop_id = ?")
		whereArgs = append(whereArgs, req.ShopID)
	}
	if req.Keyword != "" {
		whereParts = append(whereParts, "(m.title LIKE ? OR m.introduction LIKE ?)")
		whereArgs = append(whereArgs, "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 计算字段子查询
	// start_price: 服务 join 后 min(price)，无服务时为 0
	// available_slots: 当日 is_available=true 的排班数
	// business_hours: 当日排班 min(start_time) - max(end_time)，无排班为 NULL
	selectExtras := `
		COALESCE((SELECT MIN(s.price) FROM services s
		 JOIN merchant_services ms ON s.id = ms.service_id
		 WHERE ms.merchant_id = m.id AND s.status = ?), 0) AS start_price,
		COALESCE((SELECT COUNT(*) FROM schedules
		 WHERE merchant_id = m.id AND work_date = ? AND is_available = true), 0) AS available_slots,
		(SELECT MIN(start_time) || ' - ' || MAX(end_time) FROM schedules
		 WHERE merchant_id = m.id AND work_date = ? AND is_available = true) AS business_hours`
	selectArgs := []interface{}{models.ServiceStatusOnSale, workDate, workDate}

	// 距离：Haversine 球面公式，仅当用户传 lat/lng 才计算
	hasDistance := req.Latitude != 0 || req.Longitude != 0
	if hasDistance {
		selectExtras += `,
		(6371 * acos(
			cos(radians(?)) * cos(radians(m.latitude)) *
			cos(radians(m.longitude) - radians(?)) +
			sin(radians(?)) * sin(radians(m.latitude))
		)) AS distance`
		selectArgs = append(selectArgs, req.Latitude, req.Longitude, req.Latitude)
	}

	// 排序
	orderBy := "m.is_top DESC, m.sort_order DESC, m.rating DESC"
	switch req.SortBy {
	case "distance":
		if hasDistance {
			orderBy = "distance ASC NULLS LAST"
		}
	case "service_count":
		orderBy = "m.service_count DESC"
	case "rating":
		orderBy = "m.rating DESC, m.review_count DESC"
	}

	whereClause := strings.Join(whereParts, " AND ")

	// COUNT
	var total int64
	countSQL := "SELECT COUNT(*) FROM merchants m WHERE " + whereClause
	if err := config.GetDB().Raw(countSQL, whereArgs...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// LIST（参数顺序：先 selectArgs 再 whereArgs 再 pageSize/offset）
	listSQL := "SELECT m.*, " + selectExtras + " FROM merchants m WHERE " + whereClause +
		" ORDER BY " + orderBy + " LIMIT ? OFFSET ?"
	listArgs := append([]interface{}{}, selectArgs...)
	listArgs = append(listArgs, whereArgs...)
	listArgs = append(listArgs, pageSize, offset)

	var merchants []models.Merchant
	if err := config.GetDB().Raw(listSQL, listArgs...).Scan(&merchants).Error; err != nil {
		return nil, 0, err
	}

	return merchants, total, nil
}

// GetMerchantDetail 商家详情
func (s *MerchantService) GetMerchantDetail(id string) (*models.Merchant, error) {
	var merchant models.Merchant
	err := config.GetDB().Where("id = ?", id).First(&merchant).Error
	return &merchant, err
}