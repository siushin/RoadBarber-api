package customer

import (
	"roadbarber/backend/internal/middleware"
	"roadbarber/backend/internal/modules/customer/handler"

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

	// 公开 API（顾客端部分接口无需登录也可访问）
	api := app.Group("/api")

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

	// 需要登录的 API
	customer := api.Group("/", middleware.Auth(), middleware.CustomerOnly())

	// 预约
	customer.Post("/bookings", bookingHandler.Create)
	customer.Get("/bookings", bookingHandler.List)
	customer.Get("/bookings/:id", bookingHandler.Detail)
	customer.Patch("/bookings/:id/cancel", bookingHandler.Cancel)

	// 评价
	customer.Post("/reviews", reviewHandler.Create)

	// 收藏
	customer.Post("/merchants/:id/favorite", favoriteHandler.Add)
	customer.Delete("/merchants/:id/favorite", favoriteHandler.Remove)
	customer.Get("/favorites", favoriteHandler.List)
}
