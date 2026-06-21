-- ================================================================
-- 000004_extend_merchants.up.sql
-- merchants 表的位置 + 4 个计算字段
--
-- 注意：以下 6 个字段已经写在 000001_init.up.sql 的 CREATE TABLE merchants
-- 中（位于 sort_order 之后、created_at 之前）。本迁移保留为幂等兜底——
-- 使用 ADD COLUMN IF NOT EXISTS + CREATE INDEX IF NOT EXISTS，在历史数据库
-- （字段被 ALTER 追加到末尾）上重复跑也不会报错。
--
-- 字段语义：
--   - latitude / longitude: 商家位置（用于距离计算）
--   - start_price: 起价（service join 计算的最低价，冗余存储避免每次聚合）
--   - business_hours: 营业时间（当日排班汇总，格式 "HH:mm - HH:mm"）
--   - distance: 距离用户位置 km（请求时按 lat/lng 计算后写入）
--   - available_slots: 当日可用时段数（按日期计算）
-- ================================================================

ALTER TABLE merchants
    ADD COLUMN IF NOT EXISTS latitude        DECIMAL(10, 7),
    ADD COLUMN IF NOT EXISTS longitude       DECIMAL(10, 7),
    ADD COLUMN IF NOT EXISTS start_price     DECIMAL(10, 2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS business_hours  VARCHAR(50),
    ADD COLUMN IF NOT EXISTS distance        DECIMAL(10, 2),
    ADD COLUMN IF NOT EXISTS available_slots INT             NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_merchants_lat_lng ON merchants(latitude, longitude);

COMMENT ON COLUMN merchants.latitude        IS '商家纬度（用于距离计算，NULL 表示未定位）';
COMMENT ON COLUMN merchants.longitude       IS '商家经度（用于距离计算，NULL 表示未定位）';
COMMENT ON COLUMN merchants.start_price     IS '起价：从 services join 计算的最低价（冗余字段，商家设置服务时回填）';
COMMENT ON COLUMN merchants.business_hours  IS '营业时间：当日排班 min(start_time)-max(end_time)';
COMMENT ON COLUMN merchants.distance        IS '距离用户位置 km（Haversine 球面距离，请求时计算）';
COMMENT ON COLUMN merchants.available_slots IS '当日可用时段数：is_available=true 的排班数';