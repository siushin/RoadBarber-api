-- ================================================================
-- 000008_seed_home_content.down.sql
-- 回滚：清空 banners，回退 merchant lat/lng/start_price
-- ================================================================

DELETE FROM banners;

UPDATE merchants SET
    latitude    = NULL,
    longitude   = NULL,
    start_price = 0
WHERE id IN (
    'ffffffff-ffff-ffff-ffff-fffffffffff1',
    'ffffffff-ffff-ffff-ffff-fffffffffff2'
);