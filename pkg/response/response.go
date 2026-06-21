package response

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`           // 状态码：0 成功，非 0 失败
	Message string      `json:"message"`        // 提示信息
	Data    interface{} `json:"data,omitempty"` // 响应数据
}

// PageData 分页数据
type PageData struct {
	List     interface{} `json:"list"`      // 数据列表
	Total    int64       `json:"total"`     // 总数
	Page     int         `json:"page"`      // 当前页
	PageSize int         `json:"page_size"` // 每页数量
}

// normalizeNilSlice 把 nil slice 替换成同类型的空切片，
// 避免 GORM 没查到记录时返回的 nil slice 被 encoding/json 序列化成 null，
// 导致前端 .map(null) / notices.value.length 抛错。
func normalizeNilSlice(v interface{}) interface{} {
	if v == nil {
		return []interface{}{}
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Slice && rv.IsNil() {
		return reflect.MakeSlice(rv.Type(), 0, 0).Interface()
	}
	return v
}

// Success 成功响应
func Success(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    0,
		Message: "success",
		Data:    normalizeNilSlice(data),
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    0,
		Message: message,
		Data:    normalizeNilSlice(data),
	})
}

// PageSuccess 分页成功响应
func PageSuccess(c *fiber.Ctx, list interface{}, total int64, page, pageSize int) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    0,
		Message: "success",
		Data: PageData{
			List:     normalizeNilSlice(list),
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// Fail 失败响应
func Fail(c *fiber.Ctx, code int, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 参数错误
func BadRequest(c *fiber.Ctx, message string) error {
	return Fail(c, 400, message)
}

// Unauthorized 未授权
func Unauthorized(c *fiber.Ctx, message string) error {
	return Fail(c, 401, message)
}

// Forbidden 无权限
func Forbidden(c *fiber.Ctx, message string) error {
	return Fail(c, 403, message)
}

// NotFound 资源不存在
func NotFound(c *fiber.Ctx, message string) error {
	return Fail(c, 404, message)
}

// ServerError 服务器错误
func ServerError(c *fiber.Ctx, message string) error {
	return Fail(c, 500, message)
}
