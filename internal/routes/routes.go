package routes

import (
	"roadbarber/backend/internal/config"
	"roadbarber/backend/internal/modules/admin"
	"roadbarber/backend/internal/modules/common"
	"roadbarber/backend/internal/modules/customer"
	"roadbarber/backend/internal/modules/merchant"

	"github.com/gofiber/fiber/v2"
)

// Setup 注册所有路由
func Setup(app *fiber.App, cfg *config.Config) {
	// 公共模块（认证、地区等所有端共用）
	common.RegisterRoutes(app)

	// 顾客端模块
	customer.RegisterRoutes(app)

	// 商家端模块
	merchant.RegisterRoutes(app)

	// 管理员端模块
	admin.RegisterRoutes(app)
}
