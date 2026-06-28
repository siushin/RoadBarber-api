package handler

import (
	"roadbarber/api/internal/modules/customer/service"
	"roadbarber/api/pkg/response"

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

// Create 创建预约
func (h *BookingHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	var req service.CreateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	booking, err := h.bookingService.CreateBooking(userID, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, booking)
}

// List 我的预约
func (h *BookingHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	req := service.ListMyBookingsRequest{CustomerID: userID}
	if err := c.QueryParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	bookings, total, err := h.bookingService.ListMyBookings(&req)
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

// Detail 预约详情
func (h *BookingHandler) Detail(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	id := c.Params("id")
	booking, err := h.bookingService.GetBookingDetail(id, userID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, booking)
}

// Cancel 取消预约
func (h *BookingHandler) Cancel(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	id := c.Params("id")
	reason := c.Query("reason", "用户取消")

	if err := h.bookingService.CancelBooking(id, userID, reason); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "取消成功", nil)
}
