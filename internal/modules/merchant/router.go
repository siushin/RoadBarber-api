package merchant

import (
	"roadbarber/backend/internal/middleware"
	"roadbarber/backend/internal/modules/merchant/handler"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes 注册商家端路由
func RegisterRoutes(app *fiber.App) {
	scheduleHandler := handler.NewScheduleHandler()
	bookingHandler := handler.NewBookingHandler()

	api := app.Group("/api")

	// 商家入驻申请（公开）
	apply := api.Group("/merchant-apply")
	apply.Post("/apply", handler.NewMerchantApplyHandler().Apply)
	apply.Get("/apply", handler.NewMerchantApplyHandler().MyApplies)

	// 需要商家登录
	merchant := api.Group("/", middleware.Auth(), middleware.MerchantOnly())

	// 排班管理
	merchant.Post("/schedules", scheduleHandler.Create)
	merchant.Post("/schedules/batch", scheduleHandler.BatchCreate)
	merchant.Get("/schedules", scheduleHandler.List)
	merchant.Delete("/schedules/:id", scheduleHandler.Delete)

	// 预约管理
	merchant.Get("/bookings", bookingHandler.List)
	merchant.Patch("/bookings/:id/confirm", bookingHandler.Confirm)
	merchant.Patch("/bookings/:id/reject", bookingHandler.Reject)
	merchant.Patch("/bookings/:id/start", bookingHandler.Start)
	merchant.Patch("/bookings/:id/finish", bookingHandler.Finish)
}
