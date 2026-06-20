package handler

import (
	"roadbarber/backend/internal/modules/customer/service"
	"roadbarber/backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type FavoriteHandler struct {
	favoriteService *service.FavoriteService
}

func NewFavoriteHandler() *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: &service.FavoriteService{},
	}
}

// Add 收藏
func (h *FavoriteHandler) Add(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	merchantID := c.Params("id")
	if merchantID == "" {
		return response.BadRequest(c, "商家ID不能为空")
	}

	if err := h.favoriteService.AddFavorite(userID, merchantID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "收藏成功", nil)
}

// Remove 取消收藏
func (h *FavoriteHandler) Remove(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	merchantID := c.Params("id")
	if merchantID == "" {
		return response.BadRequest(c, "商家ID不能为空")
	}

	if err := h.favoriteService.RemoveFavorite(userID, merchantID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "取消收藏成功", nil)
}

// List 我的收藏
func (h *FavoriteHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)

	merchants, total, err := h.favoriteService.ListMyFavorites(userID, page, pageSize)
	if err != nil {
		return response.ServerError(c, "查询收藏失败")
	}

	return response.PageSuccess(c, merchants, total, page, pageSize)
}
