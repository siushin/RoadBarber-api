package handler

import (
	"roadbarber/api/internal/modules/customer/service"
	"roadbarber/api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{
		reviewService: &service.ReviewService{},
	}
}

// Create 提交评价
func (h *ReviewHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	var req service.CreateReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	review, err := h.reviewService.CreateReview(userID, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, review)
}

// ListByMerchant 商家评价列表
func (h *ReviewHandler) ListByMerchant(c *fiber.Ctx) error {
	merchantID := c.Params("id")
	if merchantID == "" {
		return response.BadRequest(c, "商家ID不能为空")
	}

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)

	reviews, total, err := h.reviewService.ListByMerchant(merchantID, page, pageSize)
	if err != nil {
		return response.ServerError(c, "查询评价失败")
	}

	return response.PageSuccess(c, reviews, total, page, pageSize)
}
