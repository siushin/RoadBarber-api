-- ================================================================
-- 000014_notices_icon_not_null.down.sql
-- 撤销 000014：解除 notices.icon 的 NOT NULL 约束
-- ================================================================

ALTER TABLE notices ALTER COLUMN icon DROP NOT NULL;
