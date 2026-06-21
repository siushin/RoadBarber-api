-- ================================================================
-- 000004_extend_merchants.down.sql
-- 回滚：删除 merchants 表扩展字段
-- ================================================================

DROP INDEX IF EXISTS idx_merchants_lat_lng;

ALTER TABLE merchants
    DROP COLUMN IF EXISTS latitude,
    DROP COLUMN IF EXISTS longitude,
    DROP COLUMN IF EXISTS start_price,
    DROP COLUMN IF EXISTS business_hours,
    DROP COLUMN IF EXISTS distance,
    DROP COLUMN IF EXISTS available_slots;