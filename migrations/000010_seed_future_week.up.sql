-- ================================================================
-- 000010_seed_future_week.up.sql
-- 未来 7 天排班 + 预约种子数据（不含当天 CURRENT_DATE + 0）
--
-- 与 000002_seed_demo 的分工：
--   - 000002 已经给 +1 (Tony 09-10/10-11/14-15, Kevin 13-15) 和
--     +2 (Kevin 10-12) 种了排班。
--   - 本次迁移用 ON CONFLICT DO NOTHING 兜底（重复时段自动跳过），
--     额外补充 +1/+2 的新时段（不冲突的），并完整覆盖 +3 ~ +7 五天。
--
-- 演示商家：
--   - Tony（ffffffff-ffff-ffff-ffff-fffffffffff1）
--   - Kevin（ffffffff-ffff-ffff-ffff-fffffffffff2）
--
-- 演示顾客：
--   - 小王（bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb01）
--   - 小李（bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb02）
--
-- 排班分布（每天随机 0-4 条，部分天少于 4，体现忙闲不一）：
--   +1（已有 4 条，再补 1 条新时段）：Kevin 16-18                          = 1 条
--   +2（已有 1 条，再补 2 条新时段）：Tony 15-16 / Kevin 14-15             = 2 条
--   +3（Tony 休息）：Kevin 14-16                                          = 1 条
--   +4：Tony 09-10/14-15/15-16/16-17 + Kevin 10-12/13-14                  = 6 条
--   +5：Tony 09-10/10-11 + Kevin 13-15                                    = 3 条
--   +6（Kevin 休息）：Tony 09-10/10-11/14-15                              = 3 条
--   +7：Tony 09-10 + Kevin 10-12                                          = 2 条
--
-- Booking：从排班中抽 ~30% 转 confirmed (status=2)，覆盖核心字段。
-- 当天（CURRENT_DATE + 0）一条都不放。
-- ================================================================

BEGIN;

-- ================================================================
-- Schedules：补充 +1/+2 不冲突的时段，完整覆盖 +3 ~ +7
-- ON CONFLICT (merchant_id, work_date, start_time) DO NOTHING 兜底重复
-- ================================================================

-- +1 补充：Kevin 16:00-18:00（与已有 13:00-15:00 不冲突）
INSERT INTO schedules (id, merchant_id, work_date, start_time, end_time, is_available) VALUES
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff2', CURRENT_DATE + 1, '16:00', '18:00', TRUE)
ON CONFLICT (merchant_id, work_date, start_time) DO NOTHING;

-- +2 补充：Tony 15:00-16:00 + Kevin 14:00-15:00（与已有 Kevin 10:00-12:00 不冲突）
INSERT INTO schedules (id, merchant_id, work_date, start_time, end_time, is_available) VALUES
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff1', CURRENT_DATE + 2, '15:00', '16:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff2', CURRENT_DATE + 2, '14:00', '15:00', TRUE)
ON CONFLICT (merchant_id, work_date, start_time) DO NOTHING;

-- +3（Tony 休息）
INSERT INTO schedules (id, merchant_id, work_date, start_time, end_time, is_available) VALUES
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff2', CURRENT_DATE + 3, '14:00', '16:00', TRUE)
ON CONFLICT (merchant_id, work_date, start_time) DO NOTHING;

-- +4
INSERT INTO schedules (id, merchant_id, work_date, start_time, end_time, is_available) VALUES
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff1', CURRENT_DATE + 4, '09:00', '10:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff1', CURRENT_DATE + 4, '14:00', '15:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff1', CURRENT_DATE + 4, '15:00', '16:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff1', CURRENT_DATE + 4, '16:00', '17:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff2', CURRENT_DATE + 4, '10:00', '12:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff2', CURRENT_DATE + 4, '13:00', '14:00', TRUE)
ON CONFLICT (merchant_id, work_date, start_time) DO NOTHING;

-- +5
INSERT INTO schedules (id, merchant_id, work_date, start_time, end_time, is_available) VALUES
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff1', CURRENT_DATE + 5, '09:00', '10:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff1', CURRENT_DATE + 5, '10:00', '11:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff2', CURRENT_DATE + 5, '13:00', '15:00', TRUE)
ON CONFLICT (merchant_id, work_date, start_time) DO NOTHING;

-- +6（Kevin 休息）
INSERT INTO schedules (id, merchant_id, work_date, start_time, end_time, is_available) VALUES
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff1', CURRENT_DATE + 6, '09:00', '10:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff1', CURRENT_DATE + 6, '10:00', '11:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff1', CURRENT_DATE + 6, '14:00', '15:00', TRUE)
ON CONFLICT (merchant_id, work_date, start_time) DO NOTHING;

-- +7
INSERT INTO schedules (id, merchant_id, work_date, start_time, end_time, is_available) VALUES
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff1', CURRENT_DATE + 7, '09:00', '10:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff2', CURRENT_DATE + 7, '10:00', '12:00', TRUE)
ON CONFLICT (merchant_id, work_date, start_time) DO NOTHING;

-- ================================================================
-- Bookings：把 000010 新增时段（约 18 条）中 ~30% 转 confirmed
-- 用 CTE 锁定"本次新增 schedule"：通过 order_no 还没生成（不存在关联 booking）来识别。
-- 更稳的做法：直接从 schedules 表筛未关联 booking 的时段。
-- ================================================================
WITH s AS (
    SELECT sch.id, sch.merchant_id, sch.work_date, sch.start_time,
           (sch.end_time::time - sch.start_time::time) AS dur
      FROM schedules sch
     WHERE sch.work_date BETWEEN CURRENT_DATE + 1 AND CURRENT_DATE + 7
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
    (CASE WHEN merchant_id = 'ffffffff-ffff-ffff-ffff-fffffffffff1'
         THEN 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee01'
         ELSE 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee02'
    END)::uuid,
    (CASE WHEN merchant_id = 'ffffffff-ffff-ffff-ffff-fffffffffff1'
         THEN '99999999-9999-9999-9999-999999999901'
         ELSE '99999999-9999-9999-9999-999999999903'
    END)::uuid,
    id,
    work_date,
    start_time,
    EXTRACT(EPOCH FROM dur) / 60,
    CASE WHEN merchant_id = 'ffffffff-ffff-ffff-ffff-fffffffffff1'
         THEN 68.00
         ELSE 388.00
    END,
    2,
    '演示预约',
    NOW()
  FROM s
 WHERE RANDOM() < 0.30;

COMMIT;