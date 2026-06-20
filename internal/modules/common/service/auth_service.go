package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"roadbarber/backend/internal/config"
	"roadbarber/backend/internal/models"
	"roadbarber/backend/pkg/utils"

	"gorm.io/gorm"
)

type AuthService struct{}

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Phone string `json:"phone" validate:"required"`
}

// LoginByCodeRequest 验证码登录请求
type LoginByCodeRequest struct {
	Phone string `json:"phone" validate:"required"`
	Code  string `json:"code" validate:"required"`
}

// LoginByPasswordRequest 密码登录请求
type LoginByPasswordRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	Nickname string `json:"nickname" validate:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string      `json:"token"`
	UserInfo models.User `json:"user_info"`
}

// SendCode 发送短信验证码
func (s *AuthService) SendCode(phone string) error {
	if phone == "" {
		return errors.New("手机号不能为空")
	}

	// 生成 6 位验证码
	code := utils.GenerateCode()

	// 存储到 Redis，5 分钟过期
	ctx := context.Background()
	key := fmt.Sprintf("sms:code:%s", phone)
	err := config.RedisClient.Set(ctx, key, code, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("保存验证码失败: %w", err)
	}

	// TODO: 接入实际短信服务商发送验证码
	// 调试阶段直接打印到日志
	fmt.Printf("[SMS] Phone: %s, Code: %s\n", phone, code)

	return nil
}

// LoginByCode 验证码登录
func (s *AuthService) LoginByCode(phone, code string) (*LoginResponse, error) {
	if phone == "" || code == "" {
		return nil, errors.New("手机号和验证码不能为空")
	}

	// 校验验证码
	ctx := context.Background()
	key := fmt.Sprintf("sms:code:%s", phone)
	storedCode, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, errors.New("验证码已过期或不存在")
	}
	if storedCode != code {
		return nil, errors.New("验证码错误")
	}

	// 删除验证码
	config.RedisClient.Del(ctx, key)

	// 查找或创建用户
	var user models.User
	err = config.GetDB().Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 自动注册
			user = models.User{
				Phone:    phone,
				Nickname: "用户" + phone[7:],
				Role:     models.RoleCustomer,
				Status:   models.UserStatusNormal,
			}
			if err := config.GetDB().Create(&user).Error; err != nil {
				return nil, fmt.Errorf("创建用户失败: %w", err)
			}
		} else {
			return nil, fmt.Errorf("查询用户失败: %w", err)
		}
	}

	// 检查用户状态
	if user.Status == models.UserStatusDisabled {
		return nil, errors.New("账号已被禁用")
	}

	// 生成 Token
	token, err := utils.GenerateToken(user.ID, int(user.Role), config.GetJWTSecret(), config.GetJWTExpiresIn())
	if err != nil {
		return nil, fmt.Errorf("生成 Token 失败: %w", err)
	}

	// 更新最后登录时间
	now := time.Now()
	config.GetDB().Model(&user).Updates(map[string]interface{}{
		"last_login_at": now,
	})

	return &LoginResponse{
		Token:    token,
		UserInfo: user,
	}, nil
}

// LoginByPassword 密码登录
func (s *AuthService) LoginByPassword(phone, password string) (*LoginResponse, error) {
	if phone == "" || password == "" {
		return nil, errors.New("手机号和密码不能为空")
	}

	var user models.User
	err := config.GetDB().Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账号不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 校验密码
	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, errors.New("密码错误")
	}

	// 检查用户状态
	if user.Status == models.UserStatusDisabled {
		return nil, errors.New("账号已被禁用")
	}

	// 生成 Token
	token, err := utils.GenerateToken(user.ID, int(user.Role), config.GetJWTSecret(), config.GetJWTExpiresIn())
	if err != nil {
		return nil, fmt.Errorf("生成 Token 失败: %w", err)
	}

	// 更新最后登录时间
	now := time.Now()
	config.GetDB().Model(&user).Updates(map[string]interface{}{
		"last_login_at": now,
	})

	return &LoginResponse{
		Token:    token,
		UserInfo: user,
	}, nil
}

// Register 顾客注册
func (s *AuthService) Register(phone, password, nickname string) (*LoginResponse, error) {
	if phone == "" || password == "" || nickname == "" {
		return nil, errors.New("参数不完整")
	}
	if len(password) < 6 {
		return nil, errors.New("密码长度不能少于 6 位")
	}

	// 检查手机号是否已注册
	var existUser models.User
	err := config.GetDB().Where("phone = ?", phone).First(&existUser).Error
	if err == nil {
		return nil, errors.New("手机号已注册")
	}

	// 密码哈希
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建用户
	user := models.User{
		Phone:        phone,
		PasswordHash: hashedPassword,
		Nickname:     nickname,
		Role:         models.RoleCustomer,
		Status:       models.UserStatusNormal,
	}
	if err := config.GetDB().Create(&user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 生成 Token
	token, err := utils.GenerateToken(user.ID, int(user.Role), config.GetJWTSecret(), config.GetJWTExpiresIn())
	if err != nil {
		return nil, fmt.Errorf("生成 Token 失败: %w", err)
	}

	return &LoginResponse{
		Token:    token,
		UserInfo: user,
	}, nil
}

// GetUserInfo 获取当前用户信息
func (s *AuthService) GetUserInfo(userID string) (*models.User, error) {
	var user models.User
	err := config.GetDB().Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}
