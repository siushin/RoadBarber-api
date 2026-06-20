package admin

import (
	"roadbarber/backend/internal/middleware"
	"roadbarber/backend/internal/modules/admin/handler"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes 注册管理员端路由
func RegisterRoutes(app *fiber.App) {
	adminHandler := handler.NewAdminHandler()

	api := app.Group("/api")

	// 需要管理员登录
	admin := api.Group("/admin", middleware.Auth(), middleware.AdminOnly())

	// 仪表盘
	admin.Get("/dashboard", adminHandler.Dashboard)

	// 商家管理
	admin.Get("/merchants", adminHandler.ListMerchants)
	admin.Patch("/merchants/:id/verify", adminHandler.VerifyMerchant)

	// 商家入驻申请
	admin.Get("/merchant-apply/applies", adminHandler.ListApplies)
	admin.Patch("/merchant-apply/applies/:id/approve", adminHandler.ApproveApply)
	admin.Patch("/merchant-apply/applies/:id/reject", adminHandler.RejectApply)
}
