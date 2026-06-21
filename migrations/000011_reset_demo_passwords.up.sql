-- ================================================================
-- 000011_reset_demo_passwords.up.sql
-- 重置演示账号密码：管理员 admin / 其余用户 123456
-- bcrypt cost=10
-- 000002_seed_demo 里插入的是 placeholder 哈希，本迁移覆盖为可登录哈希
-- ================================================================

-- 超级管理员
UPDATE users
   SET password_hash = '$2a$10$uAkngir84e5pPbtx933Ay.monIig59tqrq8cK183JFiBzlVe1T5EC'
 WHERE id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

-- 其余用户（顾客 + 商家）
UPDATE users
   SET password_hash = '$2a$10$RKFxHM2Q7FDdghaQ02YqOegjvvCFNid3M3m5u7UV1RmuuGpGwidOq'
 WHERE id <> 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';