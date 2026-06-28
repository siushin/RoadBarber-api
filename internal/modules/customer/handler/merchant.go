package handler

import (
	"roadbarber/api/internal/modules/customer/service"
	"roadbarber/api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type MerchantHandler struct {
	merchantService *service.MerchantService
}

func NewMerchantHandler() *MerchantHandler {
	return &MerchantHandler{
		merchantService: &service.MerchantService{},
	}
}

// List 商家列表
func (h *MerchantHandler) List(c *fiber.Ctx) error {
	var req service.ListMerchantsRequest
	if err := c.QueryParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	merchants, total, err := h.merchantService.ListMerchants(&req)
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

// Detail 商家详情
func (h *MerchantHandler) Detail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "商家ID不能为空")
	}

	merchant, err := h.merchantService.GetMerchantDetail(id)
	if err != nil {
		return response.NotFound(c, "商家不存在")
	}

	return response.Success(c, merchant)
}
