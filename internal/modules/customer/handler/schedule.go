package handler

import (
	"roadbarber/api/internal/modules/customer/service"
	"roadbarber/api/pkg/response"

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

// List 商家可用排班
func (h *ScheduleHandler) List(c *fiber.Ctx) error {
	merchantID := c.Params("id")
	if merchantID == "" {
		return response.BadRequest(c, "商家ID不能为空")
	}
	date := c.Query("date")

	schedules, err := h.scheduleService.ListAvailableSchedules(merchantID, date)
	if err != nil {
		return response.ServerError(c, "查询排班失败")
	}

	return response.Success(c, schedules)
}
