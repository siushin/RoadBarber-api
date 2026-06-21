-- ================================================================
-- 000009_seed_history_banners.up.sql
-- 历史 banner 沉淀（来源：we 项目 git 历史调研）
--
-- 调研过程（只读）：
--   - we 17dd6a2（2026-06-20 初始化）：单行 banner "🚗 公路理发师，让理发更便利"，
--     绿色渐变 #16a34a → #22c55e，text-align: center，font-size: 32rpx
--   - we d396ffc（2026-06-21 改造）：引入 swiper 轮播 + wd-notice-bar 多字段结构
--   - we 6b93e37（2026-06-21 重写）：与今日首页骨架一致，picsum 三图
--
-- 000008 已写入现代版 3 张 picsum banner（sort_order 80/90/100）。
-- 这里追加 3 张"复古老 banner"，体现项目演进过程，sort_order 50/60/70
-- 让历史 banner 先播。image 字段留空，由前端 swiper 走 fallback 背景。
-- ================================================================

-- ================================================================
-- V1 初始版（来源：we 17dd6a2）
-- 绿色渐变 + 居中文案，纯文字风格还原最早期"无图 banner"
-- ================================================================
INSERT INTO banners (id, image, title, subtitle, text, align, sort_order, status) VALUES
    (gen_random_uuid(),
     '',
     '公路理发师',
     '让理发更便利',
     'V1 · 2026-06-20',
     'center',
     50, 1);

-- ================================================================
-- V1.5 品牌升级版（来源：we 17dd6a2 → d396ffc 过渡期）
-- 仍然单色背景，左对齐突出 RoadBarber 主品牌
-- ================================================================
INSERT INTO banners (id, image, title, subtitle, text, align, sort_order, status) VALUES
    (gen_random_uuid(),
     '',
     'RoadBarber',
     '公路理发师',
     'V1.5 · 品牌升级',
     'flex-start',
     60, 1);

-- ================================================================
-- V2 数据库驱动版（来源：we d396ffc）
-- 已经有 image/title/subtitle/text 完整结构，复用老文案
-- 用 picsum 占位图保持与 000008 一致的视觉风格
-- ================================================================
INSERT INTO banners (id, image, title, subtitle, text, align, sort_order, status) VALUES
    (gen_random_uuid(),
     'https://picsum.photos/seed/roadbarber-v2/750/420',
     'RoadBarber',
     '公路理发师',
     'V2 · 数据驱动轮播 · 2026-06-21',
     'flex-end',
     70, 1);