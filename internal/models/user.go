package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 用户表：存储所有用户（顾客、商家、管理员）
type User struct {
	ID           string     `gorm:"type:uuid;primaryKey" json:"id"`
	Phone        string     `gorm:"type:varchar(20);uniqueIndex:idx_users_phone" json:"phone"`
	Username     *string    `gorm:"type:varchar(50);column:username;index:idx_users_username" json:"username"`
	Email        *string    `gorm:"type:varchar(255);column:email;index:idx_users_email" json:"email"`
	PasswordHash string     `gorm:"type:varchar(255)" json:"-"`
	Nickname     string     `gorm:"type:varchar(50);not null" json:"nickname"`
	Avatar       string     `gorm:"type:varchar(500)" json:"avatar"`
	Gender       int8       `gorm:"type:smallint;default:0" json:"gender"`
	Role         int8       `gorm:"type:smallint;not null;default:1;index:idx_users_role" json:"role"`
	Status       int8       `gorm:"type:smallint;not null;default:1;index:idx_users_status" json:"status"`
	OpenID       *string    `gorm:"type:varchar(128);column:openid;index:idx_users_openid" json:"openid"`
	UnionID      *string    `gorm:"type:varchar(128);column:unionid;index:idx_users_unionid" json:"unionid"`
	WxNickname   *string    `gorm:"type:varchar(100);column:wx_nickname" json:"wx_nickname"`
	WxAvatar     *string    `gorm:"type:varchar(500);column:wx_avatar" json:"wx_avatar"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	LastLoginIP  string     `gorm:"type:varchar(50)" json:"last_login_ip"`
	CreatedAt    time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前生成 UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// 用户角色常量
const (
	RoleCustomer = 1 // 顾客
	RoleMerchant = 2 // 商家
	RoleAdmin    = 3 // 管理员
)

// 用户状态常量
const (
	UserStatusNormal    = 1 // 正常
	UserStatusPending   = 2 // 待审核
	UserStatusDisabled  = 3 // 禁用
)