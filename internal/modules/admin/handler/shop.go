package handler

import (
	"roadbarber/backend/internal/modules/admin/service"
	"roadbarber/backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AdminShopHandler struct {
	shopService *service.AdminShopService
}

func NewAdminShopHandler() *AdminShopHandler {
	return &AdminShopHandler{
		shopService: &service.AdminShopService{},
	}
}

// List 店铺列表（管理员端）
func (h *AdminShopHandler) List(c *fiber.Ctx) error {
	var req service.AdminListShopsRequest
	if err := c.QueryParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	shops, total, err := h.shopService.ListShops(&req)
	if err != nil {
		return response.ServerError(c, "查询店铺失败")
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	return response.PageSuccess(c, shops, total, page, pageSize)
}

// Create 创建店铺
func (h *AdminShopHandler) Create(c *fiber.Ctx) error {
	var req service.AdminShopUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}
	if req.Name == "" {
		return response.BadRequest(c, "店铺名称不能为空")
	}

	creatorID, _ := c.Locals("user_id").(string)
	shop, err := h.shopService.CreateShop(&req, creatorID)
	if err != nil {
		return response.ServerError(c, "创建店铺失败: "+err.Error())
	}
	return response.SuccessWithMessage(c, "创建成功", shop)
}

// Update 更新店铺
func (h *AdminShopHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "店铺ID不能为空")
	}

	var req service.AdminShopUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	shop, err := h.shopService.UpdateShop(id, &req)
	if err != nil {
		return response.ServerError(c, "更新店铺失败: "+err.Error())
	}
	return response.SuccessWithMessage(c, "更新成功", shop)
}