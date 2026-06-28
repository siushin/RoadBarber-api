package service

import (
	"roadbarber/api/internal/config"
	"roadbarber/api/internal/models"

	"gorm.io/gorm"
)

type HomeService struct{}

// ListBanners 首页 Banner 列表（按 sort_order DESC, created_at DESC）
func (s *HomeService) ListBanners() ([]models.Banner, error) {
	var banners []models.Banner
	err := config.GetDB().
		Where("status = ?", models.BannerStatusActive).
		Order("sort_order DESC, created_at DESC").
		Find(&banners).Error
	return banners, err
}

// ListNotices 首页公告列表（按 sort_order DESC, created_at DESC）
// userID 非空时按 notice_reads 过滤掉当前用户已读的公告。
// 不再合并 icon + content：icon 作为 prefix 独立渲染在 wd-notice-bar 左侧，
// content 是纯文本，text_color 控制文字颜色。
// 支持 page / page_size 分页，前端关闭一条后 page+1 拿下一条。
func (s *HomeService) ListNotices(userID string, page, pageSize int) ([]models.Notice, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	db := config.GetDB().Model(&models.Notice{}).
		Where("status = ?", models.NoticeStatusActive)

	if userID != "" {
		// 找到当前用户所有已读的 notice_id
		var readIDs []string
		if err := config.GetDB().
			Model(&models.NoticeRead{}).
			Where("user_id = ?", userID).
			Pluck("notice_id", &readIDs).Error; err != nil {
			return nil, 0, err
		}
		if len(readIDs) > 0 {
			db = db.Where("id NOT IN ?", readIDs)
		}
	}

	var total int64
	countDB := db.Session(&gorm.Session{})
	if err := countDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var notices []models.Notice
	if err := db.Order("sort_order DESC, created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&notices).Error; err != nil {
		return nil, 0, err
	}
	return notices, total, nil
}

// MarkNoticeRead 标记某条公告为已读（upsert）
func (s *HomeService) MarkNoticeRead(userID, noticeID string) error {
	read := models.NoticeRead{
		UserID:   userID,
		NoticeID: noticeID,
	}
	// OnConflict 走 unique (user_id, notice_id) 索引，已存在则跳过
	return config.GetDB().
		Where("user_id = ? AND notice_id = ?", userID, noticeID).
		FirstOrCreate(&read).Error
}