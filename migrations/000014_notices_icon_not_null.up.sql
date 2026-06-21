-- ================================================================
-- 000014_notices_icon_not_null.up.sql
-- notices.icon 加 NOT NULL 约束
--
-- 背景：前端 wd-notice-bar 依赖 icon 字段渲染左侧 lucide 图标，
-- NULL 时必须用兜底 'info-circle-fill'，但数据库层面也要求 NOT NULL
-- （避免后续 seed 漏填导致前端兜底）。先把 NULL 数据回填，再加约束。
-- ================================================================

-- 1) 历史 NULL 数据回填成兜底图标
UPDATE notices SET icon = 'info-circle-fill' WHERE icon IS NULL OR icon = '';

-- 2) 加 NOT NULL 约束
ALTER TABLE notices ALTER COLUMN icon SET NOT NULL;

COMMENT ON COLUMN notices.icon IS '图标：lucide 图标名（NOT NULL，前端作为 wd-notice-bar 的 prefix 显示）';
