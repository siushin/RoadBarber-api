package handler

import (
	"roadbarber/backend/internal/modules/merchant/service"
	"roadbarber/backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	svc *service.CustomerService
}

func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{svc: &service.CustomerService{}}
}

// List 商家端顾客列表（去重自 bookings）
func (h *CustomerHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	customers, err := h.svc.ListCustomersByMerchantUser(userID)
	if err != nil {
		return response.ServerError(c, "查询顾客失败: "+err.Error())
	}
	return response.Success(c, customers)
}