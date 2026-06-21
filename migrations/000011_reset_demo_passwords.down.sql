-- ================================================================
-- 000011_reset_demo_passwords.down.sql
-- 回滚：恢复旧密码哈希（admin123 / customer123 / merchant123）
-- ================================================================

UPDATE users
   SET password_hash = '$2a$10$4mmWi1bOU0F7RVnhpl42eeS6My/i2AmwI4/mcu.yJiaTgMLlWx0YO'
 WHERE id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

UPDATE users
   SET password_hash = '$2a$10$uViD7cCPtlkdVH6Ku0D/xuIXkvNbz.2k2pF8Km5a/AwI0WtRAVQLm'
 WHERE id IN (
     'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb01',
     'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb02'
     );

UPDATE users
   SET password_hash = '$2a$10$xeaTjsIRlCepEatQYwpWv.1uppgHYG/9tI34dAWy.XPWS/by4VAaW'
 WHERE id IN (
     'cccccccc-cccc-cccc-cccc-ccccccccccc1',
     'cccccccc-cccc-cccc-cccc-ccccccccccc2'
     );