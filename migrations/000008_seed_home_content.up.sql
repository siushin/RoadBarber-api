-- ================================================================
-- 000008_seed_home_content.up.sql
-- 首页对接：补商家位置 + 起价 + Banner seed 数据
--
-- 与 000006_enhance_notices 的分工：
--   - 000006 只管 notices，本迁移只管 banners + merchants 字段
--   - 不动 000002_seed_demo 里的演示数据
-- ================================================================

-- ================================================================
-- 商家位置（lat/lng）+ 起价（来自 services 最低价）
-- Tony：深圳南山 22.5431, 113.9528；精剪（男）68 起
-- Kevin：广州天河 23.1351, 113.3252；冷烫 388 起
-- ================================================================
UPDATE merchants SET
    latitude    = 22.5431000,
    longitude   = 113.9528000,
    start_price = 68.00,
    shop_id     = 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee01'
WHERE id = 'ffffffff-ffff-ffff-ffff-fffffffffff1';

UPDATE merchants SET
    latitude    = 23.1351000,
    longitude   = 113.3252000,
    start_price = 388.00,
    shop_id     = 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee02'
WHERE id = 'ffffffff-ffff-ffff-ffff-fffffffffff2';

-- ================================================================
-- Banner：3 张现代版（首页轮播）
-- ================================================================
INSERT INTO banners (id, image, title, subtitle, text, align, sort_order, status) VALUES
    (gen_random_uuid(),
     'https://picsum.photos/seed/roadbarber1/750/420',
     'RoadBarber', '公路理发师', '140+ 覆盖城市数', 'flex-start', 100, 1),
    (gen_random_uuid(),
     'https://picsum.photos/seed/roadbarber2/750/420',
     'OnTheWay', '在路上为你理发', '7×24 上门服务', 'flex-end', 90, 1),
    (gen_random_uuid(),
     'https://picsum.photos/seed/roadbarber3/750/420',
     'ProBarber', '专业认证师傅', '10 万 + 用户信赖', 'center', 80, 1);