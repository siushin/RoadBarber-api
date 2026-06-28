package common

import (
	"roadbarber/api/internal/middleware"
	"roadbarber/api/internal/modules/common/handler"
	"roadbarber/api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes 注册公共模块路由（认证、地区等所有端共用的）
func RegisterRoutes(app *fiber.App, sms utils.SMSProvider) {
	authHandler := handler.NewAuthHandler(sms)
	locationHandler := handler.NewLocationHandler()

	// API 分组
	api := app.Group("/api")

	// 认证接口（无需登录）
	auth := api.Group("/auth")
	auth.Post("/send-code", authHandler.SendCode)
	auth.Post("/login", authHandler.LoginByCode)
	auth.Post("/login/pwd", authHandler.LoginByPassword)
	auth.Post("/register", authHandler.Register)
	auth.Post("/logout", authHandler.Logout)

	// 需要登录的接口
	authProtected := api.Group("/auth", middleware.Auth())
	authProtected.Get("/userinfo", authHandler.GetUserInfo)

	// 地区接口（公开）
	locations := api.Group("/locations")
	locations.Get("/", locationHandler.List)
	locations.Get("/:id", locationHandler.GetChildren)
	locations.Get("/tree/all", locationHandler.Tree)
}