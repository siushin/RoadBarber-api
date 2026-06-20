package handler

import (
	"roadbarber/backend/internal/modules/customer/service"
	"roadbarber/backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ServiceCatalogHandler struct {
	service *service.ServiceCatalogService
}

func NewServiceCatalogHandler() *ServiceCatalogHandler {
	return &ServiceCatalogHandler{service: &service.ServiceCatalogService{}}
}

// ListByShop 店铺下的服务
func (h *ServiceCatalogHandler) ListByShop(c *fiber.Ctx) error {
	shopID := c.Params("id")
	if shopID == "" {
		return response.BadRequest(c, "店铺ID不能为空")
	}

	services, err := h.service.ListByShop(shopID)
	if err != nil {
		return response.ServerError(c, "查询服务失败")
	}

	return response.Success(c, services)
}
