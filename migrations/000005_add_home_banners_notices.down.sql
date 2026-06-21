-- ================================================================
-- 000005_add_home_banners_notices.down.sql
-- 回滚：删除 banners / notices 表
-- ================================================================

DROP TABLE IF EXISTS banners;
DROP TABLE IF EXISTS notices;