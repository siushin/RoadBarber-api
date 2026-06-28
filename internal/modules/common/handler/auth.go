package handler

import (
	"roadbarber/api/internal/modules/common/service"
	"roadbarber/api/pkg/response"
	"roadbarber/api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(sms utils.SMSProvider) *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(sms),
	}
}

// SendCode 发送短信验证码
func (h *AuthHandler) SendCode(c *fiber.Ctx) error {
	var req service.SendCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	if err := h.authService.SendCode(req.Phone); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "验证码已发送", nil)
}

// LoginByCode 验证码登录
func (h *AuthHandler) LoginByCode(c *fiber.Ctx) error {
	var req service.LoginByCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	data, err := h.authService.LoginByCode(req.Phone, req.Code)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, data)
}

// LoginByPassword 密码登录（支持用户名/手机号/邮箱任一账号）
func (h *AuthHandler) LoginByPassword(c *fiber.Ctx) error {
	var req service.LoginByPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	// 兼容：小程序端仍传 phone 字段；后台传 account
	account := req.Account
	if account == "" {
		account = req.Phone
	}

	data, err := h.authService.LoginByPassword(account, req.Password, int8(req.Role))
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, data)
}

// Register 顾客注册
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req service.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "参数解析失败")
	}

	data, err := h.authService.Register(req.Phone, req.Password, req.Nickname)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, data)
}

// Logout 登出（前端清除 token 即可，后端仅做记录）
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return response.SuccessWithMessage(c, "登出成功", nil)
}

// GetUserInfo 获取当前用户信息
func (h *AuthHandler) GetUserInfo(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "未登录")
	}

	user, err := h.authService.GetUserInfo(userID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, user)
}
