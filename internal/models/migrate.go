package models

import (
	"roadbarber/backend/internal/config"
)

// AutoMigrate 自动迁移所有数据表
func AutoMigrate() error {
	db := config.GetDB()
	return db.AutoMigrate(
		&User{},
		&MerchantProfile{},
		&Location{},
		&Shop{},
		&Merchant{},
		&Service{},
		&MerchantService{},
		&Schedule{},
		&Booking{},
		&Review{},
		&MerchantApply{},
		&Favorite{},
	)
}
