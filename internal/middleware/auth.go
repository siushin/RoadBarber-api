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

// OptionalAuth 可选 JWT 中间件：有 token 时解析 user_id 进 locals，无 token 或非法 token
// 时直接放行（不报错）。适用于公开接口但登录态下需要按用户上下文做差异化的场景。
func OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		token := parts[1]
		claims, err := utils.ParseToken(token, config.GetJWTSecret())
		if err != nil {
			return c.Next()
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

// RequireRole 角色权限中间件
// 超级管理员（role=3）拥有所有角色权限，无视 roles 列表直接放行
func RequireRole(roles ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(int)
		if !ok {
			return response.Unauthorized(c, "未登录")
		}

		// 超级管理员直通
		if role == 3 {
			return c.Next()
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
