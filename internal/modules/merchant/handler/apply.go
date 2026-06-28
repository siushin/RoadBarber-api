package handler

import (
	"roadbarber/api/internal/modules/merchant/service"
	"roadbarber/api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type MerchantApplyHandler struct {
	applyService *service.MerchantApplyService
}

func NewMerchantApplyHandler() *MerchantApplyHandler {
	return &MerchantApplyHandler{
		applyService: &service.MerchantApplyService{},
	}
}

// Apply 提交入驻申请
func (h *MerchantApplyHandler) Apply(c *fiber.Ctx) error {
	var req service.CreateApplyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	apply, err := h.applyService.SubmitApply(&req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, apply)
}

// MyApplies 我的申请记录
func (h *MerchantApplyHandler) MyApplies(c *fiber.Ctx) error {
	phone := c.Query("phone")
	if phone == "" {
		return response.BadRequest(c, "手机号不能为空")
	}

	applies, err := h.applyService.ListMyApplies(phone)
	if err != nil {
		return response.ServerError(c, "查询失败")
	}

	return response.Success(c, applies)
}
