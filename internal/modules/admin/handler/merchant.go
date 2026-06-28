package handler

import (
	"roadbarber/api/internal/models"
	"roadbarber/api/internal/modules/admin/service"
	"roadbarber/api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		adminService: &service.AdminService{},
	}
}

// Dashboard 仪表盘统计
func (h *AdminHandler) Dashboard(c *fiber.Ctx) error {
	stats, err := h.adminService.GetDashboardStats()
	if err != nil {
		return response.ServerError(c, "查询统计失败")
	}
	return response.Success(c, stats)
}

// ListMerchants 商家列表
func (h *AdminHandler) ListMerchants(c *fiber.Ctx) error {
	var req service.ListMerchantsRequest
	if err := c.QueryParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	merchants, total, err := h.adminService.ListMerchants(&req)
	if err != nil {
		return response.ServerError(c, "查询商家失败")
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	return response.PageSuccess(c, merchants, total, page, pageSize)
}

// VerifyMerchant 审核商家资质
func (h *AdminHandler) VerifyMerchant(c *fiber.Ctx) error {
	auditorID, _ := c.Locals("user_id").(string)
	id := c.Params("id")

	var req struct {
		AuditStatus int8   `json:"audit_status"`
		Remark      string `json:"remark"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	if req.AuditStatus != models.AuditStatusApproved && req.AuditStatus != models.AuditStatusRejected {
		return response.BadRequest(c, "审核状态错误")
	}

	if err := h.adminService.VerifyMerchant(id, auditorID, req.AuditStatus, req.Remark); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "审核成功", nil)
}

// ListApplies 入驻申请列表
func (h *AdminHandler) ListApplies(c *fiber.Ctx) error {
	var req service.ListAppliesRequest
	if err := c.QueryParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	applies, total, err := h.adminService.ListApplies(&req)
	if err != nil {
		return response.ServerError(c, "查询申请失败")
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	return response.PageSuccess(c, applies, total, page, pageSize)
}

// ApproveApply 审核通过
func (h *AdminHandler) ApproveApply(c *fiber.Ctx) error {
	auditorID, _ := c.Locals("user_id").(string)
	id := c.Params("id")

	if err := h.adminService.ApproveApply(id, auditorID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "审核通过", nil)
}

// RejectApply 审核拒绝
func (h *AdminHandler) RejectApply(c *fiber.Ctx) error {
	auditorID, _ := c.Locals("user_id").(string)
	id := c.Params("id")

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	if err := h.adminService.RejectApply(id, auditorID, req.Reason); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "已拒绝", nil)
}