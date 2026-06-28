package handler

import (
	"roadbarber/api/internal/modules/merchant/service"
	"roadbarber/api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	svc *service.MerchantDashboardService
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{svc: &service.MerchantDashboardService{}}
}

// Get 商家仪表盘统计
func (h *DashboardHandler) Get(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	data, err := h.svc.GetStats(userID)
	if err != nil {
		return response.ServerError(c, "查询失败: "+err.Error())
	}
	return response.Success(c, data)
}