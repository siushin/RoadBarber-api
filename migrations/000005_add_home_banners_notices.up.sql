-- ================================================================
-- 000005_add_home_banners_notices.up.sql
-- 首页运营内容：Banner 轮播图 + 滚动公告
-- 表结构变更（数据 seed 在 000008/000009/000010 单独跑）
-- ================================================================

-- ================================================================
-- Banner 轮播图表
-- ================================================================
CREATE TABLE IF NOT EXISTS banners (
    id           UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    image        VARCHAR(500)    NOT NULL,
    title        VARCHAR(50),
    subtitle     VARCHAR(100),
    text         VARCHAR(200),
    align        VARCHAR(20)     NOT NULL DEFAULT 'flex-start',
    link_url     VARCHAR(500),
    sort_order   INT             NOT NULL DEFAULT 0,
    status       SMALLINT        NOT NULL DEFAULT 1,
    created_at   TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_banners_status_sort ON banners(status, sort_order DESC);

COMMENT ON TABLE  banners               IS 'Banner 轮播图表：首页轮播运营位';
COMMENT ON COLUMN banners.id            IS '主键ID';
COMMENT ON COLUMN banners.image         IS '图片 URL';
COMMENT ON COLUMN banners.title         IS '主标题';
COMMENT ON COLUMN banners.subtitle      IS '副标题';
COMMENT ON COLUMN banners.text          IS '补充说明文字';
COMMENT ON COLUMN banners.align         IS '文字对齐：flex-start / center / flex-end';
COMMENT ON COLUMN banners.link_url      IS '点击跳转链接（可选）';
COMMENT ON COLUMN banners.sort_order    IS '排序权重，越大越靠前';
COMMENT ON COLUMN banners.status        IS '状态：1 启用 2 禁用';
COMMENT ON COLUMN banners.created_at    IS '创建时间';
COMMENT ON COLUMN banners.updated_at    IS '更新时间';

-- ================================================================
-- 滚动公告表
-- ================================================================
CREATE TABLE IF NOT EXISTS notices (
    id           UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    content      VARCHAR(500)    NOT NULL,
    sort_order   INT             NOT NULL DEFAULT 0,
    status       SMALLINT        NOT NULL DEFAULT 1,
    icon         VARCHAR(50),
    text_color   VARCHAR(20),
    created_at   TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notices_status_sort ON notices(status, sort_order DESC);

COMMENT ON TABLE  notices               IS '滚动公告表：首页顶部跑马灯';
COMMENT ON COLUMN notices.id            IS '主键ID';
COMMENT ON COLUMN notices.content       IS '公告内容';
COMMENT ON COLUMN notices.sort_order    IS '排序权重，越大越靠前';
COMMENT ON COLUMN notices.status        IS '状态：1 启用 2 禁用';
COMMENT ON COLUMN notices.icon         IS '图标标识：lucide 图标名（如 gift / truck / bell / info-circle-fill）';
COMMENT ON COLUMN notices.text_color   IS '文字颜色：CSS 颜色值（如 #ef4444 / #10b981 / #f59e0b）';
COMMENT ON COLUMN notices.created_at    IS '创建时间';
COMMENT ON COLUMN notices.updated_at    IS '更新时间';