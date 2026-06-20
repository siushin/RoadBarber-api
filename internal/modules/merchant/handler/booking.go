package handler

import (
	"roadbarber/backend/internal/modules/merchant/service"
	"roadbarber/backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler() *BookingHandler {
	return &BookingHandler{
		bookingService: &service.BookingService{},
	}
}

// List 商家预约列表
func (h *BookingHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	var req service.ListMerchantBookingsRequest
	if err := c.QueryParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}
	req.MerchantID = userID

	bookings, total, err := h.bookingService.ListMerchantBookings(&req)
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

	return response.PageSuccess(c, bookings, total, page, pageSize)
}

// Confirm 确认预约
func (h *BookingHandler) Confirm(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	id := c.Params("id")
	if err := h.bookingService.ConfirmBooking(id, userID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "确认成功", nil)
}

// Reject 拒绝预约
func (h *BookingHandler) Reject(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	id := c.Params("id")
	reason := c.Query("reason", "商家拒绝")

	if err := h.bookingService.RejectBooking(id, userID, reason); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "已拒绝", nil)
}

// Start 开始服务
func (h *BookingHandler) Start(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	id := c.Params("id")
	if err := h.bookingService.StartService(id, userID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "已开始服务", nil)
}

// Finish 完成服务
func (h *BookingHandler) Finish(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	id := c.Params("id")
	if err := h.bookingService.FinishService(id, userID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "服务已完成", nil)
}
