-- ================================================================
-- 000002_seed_demo.down.sql
-- 回滚：清空演示种子数据
-- ================================================================

BEGIN;

-- 按依赖反序清理（保留管理员账号可选；这里彻底清空）
DELETE FROM schedules;
DELETE FROM merchant_services;
DELETE FROM services;
DELETE FROM favorites;
DELETE FROM reviews;
DELETE FROM bookings;
DELETE FROM merchants;
DELETE FROM shops;
DELETE FROM merchant_applies;
DELETE FROM merchant_profiles;
DELETE FROM users;
DELETE FROM locations WHERE code IN ('440000', '440300', '440100', '440305', '440106');

COMMIT;