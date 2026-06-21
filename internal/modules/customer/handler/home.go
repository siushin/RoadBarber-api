package handler

import (
	"roadbarber/backend/internal/modules/customer/service"
	"roadbarber/backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type HomeHandler struct {
	homeService *service.HomeService
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{
		homeService: &service.HomeService{},
	}
}

// Banners 首页 Banner 列表
func (h *HomeHandler) Banners(c *fiber.Ctx) error {
	banners, err := h.homeService.ListBanners()
	if err != nil {
		return response.ServerError(c, "查询 Banner 失败")
	}
	return response.Success(c, banners)
}

// Notices 首页公告列表（登录态按 user_id 过滤已读；支持分页）
func (h *HomeHandler) Notices(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)

	notices, total, err := h.homeService.ListNotices(userID, page, pageSize)
	if err != nil {
		return response.ServerError(c, "查询公告失败")
	}
	return response.PageSuccess(c, notices, total, page, pageSize)
}

// MarkRead 标记通知已读
func (h *HomeHandler) MarkRead(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	noticeID := c.Params("id")
	if noticeID == "" {
		return response.BadRequest(c, "通知 ID 不能为空")
	}

	if err := h.homeService.MarkNoticeRead(userID, noticeID); err != nil {
		return response.ServerError(c, "标记已读失败")
	}
	return response.SuccessWithMessage(c, "已标记已读", nil)
}