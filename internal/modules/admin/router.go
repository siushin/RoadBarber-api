package admin

import (
	"roadbarber/api/internal/middleware"
	"roadbarber/api/internal/modules/admin/handler"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes 注册管理员端路由
func RegisterRoutes(app *fiber.App) {
	adminHandler := handler.NewAdminHandler()
	shopHandler := handler.NewAdminShopHandler()
	serviceHandler := handler.NewAdminServiceHandler()
	bookingHandler := handler.NewAdminBookingHandler()
	reviewHandler := handler.NewAdminReviewHandler()

	api := app.Group("/api")

	// 需要管理员登录
	admin := api.Group("/admin", middleware.Auth(), middleware.AdminOnly())

	// 仪表盘
	admin.Get("/dashboard", adminHandler.Dashboard)

	// 商家管理
	admin.Get("/merchants", adminHandler.ListMerchants)
	admin.Patch("/merchants/:id/verify", adminHandler.VerifyMerchant)

	// 商家入驻申请
	admin.Get("/merchant/applies", adminHandler.ListApplies)
	admin.Patch("/merchant/applies/:id/approve", adminHandler.ApproveApply)
	admin.Patch("/merchant/applies/:id/reject", adminHandler.RejectApply)

	// 店铺管理
	admin.Get("/shops", shopHandler.List)
	admin.Post("/shops", shopHandler.Create)
	admin.Put("/shops/:id", shopHandler.Update)

	// 服务项目
	admin.Get("/services", serviceHandler.List)
	admin.Post("/services", serviceHandler.Create)
	admin.Put("/services/:id", serviceHandler.Update)

	// 预约管理
	admin.Get("/bookings", bookingHandler.List)

	// 评价管理
	admin.Get("/reviews", reviewHandler.List)
	admin.Delete("/reviews/:id", reviewHandler.Delete)
}