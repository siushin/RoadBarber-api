-- ================================================================
-- 000009_seed_history_banners.down.sql
-- 回滚历史 banner：只删 sort_order = 50/60/70 的 3 条 V1/V1.5/V2，
-- 不影响 000008 已写入的现代版 picsum banner (sort_order 80/90/100)
-- ================================================================

DELETE FROM banners WHERE sort_order IN (50, 60, 70);