-- ================================================================
-- 公路理发师预约小程序 - 数据库迁移脚本
-- ================================================================

-- 启用 UUID 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ================================================================
-- 用户表 users
-- ================================================================
CREATE TABLE IF NOT EXISTS users (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    phone           VARCHAR(20)     NOT NULL UNIQUE,
    password_hash   VARCHAR(255),
    nickname        VARCHAR(50)     NOT NULL,
    avatar          VARCHAR(500),
    gender          SMALLINT        DEFAULT 0,
    role            SMALLINT        NOT NULL DEFAULT 1,
    status          SMALLINT        NOT NULL DEFAULT 1,
    last_login_at   TIMESTAMP,
    last_login_ip   VARCHAR(50),
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

COMMENT ON TABLE users IS '用户表：存储所有用户（顾客、商家、管理员）';
COMMENT ON COLUMN users.id IS '用户ID';
COMMENT ON COLUMN users.phone IS '手机号（登录账号）';
COMMENT ON COLUMN users.password_hash IS '密码哈希';
COMMENT ON COLUMN users.nickname IS '昵称';
COMMENT ON COLUMN users.avatar IS '头像URL';
COMMENT ON COLUMN users.gender IS '性别：0未知 1男 2女';
COMMENT ON COLUMN users.role IS '角色：1顾客 2商家 3管理员';
COMMENT ON COLUMN users.status IS '状态：1正常 2待审核 3禁用';
COMMENT ON COLUMN users.last_login_at IS '最后登录时间';
COMMENT ON COLUMN users.last_login_ip IS '最后登录IP';
COMMENT ON COLUMN users.created_at IS '创建时间';
COMMENT ON COLUMN users.updated_at IS '更新时间';

-- ================================================================
-- 商家扩展表 merchant_profiles
-- ================================================================
CREATE TABLE IF NOT EXISTS merchant_profiles (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID            NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    merchant_type       SMALLINT        NOT NULL,
    id_card             VARCHAR(20),
    id_card_front       VARCHAR(500),
    id_card_back        VARCHAR(500),
    business_license    VARCHAR(500),
    company_name        VARCHAR(200),
    tax_number          VARCHAR(50),
    qualification_docs  JSONB,
    audit_status        SMALLINT        NOT NULL DEFAULT 2,
    audit_remark        VARCHAR(500),
    audit_time          TIMESTAMP,
    auditor_id          UUID,
    created_at          TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_merchant_profiles_user ON merchant_profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_merchant_profiles_audit ON merchant_profiles(audit_status);

COMMENT ON TABLE merchant_profiles IS '商家扩展表：存储商家资质信息';
COMMENT ON COLUMN merchant_profiles.merchant_type IS '商家类型：1个人 2个体户 3公司';
COMMENT ON COLUMN merchant_profiles.audit_status IS '审核状态：1通过 2待审核 3拒绝';

-- ================================================================
-- 地区表 locations
-- ================================================================
CREATE TABLE IF NOT EXISTS locations (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id       UUID            REFERENCES locations(id) ON DELETE CASCADE,
    name            VARCHAR(100)    NOT NULL,
    code            VARCHAR(20)     NOT NULL UNIQUE,
    level           SMALLINT        NOT NULL,
    sort_order      INT             DEFAULT 0,
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_locations_parent ON locations(parent_id);
CREATE INDEX IF NOT EXISTS idx_locations_code ON locations(code);
CREATE INDEX IF NOT EXISTS idx_locations_level ON locations(level);

COMMENT ON TABLE locations IS '地区表：支持省市区街道小区楼栋多级下钻';
COMMENT ON COLUMN locations.parent_id IS '上级地区ID，NULL表示顶级';
COMMENT ON COLUMN locations.code IS '地区编码（如 110000）';
COMMENT ON COLUMN locations.level IS '地区级别：1省 2市 3区/县 4街道/乡镇 5村/居委会 6小区 7楼栋';

-- ================================================================
-- 店铺表 shops
-- ================================================================
CREATE TABLE IF NOT EXISTS shops (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100)    NOT NULL,
    location_id     UUID            REFERENCES locations(id),
    address         VARCHAR(255)    NOT NULL,
    longitude       DECIMAL(10,7),
    latitude        DECIMAL(10,7),
    phone           VARCHAR(20),
    business_hours  JSONB,
    images          JSONB,
    description     TEXT,
    status          SMALLINT        NOT NULL DEFAULT 1,
    creator_id      UUID,
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_shops_location ON shops(location_id);
CREATE INDEX IF NOT EXISTS idx_shops_status ON shops(status);
CREATE INDEX IF NOT EXISTS idx_shops_creator ON shops(creator_id);

COMMENT ON TABLE shops IS '店铺表：商家入驻的店铺';
COMMENT ON COLUMN shops.status IS '状态：1正常 2歇业 3停用';

-- ================================================================
-- 商家表 merchants
-- ================================================================
CREATE TABLE IF NOT EXISTS merchants (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID            NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    shop_id         UUID            REFERENCES shops(id),
    title           VARCHAR(100),
    specialties     JSONB,
    experience_years INT             DEFAULT 0,
    introduction    TEXT,
    rating          DECIMAL(2,1)    DEFAULT 5.0,
    review_count    INT             DEFAULT 0,
    service_count   INT             DEFAULT 0,
    avatar          VARCHAR(500),
    status          SMALLINT        NOT NULL DEFAULT 1,
    is_verified     BOOLEAN         DEFAULT FALSE,
    is_top          BOOLEAN         DEFAULT FALSE,
    sort_order      INT             DEFAULT 0,
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_merchants_user ON merchants(user_id);
CREATE INDEX IF NOT EXISTS idx_merchants_shop ON merchants(shop_id);
CREATE INDEX IF NOT EXISTS idx_merchants_status ON merchants(status);
CREATE INDEX IF NOT EXISTS idx_merchants_rating ON merchants(rating DESC);

COMMENT ON TABLE merchants IS '商家表：商家基本信息';
COMMENT ON COLUMN merchants.status IS '状态：1正常 2休息 3离职';
COMMENT ON COLUMN merchants.is_verified IS '是否通过资质审核';

-- ================================================================
-- 服务项目表 services
-- ================================================================
CREATE TABLE IF NOT EXISTS services (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    shop_id         UUID            NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    name            VARCHAR(100)    NOT NULL,
    description     TEXT,
    duration        INT             NOT NULL DEFAULT 60,
    price           DECIMAL(10,2)   NOT NULL DEFAULT 0,
    category        VARCHAR(50),
    images          JSONB,
    status          SMALLINT        NOT NULL DEFAULT 1,
    sort_order      INT             DEFAULT 0,
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_services_shop ON services(shop_id);
CREATE INDEX IF NOT EXISTS idx_services_category ON services(category);
CREATE INDEX IF NOT EXISTS idx_services_status ON services(status);

COMMENT ON TABLE services IS '服务项目表：店铺提供的服务';
COMMENT ON COLUMN services.duration IS '服务时长（分钟）';
COMMENT ON COLUMN services.status IS '状态：1上架 2下架';

-- ================================================================
-- 商家服务关联表 merchant_services
-- ================================================================
CREATE TABLE IF NOT EXISTS merchant_services (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID            NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    service_id      UUID            NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    price           DECIMAL(10,2),
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    UNIQUE(merchant_id, service_id)
);

CREATE INDEX IF NOT EXISTS idx_merchant_services_merchant ON merchant_services(merchant_id);
CREATE INDEX IF NOT EXISTS idx_merchant_services_service ON merchant_services(service_id);

COMMENT ON TABLE merchant_services IS '商家服务关联表';

-- ================================================================
-- 排班时段表 schedules
-- ================================================================
CREATE TABLE IF NOT EXISTS schedules (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID            NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    work_date       DATE            NOT NULL,
    start_time      VARCHAR(8)      NOT NULL,
    end_time        VARCHAR(8)      NOT NULL,
    is_available    BOOLEAN         DEFAULT TRUE,
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    UNIQUE(merchant_id, work_date, start_time)
);

CREATE INDEX IF NOT EXISTS idx_schedules_merchant ON schedules(merchant_id);
CREATE INDEX IF NOT EXISTS idx_schedules_date ON schedules(work_date);
CREATE INDEX IF NOT EXISTS idx_schedules_available ON schedules(is_available);

COMMENT ON TABLE schedules IS '排班时段表：商家发布的可用时段';

-- ================================================================
-- 预约表 bookings
-- ================================================================
CREATE TABLE IF NOT EXISTS bookings (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    order_no            VARCHAR(32)     NOT NULL UNIQUE,
    customer_id         UUID            NOT NULL REFERENCES users(id),
    merchant_id         UUID            NOT NULL REFERENCES merchants(id),
    shop_id             UUID            REFERENCES shops(id),
    service_id          UUID            NOT NULL REFERENCES services(id),
    schedule_id         UUID            NOT NULL REFERENCES schedules(id),
    appointment_date    DATE            NOT NULL,
    appointment_time    VARCHAR(8)      NOT NULL,
    duration            INT             NOT NULL,
    price               DECIMAL(10,2)   NOT NULL DEFAULT 0,
    status              SMALLINT        NOT NULL DEFAULT 1,
    cancel_reason       VARCHAR(255),
    cancel_time         TIMESTAMP,
    remark              TEXT,
    internal_note       TEXT,
    confirm_time        TIMESTAMP,
    start_time          TIMESTAMP,
    finish_time         TIMESTAMP,
    created_at          TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bookings_customer ON bookings(customer_id);
CREATE INDEX IF NOT EXISTS idx_bookings_merchant ON bookings(merchant_id);
CREATE INDEX IF NOT EXISTS idx_bookings_shop ON bookings(shop_id);
CREATE INDEX IF NOT EXISTS idx_bookings_date ON bookings(appointment_date);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
CREATE INDEX IF NOT EXISTS idx_bookings_order_no ON bookings(order_no);

COMMENT ON TABLE bookings IS '预约表：用户预约理发服务';
COMMENT ON COLUMN bookings.status IS '状态：1待确认 2已确认 3服务中 4已完成 5已取消 6已拒绝';

-- ================================================================
-- 评价表 reviews
-- ================================================================
CREATE TABLE IF NOT EXISTS reviews (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id      UUID            NOT NULL UNIQUE REFERENCES bookings(id) ON DELETE CASCADE,
    customer_id     UUID            NOT NULL REFERENCES users(id),
    merchant_id     UUID            NOT NULL REFERENCES merchants(id),
    shop_id         UUID            REFERENCES shops(id),
    service_id      UUID            REFERENCES services(id),
    rating          SMALLINT        NOT NULL,
    content         TEXT,
    images          JSONB,
    is_anonymous    BOOLEAN         DEFAULT FALSE,
    reply_content   TEXT,
    reply_time      TIMESTAMP,
    status          SMALLINT        DEFAULT 1,
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reviews_merchant ON reviews(merchant_id);
CREATE INDEX IF NOT EXISTS idx_reviews_customer ON reviews(customer_id);
CREATE INDEX IF NOT EXISTS idx_reviews_booking ON reviews(booking_id);

COMMENT ON TABLE reviews IS '评价表：顾客对服务的评价';
COMMENT ON COLUMN reviews.status IS '状态：1显示 2隐藏';

-- ================================================================
-- 商家入驻申请表 merchant_applies
-- ================================================================
CREATE TABLE IF NOT EXISTS merchant_applies (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_name      VARCHAR(100)    NOT NULL,
    applicant_phone     VARCHAR(20)     NOT NULL,
    applicant_type      SMALLINT        NOT NULL,
    id_card             VARCHAR(20),
    company_name        VARCHAR(200),
    business_license    VARCHAR(500),
    location_id         UUID            REFERENCES locations(id),
    address             VARCHAR(255),
    longitude           DECIMAL(10,7),
    latitude            DECIMAL(10,7),
    status              SMALLINT        NOT NULL DEFAULT 2,
    reject_reason       VARCHAR(500),
    audit_time          TIMESTAMP,
    auditor_id          UUID,
    created_at          TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_merchant_applies_status ON merchant_applies(status);
CREATE INDEX IF NOT EXISTS idx_merchant_applies_phone ON merchant_applies(applicant_phone);

COMMENT ON TABLE merchant_applies IS '商家入驻申请表';
COMMENT ON COLUMN merchant_applies.applicant_type IS '申请人类型：1个人 2个体户 3公司';
COMMENT ON COLUMN merchant_applies.status IS '审核状态：1通过 2待审核 3拒绝';

-- ================================================================
-- 收藏表 favorites
-- ================================================================
CREATE TABLE IF NOT EXISTS favorites (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID            NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    merchant_id     UUID            NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, merchant_id)
);

CREATE INDEX IF NOT EXISTS idx_favorites_user ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_merchant ON favorites(merchant_id);

COMMENT ON TABLE favorites IS '收藏表：顾客收藏商家';
