package utils

import (
	"fmt"
	"log"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

// SMSProvider 短信服务商接口
type SMSProvider interface {
	Send(phone, code string) error
}

// ConsoleProvider 调试用：直接把验证码打印到日志
type ConsoleProvider struct{}

// Send 控制台输出验证码（仅用于开发环境）
func (p *ConsoleProvider) Send(phone, code string) error {
	log.Printf("[SMS-CONSOLE] phone=%s code=%s", phone, code)
	return nil
}

// AliyunProvider 阿里云短信实现
type AliyunProvider struct {
	client       *dysmsapi.Client
	SignName     string
	TemplateCode string
}

// NewAliyunProvider 创建阿里云短信客户端
func NewAliyunProvider(accessKeyID, accessKeySecret, signName, templateCode string) (*AliyunProvider, error) {
	if accessKeyID == "" || accessKeySecret == "" || signName == "" || templateCode == "" {
		return nil, fmt.Errorf("阿里云短信配置不完整")
	}

	cfg := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyID),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("dysmsapi.aliyuncs.com"),
	}

	client, err := dysmsapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("创建阿里云短信客户端失败: %w", err)
	}

	return &AliyunProvider{
		client:       client,
		SignName:     signName,
		TemplateCode: templateCode,
	}, nil
}

// Send 通过阿里云发送短信验证码
func (p *AliyunProvider) Send(phone, code string) error {
	request := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(p.SignName),
		TemplateCode:  tea.String(p.TemplateCode),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%s\"}", code)),
	}

	response, err := p.client.SendSms(request)
	if err != nil {
		return fmt.Errorf("调用阿里云短信失败: %w", err)
	}

	if response.Body == nil || response.Body.Code == nil || *response.Body.Code != "OK" {
		msg := ""
		if response.Body != nil && response.Body.Message != nil {
			msg = *response.Body.Message
		}
		return fmt.Errorf("阿里云短信返回错误: %s", msg)
	}

	log.Printf("[SMS-ALIYUN] phone=%s code=%s requestId=%v at=%s",
		phone, code,
		tea.StringValue(response.Body.RequestId),
		time.Now().Format(time.RFC3339),
	)
	return nil
}