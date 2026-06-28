package merchant

import (
	"roadbarber/api/internal/middleware"
	"roadbarber/api/internal/modules/merchant/handler"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes 注册商家端路由
func RegisterRoutes(app *fiber.App) {
	scheduleHandler := handler.NewScheduleHandler()
	bookingHandler := handler.NewBookingHandler()
	customerHandler := handler.NewCustomerHandler()
	reviewHandler := handler.NewMerchantReviewHandler()
	profileHandler := handler.NewProfileHandler()
	dashboardHandler := handler.NewDashboardHandler()
	uploadHandler := handler.NewUploadHandler()

	api := app.Group("/api")

	// 商家入驻申请（公开）
	apply := api.Group("/merchant")
	apply.Post("/apply", handler.NewMerchantApplyHandler().Apply)
	apply.Get("/apply", handler.NewMerchantApplyHandler().MyApplies)

	// 文件上传（任意登录用户，仅 Auth 不限定角色）
	api.Post("/uploads", middleware.Auth(), uploadHandler.Upload)

	// 商家登录态：每条路由显式指定 Auth + MerchantOnly，
	// 避免被注册为 USE 全局中间件而污染 /api/* 其它路由
	// /bookings 改挂 /merchant/bookings，避免与顾客端 /bookings 同路径冲突

	// 排班管理
	api.Post("/schedules", middleware.Auth(), middleware.MerchantOnly(), scheduleHandler.Create)
	api.Post("/schedules/batch", middleware.Auth(), middleware.MerchantOnly(), scheduleHandler.BatchCreate)
	api.Get("/schedules", middleware.Auth(), middleware.MerchantOnly(), scheduleHandler.List)
	api.Delete("/schedules/:id", middleware.Auth(), middleware.MerchantOnly(), scheduleHandler.Delete)

	// 预约管理（挂 /merchant/bookings 避免与顾客端冲突）
	api.Get("/merchant/bookings", middleware.Auth(), middleware.MerchantOnly(), bookingHandler.List)
	api.Patch("/merchant/bookings/:id/confirm", middleware.Auth(), middleware.MerchantOnly(), bookingHandler.Confirm)
	api.Patch("/merchant/bookings/:id/reject", middleware.Auth(), middleware.MerchantOnly(), bookingHandler.Reject)
	api.Patch("/merchant/bookings/:id/start", middleware.Auth(), middleware.MerchantOnly(), bookingHandler.Start)
	api.Patch("/merchant/bookings/:id/finish", middleware.Auth(), middleware.MerchantOnly(), bookingHandler.Finish)

	// 顾客管理
	api.Get("/customers", middleware.Auth(), middleware.MerchantOnly(), customerHandler.List)

	// 评价管理
	api.Get("/reviews", middleware.Auth(), middleware.MerchantOnly(), reviewHandler.List)
	api.Post("/reviews/:id/reply", middleware.Auth(), middleware.MerchantOnly(), reviewHandler.Reply)

	// 商家资料
	api.Get("/profile", middleware.Auth(), middleware.MerchantOnly(), profileHandler.Get)
	api.Put("/profile", middleware.Auth(), middleware.MerchantOnly(), profileHandler.Update)

	// 商家仪表盘
	api.Get("/dashboard", middleware.Auth(), middleware.MerchantOnly(), dashboardHandler.Get)
}