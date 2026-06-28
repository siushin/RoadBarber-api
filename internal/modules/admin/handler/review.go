package handler

import (
	"roadbarber/api/internal/modules/admin/service"
	"roadbarber/api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AdminReviewHandler struct {
	svc *service.AdminReviewService
}

func NewAdminReviewHandler() *AdminReviewHandler {
	return &AdminReviewHandler{
		svc: &service.AdminReviewService{},
	}
}

// List 评价列表（管理员端）
func (h *AdminReviewHandler) List(c *fiber.Ctx) error {
	var req service.AdminListReviewsRequest
	if err := c.QueryParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	items, total, err := h.svc.ListReviews(&req)
	if err != nil {
		return response.ServerError(c, "查询评价失败")
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

// Delete 删除评价（软删：status=2）
func (h *AdminReviewHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "评价ID不能为空")
	}

	if err := h.svc.DeleteReview(id); err != nil {
		return response.ServerError(c, "删除评价失败: "+err.Error())
	}
	return response.SuccessWithMessage(c, "删除成功", nil)
}