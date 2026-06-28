-- ================================================================
-- 000013_seed_demo_merchants.up.sql
-- 演示用 20 个假商家 + 未来 7 天（+0 ~ +6）排班 + 预约种子数据
--
-- 用途：让首页 merchants 列表每页 10 条、3 页以上，"多页"效果。
-- 与 000010 的关系：000010 已经给 Tony/Kevin 在 +1 ~ +7 铺了排班，
-- 000013 给所有 22 个商家（Tony/Kevin + 20 假）补全 +0 ~ +6 的排班。
--
-- 排班分布规则：
--   - +0 / +1 / +2 / +4 / +5 / +6：22 个商家每天 1 个 slot  → 商家列表 22 条 = 3 页（10/页）
--   - +3：故意只让 2 个商家有 slot                       → 商家列表 < 4 条（7 天中唯一一天低）
--
-- 预约：~30% 转 confirmed (status=2)，覆盖核心字段。
-- ================================================================

BEGIN;

-- ================================================================
-- 1) 演示用户：20 个 role=2 老板，复用同一个 shop_id (Tony 工作室) 即可
-- ================================================================
-- phone 用 '1390009000X' 占位，X 跟 g 走，确保 UNIQUE NOT NULL 满足
-- password_hash 留空（演示数据，无须登录）
INSERT INTO users (id, phone, nickname, role, status, created_at, updated_at)
SELECT gen_random_uuid(), '1390009000' || LPAD(g::text, 2, '0'), '演示理发师 ' || g, 2, 1, NOW(), NOW()
  FROM generate_series(1, 20) g
ON CONFLICT DO NOTHING;

-- ================================================================
-- 2) 演示商家：20 个 merchants
--    - user_id 引用刚插入的 20 个演示用户
--    - shop_id 复用 Tony 工作室（不强制外键，可空）
--    - 经纬度在广州天河 / 越秀 / 海珠 / 番禺 / 白云 / 荔湾分散，避免全部重叠
--    - rating 4.4-4.9、experience_years 2-15、price 50-388
-- ================================================================
INSERT INTO merchants (
    id, user_id, shop_id, title, specialties, experience_years, introduction,
    rating, review_count, service_count, avatar, status, is_verified, is_top,
    sort_order, latitude, longitude, created_at, updated_at
)
SELECT
    gen_random_uuid(),
    u.id,
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee01',
    '演示理发师 ' || g,
    jsonb_build_array('剪发', '造型'),
    (2 + (g % 14))::bigint,
    '专注剪发造型，服务专业，价格实惠。',
    (4.4 + (g % 6) * 0.1)::numeric,
    ((g * 37) % 200 + 50)::bigint,
    ((g * 3) % 8 + 3)::bigint,
    'https://picsum.photos/seed/roadbarber-merchant-' || g || '/400/400',
    1,            -- MerchantStatusNormal
    TRUE,
    FALSE,
    (100 - g)::bigint,    -- sort_order 从 99 倒排，让新商家排在 Tony/Kevin 后
    22.5 + (g * 0.013)::numeric,           -- 22.5 ~ 22.747
    113.3 + (g * 0.018)::numeric,          -- 113.3 ~ 113.642
    NOW(),
    NOW()
  FROM users u, generate_series(1, 20) g
 WHERE u.nickname = '演示理发师 ' || g
   AND u.phone = '1390009000' || LPAD(g::text, 2, '0')
ON CONFLICT (user_id) DO NOTHING;

-- ================================================================
-- 3) 排班：+0 / +1 / +2 / +4 / +5 / +6 给 22 个商家每天 1 个 slot（09:00-10:00）
--    +3 故意只给前 2 个商家（Tony + Kevin）铺 1 个 slot
--    ON CONFLICT (merchant_id, work_date, start_time) DO NOTHING 兜底
-- ============================================================      ===+++

-- 5 天满铺：+0 / +1 / +2 / +4 / +5 / +6
INSERT INTO schedules (id, merchant_id, work_date, start_time, end_time, is_available)
SELECT
    gen_random_uuid(),
    m.id,
    d.day::date,
    '09:00',
    '10:00',
    TRUE
  FROM merchants m
 CROSS JOIN (VALUES
    (CURRENT_DATE + 0),
    (CURRENT_DATE + 1),
    (CURRENT_DATE + 2),
    (CURRENT_DATE + 4),
    (CURRENT_DATE + 5),
    (CURRENT_DATE + 6)
 ) AS d(day)
 WHERE m.status = 1
   AND d.day <> (CURRENT_DATE + 3)
ON CONFLICT (merchant_id, work_date, start_time) DO NOTHING;

-- +3 只让 2 个商家有 slot
INSERT INTO schedules (id, merchant_id, work_date, start_time, end_time, is_available)
SELECT
    gen_random_uuid(),
    m.id,
    CURRENT_DATE + 3,
    '09:00',
    '10:00',
    TRUE
  FROM merchants m
 WHERE m.id IN (
    'ffffffff-ffff-ffff-ffff-fffffffffff1',  -- Tony
    'ffffffff-ffff-ffff-ffff-fffffffffff2'   -- Kevin
 )
ON CONFLICT (merchant_id, work_date, start_time) DO NOTHING;

-- ================================================================
-- 4) 预约：~30% 转为 confirmed (status=2)
--    锁定本次新增 schedule：work_date BETWEEN +0 ~ +6 且没关联 booking
--    customer_id 随机分配小王 / 小李
-- ============================================================
WITH s AS (
    SELECT sch.id, sch.merchant_id, sch.work_date, sch.start_time,
           (sch.end_time::time - sch.start_time::time) AS dur
      FROM schedules sch
     WHERE sch.work_date BETWEEN CURRENT_DATE + 0 AND CURRENT_DATE + 6
       AND NOT EXISTS (
           SELECT 1 FROM bookings b WHERE b.schedule_id = sch.id
       )
)
INSERT INTO bookings (
    order_no, customer_id, merchant_id, shop_id, service_id, schedule_id,
    appointment_date, appointment_time, duration, price, status, remark, confirm_time
)
SELECT
    'RB' || TO_CHAR(work_date, 'YYYYMMDD') || LPAD(FLOOR(RANDOM() * 10000)::TEXT, 4, '0'),
    (CASE WHEN RANDOM() < 0.5
         THEN 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb01'
         ELSE 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb02'
    END)::uuid,
    merchant_id,
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee01',
    (ARRAY[
        '99999999-9999-9999-9999-999999999901',
        '99999999-9999-9999-9999-999999999902',
        '99999999-9999-9999-9999-999999999903'
    ])[FLOOR(RANDOM() * 3)::int + 1]::uuid,
    id,
    work_date,
    start_time,
    EXTRACT(EPOCH FROM dur) / 60,
    (ARRAY[68.00, 98.00, 388.00])[FLOOR(RANDOM() * 3)::int + 1],
    2,
    '演示预约',
    NOW()
  FROM s
 WHERE RANDOM() < 0.30;

COMMIT;
