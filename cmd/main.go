package main

import (
	"log"
	"os"

	"roadbarber/backend/internal/config"
	"roadbarber/backend/internal/models"
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

	// 初始化数据库
	if err := config.InitDB(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表
	if err := models.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化 Redis
	if err := config.InitRedis(cfg); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
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
