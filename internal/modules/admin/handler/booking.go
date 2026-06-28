package handler

import (
	"roadbarber/api/internal/modules/admin/service"
	"roadbarber/api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AdminBookingHandler struct {
	svc *service.AdminBookingService
}

func NewAdminBookingHandler() *AdminBookingHandler {
	return &AdminBookingHandler{
		svc: &service.AdminBookingService{},
	}
}

// List 预约列表（管理员端）
func (h *AdminBookingHandler) List(c *fiber.Ctx) error {
	var req service.AdminListBookingsRequest
	if err := c.QueryParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	items, total, err := h.svc.ListBookings(&req)
	if err != nil {
		return response.ServerError(c, "查询预约失败")
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	return response.PageSuccess(c, items, total, page, pageSize)
}