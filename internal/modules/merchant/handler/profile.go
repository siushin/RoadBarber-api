package handler

import (
	"roadbarber/api/internal/modules/merchant/service"
	"roadbarber/api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	svc *service.ProfileService
}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{svc: &service.ProfileService{}}
}

// Get 获取商家自己的 profile
func (h *ProfileHandler) Get(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	data, err := h.svc.GetProfile(userID)
	if err != nil {
		return response.ServerError(c, "查询失败: "+err.Error())
	}
	return response.Success(c, data)
}

// Update 更新商家自己的 profile（title/specialties/introduction/avatar）
func (h *ProfileHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	var req service.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	if err := h.svc.UpdateProfile(userID, &req); err != nil {
		return response.ServerError(c, "更新失败: "+err.Error())
	}
	return response.SuccessWithMessage(c, "更新成功", nil)
}