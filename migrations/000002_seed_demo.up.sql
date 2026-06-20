-- ================================================================
-- 000002_seed_demo.up.sql
-- 演示种子数据：管理员 / 示例顾客 / 示例商家 / 示例服务 / 广东省市区
-- 密码统一：管理员 admin123 / 顾客 customer123 / 商家 merchant123
-- ================================================================

BEGIN;

-- ================================================================
-- 地区：广东省 / 深圳市 / 广州市 / 南山区（用于演示下钻）
-- ================================================================
INSERT INTO locations (id, parent_id, name, code, level, sort_order) VALUES
    ('11111111-1111-1111-1111-111111111111', NULL, '广东省', '440000', 1, 1),
    ('22222222-2222-2222-2222-222222222201', '11111111-1111-1111-1111-111111111111', '深圳市', '440300', 2, 1),
    ('22222222-2222-2222-2222-222222222202', '11111111-1111-1111-1111-111111111111', '广州市', '440100', 2, 2),
    ('33333333-3333-3333-3333-333333333301', '22222222-2222-2222-2222-222222222201', '南山区', '440305', 3, 1),
    ('33333333-3333-3333-3333-333333333302', '22222222-2222-2222-2222-222222222202', '天河区', '440106', 3, 1);

-- ================================================================
-- 用户：1 个管理员 + 2 个顾客 + 2 个商家账号
-- ================================================================
INSERT INTO users (id, phone, password_hash, nickname, role, status, gender) VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '13800000000',
     '$2a$10$4mmWi1bOU0F7RVnhpl42eeS6My/i2AmwI4/mcu.yJiaTgMLlWx0YO',
     '超级管理员', 3, 1, 1),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb01', '13900000001',
     '$2a$10$uViD7cCPtlkdVH6Ku0D/xuIXkvNbz.2k2pF8Km5a/AwI0WtRAVQLm',
     '小王', 1, 1, 1),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbb02', '13900000002',
     '$2a$10$uViD7cCPtlkdVH6Ku0D/xuIXkvNbz.2k2pF8Km5a/AwI0WtRAVQLm',
     '小李', 1, 1, 2),
    ('cccccccc-cccc-cccc-cccc-ccccccccccc1', '13900000011',
     '$2a$10$xeaTjsIRlCepEatQYwpWv.1uppgHYG/9tI34dAWy.XPWS/by4VAaW',
     'Tony 老师', 2, 1, 1),
    ('cccccccc-cccc-cccc-cccc-ccccccccccc2', '13900000012',
     '$2a$10$xeaTjsIRlCepEatQYwpWv.1uppgHYG/9tI34dAWy.XPWS/by4VAaW',
     'Kevin 总监', 2, 1, 1);

-- ================================================================
-- 商家资质（merchant_profiles）
-- ================================================================
INSERT INTO merchant_profiles (id, user_id, merchant_type, id_card, audit_status, audit_time, auditor_id) VALUES
    ('dddddddd-dddd-dddd-dddd-dddddddddd01',
     'cccccccc-cccc-cccc-cccc-ccccccccccc1',
     1, '440301199001011234', 1, NOW(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    ('dddddddd-dddd-dddd-dddd-dddddddddd02',
     'cccccccc-cccc-cccc-cccc-ccccccccccc2',
     2, '440301199002022345', 1, NOW(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa');

-- ================================================================
-- 店铺（shops）
-- ================================================================
INSERT INTO shops (id, name, location_id, address, longitude, latitude, phone, status, creator_id) VALUES
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeee01',
     'Tony 工作室（南山店）', '33333333-3333-3333-3333-333333333301',
     '科技园南路 88 号 1 楼', 113.9528000, 22.5431000, '13900000011', 1,
     'cccccccc-cccc-cccc-cccc-ccccccccccc1'),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeee02',
     'Kevin 美发沙龙（天河店）', '33333333-3333-3333-3333-333333333302',
     '天河城购物中心 B1', 113.3252000, 23.1351000, '13900000012', 1,
     'cccccccc-cccc-cccc-cccc-ccccccccccc2');

-- ================================================================
-- 商家主档（merchants）
-- ================================================================
INSERT INTO merchants (id, user_id, shop_id, title, specialties, experience_years, introduction, rating, review_count, service_count, status, is_verified, is_top) VALUES
    ('ffffffff-ffff-ffff-ffff-fffffffffff1',
     'cccccccc-cccc-cccc-cccc-ccccccccccc1',
     'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee01',
     '首席理发师 Tony',
     '["剪发","造型","染发"]'::jsonb,
     8, '擅长日韩系剪发，服务过 5000+ 客户', 4.9, 128, 1500, 1, TRUE, TRUE),
    ('ffffffff-ffff-ffff-ffff-fffffffffff2',
     'cccccccc-cccc-cccc-cccc-ccccccccccc2',
     'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee02',
     '总监 Kevin',
     '["烫发","染发","护理"]'::jsonb,
     12, '国际认证美发师，专注高端造型', 4.8, 96, 1200, 1, TRUE, FALSE);

-- ================================================================
-- 服务项目（services）
-- ================================================================
INSERT INTO services (id, shop_id, name, description, duration, price, category, status) VALUES
    ('99999999-9999-9999-9999-999999999901',
     'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee01', '精剪（男）', '基础洗剪吹', 30, 68.00, '剪发', 1),
    ('99999999-9999-9999-9999-999999999902',
     'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee01', '精剪（女）', '基础洗剪吹', 45, 98.00, '剪发', 1),
    ('99999999-9999-9999-9999-999999999903',
     'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee02', '冷烫', '进口药水冷烫', 120, 388.00, '烫发', 1);

-- ================================================================
-- 排班（schedules）：给两个商家各发布 3 个未来时段
-- ================================================================
INSERT INTO schedules (id, merchant_id, work_date, start_time, end_time, is_available) VALUES
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff1',
     CURRENT_DATE + 1, '09:00', '10:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff1',
     CURRENT_DATE + 1, '10:00', '11:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff1',
     CURRENT_DATE + 1, '14:00', '15:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff2',
     CURRENT_DATE + 1, '13:00', '15:00', TRUE),
    (gen_random_uuid(), 'ffffffff-ffff-ffff-ffff-fffffffffff2',
     CURRENT_DATE + 2, '10:00', '12:00', TRUE);

COMMIT;