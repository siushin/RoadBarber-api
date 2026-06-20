package handler

import (
	"roadbarber/backend/internal/modules/customer/service"
	"roadbarber/backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ShopHandler struct {
	shopService *service.ShopService
}

func NewShopHandler() *ShopHandler {
	return &ShopHandler{
		shopService: &service.ShopService{},
	}
}

// List 店铺列表
func (h *ShopHandler) List(c *fiber.Ctx) error {
	var req service.ListShopsRequest
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

// Detail 店铺详情
func (h *ShopHandler) Detail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "店铺ID不能为空")
	}

	shop, err := h.shopService.GetShopDetail(id)
	if err != nil {
		return response.NotFound(c, "店铺不存在")
	}

	return response.Success(c, shop)
}

// Merchants 店铺下的商家
func (h *ShopHandler) Merchants(c *fiber.Ctx) error {
	shopID := c.Params("id")
	if shopID == "" {
		return response.BadRequest(c, "店铺ID不能为空")
	}

	merchants, err := h.shopService.ListMerchantsByShop(shopID)
	if err != nil {
		return response.ServerError(c, "查询商家失败")
	}

	return response.Success(c, merchants)
}

// Services 店铺下的服务项目
func (h *ShopHandler) Services(c *fiber.Ctx) error {
	shopID := c.Params("id")
	if shopID == "" {
		return response.BadRequest(c, "店铺ID不能为空")
	}

	services, err := h.shopService.ListServicesByShop(shopID)
	if err != nil {
		return response.ServerError(c, "查询服务失败")
	}

	return response.Success(c, services)
}
