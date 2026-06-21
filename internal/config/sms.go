package config

import (
	"log"
	"os"

	"roadbarber/backend/pkg/utils"
)

// NewSMSProvider 根据环境变量选择短信服务商：
//   - 配置齐全（ALIYUN_ACCESS_KEY_ID/SECRET/SIGN_NAME/TEMPLATE_CODE）→ AliyunProvider
//   - 否则降级到 ConsoleProvider（开发环境 stdout 输出）
func NewSMSProvider() utils.SMSProvider {
	accessKeyID := os.Getenv("ALIYUN_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	signName := os.Getenv("ALIYUN_SIGN_NAME")
	templateCode := os.Getenv("ALIYUN_TEMPLATE_CODE")

	if accessKeyID == "" || accessKeySecret == "" || signName == "" || templateCode == "" {
		log.Println("[SMS] 阿里云短信配置不完整，使用 ConsoleProvider（仅打印日志）")
		return &utils.ConsoleProvider{}
	}

	provider, err := utils.NewAliyunProvider(accessKeyID, accessKeySecret, signName, templateCode)
	if err != nil {
		log.Printf("[SMS] 创建 AliyunProvider 失败，降级到 ConsoleProvider: %v", err)
		return &utils.ConsoleProvider{}
	}
	log.Println("[SMS] 使用 AliyunProvider 发送短信")
	return provider
}