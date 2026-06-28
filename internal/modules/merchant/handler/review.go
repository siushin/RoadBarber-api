package handler

import (
	"roadbarber/api/internal/modules/merchant/service"
	"roadbarber/api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type MerchantReviewHandler struct {
	svc *service.MerchantReviewService
}

func NewMerchantReviewHandler() *MerchantReviewHandler {
	return &MerchantReviewHandler{svc: &service.MerchantReviewService{}}
}

// List 商家端评价列表（自己的店铺收到的评价）
func (h *MerchantReviewHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	var req service.MerchantListReviewsRequest
	if err := c.QueryParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	items, total, err := h.svc.ListReviews(userID, &req)
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

// Reply 商家回复评价
func (h *MerchantReviewHandler) Reply(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "评价ID不能为空")
	}

	var req struct {
		ReplyContent string `json:"reply_content"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}
	if req.ReplyContent == "" {
		return response.BadRequest(c, "回复内容不能为空")
	}

	if err := h.svc.ReplyReview(userID, id, req.ReplyContent); err != nil {
		return response.ServerError(c, "回复失败: "+err.Error())
	}
	return response.SuccessWithMessage(c, "回复成功", nil)
}