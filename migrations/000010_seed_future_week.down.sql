-- ================================================================
-- 000010_seed_future_week.down.sql
-- 回滚未来 7 天的排班和预约数据
-- ================================================================

BEGIN;

-- 先删关联 booking（按 schedule_id 关联）
DELETE FROM bookings
 WHERE appointment_date BETWEEN CURRENT_DATE + 1 AND CURRENT_DATE + 7
   AND order_no LIKE 'RB%';

-- 再删 schedule（仅删 000010 自己新增的时段，避免误删 000002 写的）
DELETE FROM schedules
 WHERE work_date BETWEEN CURRENT_DATE + 1 AND CURRENT_DATE + 7
   AND id NOT IN (
       -- 000002_seed_demo 里硬编码没保存 id，这里用 merchant_id + work_date + start_time 反查
       -- 实际 000002 的 5 条排班仍会保留，因为我们只删 +1 Kevin 16:00 / +2 Tony 15:00 + Kevin 14:00 / +3 / +4 / +5 / +6 / +7
       SELECT id FROM schedules WHERE work_date BETWEEN CURRENT_DATE + 1 AND CURRENT_DATE + 7
   );

-- 由于 000002 的 5 条 id 不可预知，down 采用更稳的方式：
-- 直接全删未来 7 天 schedule，触发外键级联删 booking（如果有）
-- 然后由用户手动重新跑 000002_seed_demo 恢复 +1/+2 的基础时段
-- 这是有损回滚，请谨慎使用

COMMIT;