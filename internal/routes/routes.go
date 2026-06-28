package routes

import (
	"roadbarber/api/internal/config"
	"roadbarber/api/internal/modules/admin"
	"roadbarber/api/internal/modules/common"
	"roadbarber/api/internal/modules/customer"
	"roadbarber/api/internal/modules/merchant"
	"roadbarber/api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Setup 注册所有路由
func Setup(app *fiber.App, cfg *config.Config, sms utils.SMSProvider) {
	// 健康检查
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// 公共模块（认证、地区等所有端共用）
	common.RegisterRoutes(app, sms)

	// 顾客端模块
	customer.RegisterRoutes(app)

	// 商家端模块
	merchant.RegisterRoutes(app)

	// 管理员端模块
	admin.RegisterRoutes(app)
}