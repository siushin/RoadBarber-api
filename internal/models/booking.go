package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Booking 预约表：用户预约理发服务
type Booking struct {
	ID              string     `gorm:"type:uuid;primaryKey" json:"id"`
	OrderNo         string     `gorm:"type:varchar(32);not null;uniqueIndex:idx_bookings_order_no" json:"order_no"`
	CustomerID      string     `gorm:"type:uuid;not null;index:idx_bookings_customer" json:"customer_id"`
	MerchantID      string     `gorm:"type:uuid;not null;index:idx_bookings_merchant" json:"merchant_id"`
	ShopID          *string    `gorm:"type:uuid;index:idx_bookings_shop" json:"shop_id"`
	ServiceID       string     `gorm:"type:uuid;not null" json:"service_id"`
	ScheduleID      string     `gorm:"type:uuid;not null" json:"schedule_id"`
	AppointmentDate time.Time  `gorm:"type:date;not null;index:idx_bookings_date" json:"appointment_date"`
	AppointmentTime string     `gorm:"type:varchar(8);not null" json:"appointment_time"`
	Duration        int        `gorm:"type:int;not null" json:"duration"`
	Price           float64    `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	Status          int8       `gorm:"type:smallint;not null;default:1;index:idx_bookings_status" json:"status"`
	CancelReason    string     `gorm:"type:varchar(255)" json:"cancel_reason"`
	CancelTime      *time.Time `json:"cancel_time"`
	Remark          string     `gorm:"type:text" json:"remark"`
	InternalNote    string     `gorm:"type:text" json:"internal_note"`
	ConfirmTime     *time.Time `json:"confirm_time"`
	StartTime       *time.Time `json:"start_time"`
	FinishTime      *time.Time `json:"finish_time"`
	CreatedAt       time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName 指定表名
func (Booking) TableName() string {
	return "bookings"
}

// BeforeCreate 创建前生成 UUID
func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// 预约状态常量
const (
	BookingStatusPending   = 1 // 待确认
	BookingStatusConfirmed = 2 // 已确认
	BookingStatusServing   = 3 // 服务中
	BookingStatusCompleted = 4 // 已完成
	BookingStatusCancelled = 5 // 已取消
	BookingStatusRejected  = 6 // 已拒绝
)
