package handler

import (
	"roadbarber/api/internal/modules/admin/service"
	"roadbarber/api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AdminServiceHandler struct {
	svc *service.AdminServiceItemService
}

func NewAdminServiceHandler() *AdminServiceHandler {
	return &AdminServiceHandler{
		svc: &service.AdminServiceItemService{},
	}
}

// List 服务列表（管理员端）
func (h *AdminServiceHandler) List(c *fiber.Ctx) error {
	var req service.AdminListServicesRequest
	if err := c.QueryParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	items, total, err := h.svc.ListServices(&req)
	if err != nil {
		return response.ServerError(c, "查询服务失败")
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

// Create 创建服务
func (h *AdminServiceHandler) Create(c *fiber.Ctx) error {
	var req service.AdminServiceUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}
	if req.ShopID == "" || req.Name == "" {
		return response.BadRequest(c, "店铺ID和服务名称必填")
	}

	item, err := h.svc.CreateService(&req)
	if err != nil {
		return response.ServerError(c, "创建服务失败: "+err.Error())
	}
	return response.SuccessWithMessage(c, "创建成功", item)
}

// Update 更新服务
func (h *AdminServiceHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "服务ID不能为空")
	}

	var req service.AdminServiceUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	item, err := h.svc.UpdateService(id, &req)
	if err != nil {
		return response.ServerError(c, "更新服务失败: "+err.Error())
	}
	return response.SuccessWithMessage(c, "更新成功", item)
}