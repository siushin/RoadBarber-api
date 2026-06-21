-- ================================================================
-- 000006_enhance_notices.up.sql
-- notices 表的 10 条公告 seed（结构已在 000005 建好）
--
-- 注意：icon / text_color 字段已经在 000005_add_home_banners_notices.up.sql
-- 的 CREATE TABLE notices 中建好（位于 status 之后、created_at 之前）。
-- 本迁移只补 COMMENT ON COLUMN 兜底 + 10 条 seed 数据。
--
-- icon 存 lucide 图标名（gift / truck / bell / star 等），
-- 前端作为 wd-notice-bar 的 prefix 显示在通知栏左侧。
-- text_color 是 CSS 颜色值，前端用作 wd-notice-bar 的 color prop 控制文字色。
-- content 是纯文本（不再带 emoji 前缀，emoji 由 icon 字段承担）。
--
-- 数据：
--   - 4 条原始公告
--   - 6 条扩展公告
--
-- 注意：本迁移是当前 notices 表的唯一数据源。
-- 若要在 000008_seed_home_content 之前扩充 banners，需要先跑本迁移。
-- ================================================================

COMMENT ON COLUMN notices.icon       IS '图标：lucide 图标名（如 gift / truck / bell / info-circle-fill），作为 wd-notice-bar 的 prefix 显示';
COMMENT ON COLUMN notices.text_color IS '文字颜色：CSS 颜色值（如 #ef4444 / #10b981 / #f59e0b），作为 wd-notice-bar 的 color prop';

-- ================================================================
-- 4 条原始公告（与 000008_seed_home_content 中的 banners 同一时期，UI 风格一致）
-- ================================================================
INSERT INTO notices (id, content, icon, text_color, sort_order, status) VALUES
    (gen_random_uuid(), '新用户首单立减 20 元',                'gift',  '#ef4444', 100, 1),
    (gen_random_uuid(), '公路理发师已覆盖 140+ 城市',          'truck', '#10b981', 90, 1),
    (gen_random_uuid(), '累计服务超 10 万次',                  'star',  '#f59e0b', 80, 1),
    (gen_random_uuid(), '邀请好友得 30 元代金券',              'gift',  '#8b5cf6', 70, 1);

-- ================================================================
-- 6 条扩展公告（不同图标 + 颜色）
-- ================================================================
INSERT INTO notices (id, content, icon, text_color, sort_order, status) VALUES
    (gen_random_uuid(), '618 大促 · 全场 8 折',                'gift',  '#ec4899', 95, 1),
    (gen_random_uuid(), '退换货政策升级，售后无忧',            'info-circle-fill', '#3b82f6', 85, 1),
    (gen_random_uuid(), '师傅评选开启，投出你心中的 TOP1',     'heart', '#ef4444', 75, 1),
    (gen_random_uuid(), 'VIP 会员上线，享 9 折专属价',         'diamond', '#f59e0b', 65, 1),
    (gen_random_uuid(), '夏季空调房护理 · 限时 9.9 元起',     'flame', '#f97316', 55, 1),
    (gen_random_uuid(), '凌晨预约通道已开放 · 0-6 点专属优惠', 'rocket', '#06b6d4', 45, 1);