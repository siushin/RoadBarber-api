package main

import (
	"fmt"
	"log"
	"os"

	"roadbarber/backend/internal/config"
	dbmigrate "roadbarber/backend/internal/migrate"
	"roadbarber/backend/internal/models"
	"roadbarber/backend/internal/routes"
	"roadbarber/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env 环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 初始化配置
	cfg := config.Load()

	// 运行数据库迁移（golang-migrate 自动 up 到最新版本）
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	if err := dbmigrate.Run("file://./migrations", databaseURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations applied successfully")

	// 初始化数据库连接（GORM）
	if err := config.InitDB(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// GORM AutoMigrate 作为安全网（golang-migrate 是主路径，AutoMigrate 兜底处理漏改的列）
	if err := models.AutoMigrate(); err != nil {
		log.Printf("Warning: AutoMigrate failed (non-fatal, schema migrations already applied): %v", err)
	} else {
		log.Println("AutoMigrate check completed")
	}

	// 初始化 Redis（失败时仅警告，开发期可无 Redis 启动）
	if err := config.InitRedis(cfg); err != nil {
		log.Printf("Warning: Redis not available, continuing without it: %v", err)
	}

	// 初始化短信服务商（按 env 切换 AliyunProvider / ConsoleProvider）
	smsProvider := config.NewSMSProvider()

	// 静默引用，避免 utils 包未使用告警
	_ = utils.GenerateCode

	// 创建 Fiber 应用
	app := fiber.New(fiber.Config{
		AppName:      "RoadBarber API",
		ErrorHandler: config.ErrorHandler,
	})

	// 注册全局中间件
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// 注册路由
	routes.Setup(app, cfg, smsProvider)

	// 静态文件
	if err := os.MkdirAll("./uploads", 0o755); err != nil {
		log.Printf("Warning: failed to ensure uploads dir: %v", err)
	}
	app.Static("/uploads", "./uploads")

	// 健康检查
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "RoadBarber API is running",
		})
	})

	// 启动服务
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}