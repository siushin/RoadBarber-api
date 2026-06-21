-- ================================================================
-- 000006_enhance_notices.down.sql
-- 回滚通知结构 + 数据：清数据 → 清字段值 → 删列
-- ================================================================

DELETE FROM notices;

ALTER TABLE notices DROP COLUMN IF EXISTS text_color;
ALTER TABLE notices DROP COLUMN IF EXISTS icon;