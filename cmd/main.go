package main

import (
	"fmt"
	"log"
	"os"

	"roadbarber/backend/internal/config"
	dbmigrate "roadbarber/backend/internal/migrate"
	"roadbarber/backend/internal/routes"

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

	// 初始化 Redis（失败时仅警告，开发期可无 Redis 启动）
	if err := config.InitRedis(cfg); err != nil {
		log.Printf("Warning: Redis not available, continuing without it: %v", err)
	}

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
	routes.Setup(app, cfg)

	// 静态文件
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
