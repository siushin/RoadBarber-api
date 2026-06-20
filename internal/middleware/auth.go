package middleware

import (
	"strings"

	"roadbarber/backend/internal/config"
	"roadbarber/backend/pkg/response"
	"roadbarber/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Auth JWT 认证中间件
func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "未提供认证令牌")
		}

		// Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Unauthorized(c, "认证格式错误")
		}

		token := parts[1]
		claims, err := utils.ParseToken(token, config.GetJWTSecret())
		if err != nil {
			return response.Unauthorized(c, "认证令牌无效或已过期")
		}

		// 将用户信息存入上下文
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// RequireRole 角色权限中间件
func RequireRole(roles ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(int)
		if !ok {
			return response.Unauthorized(c, "未登录")
		}

		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}

		return response.Forbidden(c, "无权限访问")
	}
}

// CustomerOnly 顾客端
func CustomerOnly() fiber.Handler {
	return RequireRole(1)
}

// MerchantOnly 商家端
func MerchantOnly() fiber.Handler {
	return RequireRole(2)
}

// AdminOnly 管理员端
func AdminOnly() fiber.Handler {
	return RequireRole(3)
}
