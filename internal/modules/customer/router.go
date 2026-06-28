package customer

import (
	"roadbarber/api/internal/middleware"
	"roadbarber/api/internal/modules/customer/handler"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes 注册顾客端路由
func RegisterRoutes(app *fiber.App) {
	shopHandler := handler.NewShopHandler()
	merchantHandler := handler.NewMerchantHandler()
	scheduleHandler := handler.NewScheduleHandler()
	bookingHandler := handler.NewBookingHandler()
	reviewHandler := handler.NewReviewHandler()
	favoriteHandler := handler.NewFavoriteHandler()
	homeHandler := handler.NewHomeHandler()

	// API 分组（仅作路径前缀，不会注册 USE 中间件）
	api := app.Group("/api")

	// 首页运营内容
	// - /home/banners：公开
	// - /home/notices：可选登录（中间件解析 token，未登录 locals 为空；登录则按 user_id 过滤已读）
	api.Get("/home/banners", homeHandler.Banners)
	api.Get("/home/notices", middleware.OptionalAuth(), homeHandler.Notices)

	// 通知已读：必须登录
	api.Post("/notices/:id/read", middleware.Auth(), homeHandler.MarkRead)

	// 店铺（公开）
	shops := api.Group("/shops")
	shops.Get("/", shopHandler.List)
	shops.Get("/:id", shopHandler.Detail)
	shops.Get("/:id/barbers", shopHandler.Merchants)
	shops.Get("/:id/services", shopHandler.Services)

	// 商家（公开浏览）
	merchants := api.Group("/merchants")
	merchants.Get("/", merchantHandler.List)
	merchants.Get("/:id", merchantHandler.Detail)
	merchants.Get("/:id/schedules", scheduleHandler.List)
	merchants.Get("/:id/reviews", reviewHandler.ListByMerchant)

	// 顾客登录态：每条路由显式指定 Auth + CustomerOnly，
	// 避免被注册为 USE 全局中间件而污染 /api/* 其它路由
	// 预约
	api.Post("/bookings", middleware.Auth(), middleware.CustomerOnly(), bookingHandler.Create)
	api.Get("/bookings", middleware.Auth(), middleware.CustomerOnly(), bookingHandler.List)
	api.Get("/bookings/:id", middleware.Auth(), middleware.CustomerOnly(), bookingHandler.Detail)
	api.Patch("/bookings/:id/cancel", middleware.Auth(), middleware.CustomerOnly(), bookingHandler.Cancel)

	// 评价
	api.Post("/reviews", middleware.Auth(), middleware.CustomerOnly(), reviewHandler.Create)

	// 收藏
	api.Post("/merchants/:id/favorite", middleware.Auth(), middleware.CustomerOnly(), favoriteHandler.Add)
	api.Delete("/merchants/:id/favorite", middleware.Auth(), middleware.CustomerOnly(), favoriteHandler.Remove)
	api.Get("/favorites", middleware.Auth(), middleware.CustomerOnly(), favoriteHandler.List)
}