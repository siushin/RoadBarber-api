-- ================================================================
-- 000013_seed_demo_merchants.down.sql
-- 撤销 000013：删掉演示用户（会级联删 merchants）+ 演示排班 + 演示预约
-- ================================================================

BEGIN;

-- 1) 删演示预约（schedule_id 关联的演示 schedule）
DELETE FROM bookings
 WHERE customer_id IN (
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb01',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb02'
 )
   AND remark = '演示预约'
   AND appointment_date BETWEEN CURRENT_DATE AND CURRENT_DATE + 6;

-- 2) 删演示排班（仅 +0 ~ +6 范围内，避开 000002/000010 真实数据）
DELETE FROM schedules
 WHERE work_date BETWEEN CURRENT_DATE AND CURRENT_DATE + 6
   AND start_time = '09:00'
   AND end_time = '10:00'
   AND merchant_id NOT IN (
    'ffffffff-ffff-ffff-ffff-fffffffffff1',
    'ffffffff-ffff-ffff-ffff-fffffffffff2'
   );

-- 3) 删演示商家（级联会带走 user，但 user_id ON DELETE CASCADE 是 users 上的）
--    先删 merchants，再删 users
DELETE FROM merchants
 WHERE user_id IN (
    SELECT id FROM users WHERE nickname LIKE '演示理发师 %'
 );

DELETE FROM users
 WHERE nickname LIKE '演示理发师 %';

COMMIT;
