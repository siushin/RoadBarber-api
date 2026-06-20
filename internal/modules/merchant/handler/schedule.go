package handler

import (
	"roadbarber/backend/internal/modules/merchant/service"
	"roadbarber/backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ScheduleHandler struct {
	scheduleService *service.ScheduleService
}

func NewScheduleHandler() *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: &service.ScheduleService{},
	}
}

// Create 发布排班
func (h *ScheduleHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	var req service.CreateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	schedule, err := h.scheduleService.CreateSchedule(userID, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, schedule)
}

// BatchCreate 批量发布
func (h *ScheduleHandler) BatchCreate(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	var req service.BatchCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	if err := h.scheduleService.BatchCreate(userID, &req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "发布成功", nil)
}

// List 我的排班
func (h *ScheduleHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	schedules, err := h.scheduleService.ListMySchedules(userID, startDate, endDate)
	if err != nil {
		return response.ServerError(c, "查询排班失败")
	}

	return response.Success(c, schedules)
}

// Delete 删除排班
func (h *ScheduleHandler) Delete(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	id := c.Params("id")
	if err := h.scheduleService.DeleteSchedule(id, userID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "删除成功", nil)
}
