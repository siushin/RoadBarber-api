package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"roadbarber/backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// UploadHandler 通用文件上传
type UploadHandler struct{}

// NewUploadHandler 构造函数
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

var allowedExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".pdf":  true,
}

// Upload 上传单个文件
func (h *UploadHandler) Upload(c *fiber.Ctx) error {
	// 仅商家/管理员/已登录用户可上传
	if _, ok := c.Locals("user_id").(string); !ok {
		return response.Unauthorized(c, "未登录")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return response.BadRequest(c, "请选择文件")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExt[ext] {
		return response.BadRequest(c, "不支持的文件类型: "+ext)
	}

	if file.Size > 10*1024*1024 {
		return response.BadRequest(c, "文件大小不能超过 10MB")
	}

	// 生成文件名：yyyyMM/<random>.ext
	now := time.Now()
	subdir := now.Format("200601")
	if err := os.MkdirAll(filepath.Join("./uploads", subdir), 0o755); err != nil {
		return response.ServerError(c, "无法创建目录: "+err.Error())
	}

	randBytes := make([]byte, 8)
	if _, err := rand.Read(randBytes); err != nil {
		return response.ServerError(c, "生成文件名失败: "+err.Error())
	}
	filename := fmt.Sprintf("%s%s", hex.EncodeToString(randBytes), ext)
	dst := filepath.Join("./uploads", subdir, filename)

	if err := c.SaveFile(file, dst); err != nil {
		return response.ServerError(c, "保存文件失败: "+err.Error())
	}

	url := fmt.Sprintf("/uploads/%s/%s", subdir, filename)
	return response.Success(c, fiber.Map{
		"url":      url,
		"filename": file.Filename,
		"size":     file.Size,
	})
}