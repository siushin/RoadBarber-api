-- ================================================================
-- 000001_init.up.sql
-- 公路理发师预约小程序 - 初始化所有核心表
--
-- 表按重要性排序：
--   1. users            — 用户（认证/身份，最重要）
--   2. locations        — 地区（基础字典）
--   3. merchant_profiles — 商家资质（依赖 users）
--   4. shops            — 店铺（依赖 locations）
--   5. merchants        — 商家主档（依赖 users / shops）
--   6. services         — 服务项目（依赖 shops）
--   7. merchant_services — 商家-服务多对多
--   8. schedules        — 排班时段（依赖 merchants）
--   9. bookings         — 预约（依赖 users / merchants / schedules / services）
--  10. reviews          — 评价（依赖 bookings）
--  11. favorites        — 收藏（依赖 users / merchants）
--  12. merchant_applies  — 商家入驻申请（独立）
--
-- 注意：merchants 表的位置 / 起价 / 营业时间等扩展字段在
-- 000004_extend_merchants 中通过 ALTER TABLE 单独加，
-- 本迁移只保留 init 时的最小字段集。
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
    username        VARCHAR(50),
    email           VARCHAR(255),
    openid          VARCHAR(128),
    unionid         VARCHAR(128),
    wx_nickname     VARCHAR(100),
    wx_avatar       VARCHAR(500),
    last_login_at   TIMESTAMP,
    last_login_ip   VARCHAR(50),
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_unionid ON users(unionid);

CREATE UNIQUE INDEX IF NOT EXISTS uq_users_username ON users(username) WHERE username IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_users_email    ON users(email)    WHERE email    IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_users_openid   ON users(openid)   WHERE openid   IS NOT NULL;

COMMENT ON TABLE  users              IS '用户表：存储所有用户（顾客、商家、管理员）';
COMMENT ON COLUMN users.id            IS '用户ID';
COMMENT ON COLUMN users.phone         IS '手机号（登录账号）';
COMMENT ON COLUMN users.password_hash IS '密码哈希';
COMMENT ON COLUMN users.nickname      IS '昵称';
COMMENT ON COLUMN users.avatar        IS '头像URL';
COMMENT ON COLUMN users.gender        IS '性别：0未知 1男 2女';
COMMENT ON COLUMN users.role          IS '角色：1顾客 2商家 3管理员';
COMMENT ON COLUMN users.status        IS '状态：1正常 2待审核 3禁用';
COMMENT ON COLUMN users.username      IS '用户名（后台登录账号，全局唯一）';
COMMENT ON COLUMN users.email         IS '邮箱（可作为登录账号）';
COMMENT ON COLUMN users.openid        IS '微信小程序 openid（用于微信登录）';
COMMENT ON COLUMN users.unionid       IS '微信开放平台 unionid（多端打通）';
COMMENT ON COLUMN users.wx_nickname   IS '微信昵称';
COMMENT ON COLUMN users.wx_avatar     IS '微信头像 URL';
COMMENT ON COLUMN users.last_login_at IS '最后登录时间';
COMMENT ON COLUMN users.last_login_ip IS '最后登录IP';
COMMENT ON COLUMN users.created_at    IS '创建时间';
COMMENT ON COLUMN users.updated_at    IS '更新时间';

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

COMMENT ON TABLE  merchant_profiles                 IS '商家扩展表：存储商家资质信息（与 users 一对一）';
COMMENT ON COLUMN merchant_profiles.id              IS '主键ID';
COMMENT ON COLUMN merchant_profiles.user_id         IS '关联用户ID（外键 users.id）';
COMMENT ON COLUMN merchant_profiles.merchant_type   IS '商家类型：1个人 2个体户 3公司';
COMMENT ON COLUMN merchant_profiles.id_card         IS '身份证号';
COMMENT ON COLUMN merchant_profiles.id_card_front   IS '身份证正面照片URL';
COMMENT ON COLUMN merchant_profiles.id_card_back    IS '身份证反面照片URL';
COMMENT ON COLUMN merchant_profiles.business_license IS '营业执照照片URL（个体户/公司）';
COMMENT ON COLUMN merchant_profiles.company_name    IS '公司名称';
COMMENT ON COLUMN merchant_profiles.tax_number      IS '税务登记号';
COMMENT ON COLUMN merchant_profiles.qualification_docs IS '其他资质证书列表（JSONB 数组）';
COMMENT ON COLUMN merchant_profiles.audit_status    IS '审核状态：1通过 2待审核 3拒绝';
COMMENT ON COLUMN merchant_profiles.audit_remark    IS '审核备注';
COMMENT ON COLUMN merchant_profiles.audit_time      IS '审核时间';
COMMENT ON COLUMN merchant_profiles.auditor_id      IS '审核人ID（外键 users.id）';
COMMENT ON COLUMN merchant_profiles.created_at      IS '创建时间';
COMMENT ON COLUMN merchant_profiles.updated_at      IS '更新时间';

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
CREATE INDEX IF NOT EXISTS idx_locations_level  ON locations(level);

COMMENT ON TABLE  locations              IS '地区表：省/市/区三级行政区划';
COMMENT ON COLUMN locations.id            IS '主键ID';
COMMENT ON COLUMN locations.parent_id     IS '父级ID（外键 locations.id，根节点为 NULL）';
COMMENT ON COLUMN locations.name          IS '地区名称';
COMMENT ON COLUMN locations.code          IS '行政区划代码（国标）';
COMMENT ON COLUMN locations.level         IS '级别：1省 2市 3区';
COMMENT ON COLUMN locations.sort_order    IS '排序权重';
COMMENT ON COLUMN locations.created_at    IS '创建时间';
COMMENT ON COLUMN locations.updated_at    IS '更新时间';

-- ================================================================
-- 店铺表 shops
-- ================================================================
CREATE TABLE IF NOT EXISTS shops (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100)    NOT NULL,
    location_id     UUID            REFERENCES locations(id),
    address         VARCHAR(500),
    longitude       DECIMAL(10, 7),
    latitude        DECIMAL(10, 7),
    phone           VARCHAR(20),
    cover           VARCHAR(500),
    description     TEXT,
    business_hours  VARCHAR(50),
    status          SMALLINT        NOT NULL DEFAULT 1,
    creator_id      UUID            REFERENCES users(id),
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_shops_location ON shops(location_id);
CREATE INDEX IF NOT EXISTS idx_shops_status   ON shops(status);
CREATE INDEX IF NOT EXISTS idx_shops_creator  ON shops(creator_id);

COMMENT ON TABLE  shops                IS '店铺表：商家经营的实体门店';
COMMENT ON COLUMN shops.id              IS '主键ID';
COMMENT ON COLUMN shops.name            IS '店铺名称';
COMMENT ON COLUMN shops.location_id     IS '所属地区（外键 locations.id）';
COMMENT ON COLUMN shops.address         IS '详细地址';
COMMENT ON COLUMN shops.longitude       IS '经度';
COMMENT ON COLUMN shops.latitude        IS '纬度';
COMMENT ON COLUMN shops.phone           IS '联系电话';
COMMENT ON COLUMN shops.cover           IS '店铺封面图 URL';
COMMENT ON COLUMN shops.description     IS '店铺简介';
COMMENT ON COLUMN shops.business_hours  IS '营业时间（格式 "HH:mm - HH:mm"）';
COMMENT ON COLUMN shops.status          IS '状态：1营业 2休息 3关闭';
COMMENT ON COLUMN shops.creator_id      IS '创建人ID（外键 users.id）';
COMMENT ON COLUMN shops.created_at      IS '创建时间';
COMMENT ON COLUMN shops.updated_at      IS '更新时间';

-- ================================================================
-- 商家主档表 merchants
-- ================================================================
CREATE TABLE IF NOT EXISTS merchants (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID            NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    shop_id         UUID            REFERENCES shops(id),
    title           VARCHAR(100),
    specialties     JSONB,
    experience_years INT            DEFAULT 0,
    introduction    TEXT,
    rating          DECIMAL(2,1)    DEFAULT 5.0,
    review_count    INT             DEFAULT 0,
    service_count   INT             DEFAULT 0,
    avatar          VARCHAR(500),
    status          SMALLINT        NOT NULL DEFAULT 1,
    is_verified     BOOLEAN         DEFAULT FALSE,
    is_top          BOOLEAN         DEFAULT FALSE,
    sort_order      INT             DEFAULT 0,
    latitude        DECIMAL(10, 7),
    longitude       DECIMAL(10, 7),
    start_price     DECIMAL(10, 2)  NOT NULL DEFAULT 0,
    business_hours  VARCHAR(50),
    distance        DECIMAL(10, 2),
    available_slots INT             NOT NULL DEFAULT 0,
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_merchants_user    ON merchants(user_id);
CREATE INDEX IF NOT EXISTS idx_merchants_shop    ON merchants(shop_id);
CREATE INDEX IF NOT EXISTS idx_merchants_status  ON merchants(status);
CREATE INDEX IF NOT EXISTS idx_merchants_rating  ON merchants(rating DESC);
CREATE INDEX IF NOT EXISTS idx_merchants_top     ON merchants(is_top, sort_order DESC);
CREATE INDEX IF NOT EXISTS idx_merchants_lat_lng ON merchants(latitude, longitude);

COMMENT ON TABLE  merchants                IS '商家主档表：理发师信息（与 users 一对一）';
COMMENT ON COLUMN merchants.id              IS '主键ID';
COMMENT ON COLUMN merchants.user_id         IS '关联用户ID（外键 users.id）';
COMMENT ON COLUMN merchants.shop_id         IS '所属店铺（外键 shops.id）';
COMMENT ON COLUMN merchants.title           IS '商家职称（如"首席理发师 Tony"）';
COMMENT ON COLUMN merchants.specialties     IS '擅长项目 JSONB 数组（如 ["剪发","造型"]）';
COMMENT ON COLUMN merchants.experience_years IS '从业年限';
COMMENT ON COLUMN merchants.introduction    IS '个人介绍';
COMMENT ON COLUMN merchants.rating          IS '综合评分（0-5）';
COMMENT ON COLUMN merchants.review_count    IS '评价总数';
COMMENT ON COLUMN merchants.service_count   IS '服务人次';
COMMENT ON COLUMN merchants.avatar          IS '头像 URL';
COMMENT ON COLUMN merchants.status          IS '状态：1正常 2休息 3离职';
COMMENT ON COLUMN merchants.is_verified     IS '是否平台认证';
COMMENT ON COLUMN merchants.is_top          IS '是否置顶推荐';
COMMENT ON COLUMN merchants.sort_order      IS '排序权重';
COMMENT ON COLUMN merchants.latitude         IS '商家纬度（用于距离计算，NULL 表示未定位）';
COMMENT ON COLUMN merchants.longitude        IS '商家经度（用于距离计算，NULL 表示未定位）';
COMMENT ON COLUMN merchants.start_price      IS '起价：从 services join 计算的最低价（冗余字段，商家设置服务时回填）';
COMMENT ON COLUMN merchants.business_hours   IS '营业时间：当日排班 min(start_time)-max(end_time)';
COMMENT ON COLUMN merchants.distance         IS '距离用户位置 km（Haversine 球面距离，请求时计算）';
COMMENT ON COLUMN merchants.available_slots  IS '当日可用时段数：is_available=true 的排班数';
COMMENT ON COLUMN merchants.created_at      IS '创建时间';
COMMENT ON COLUMN merchants.updated_at      IS '更新时间';

-- ================================================================
-- 服务项目表 services
-- ================================================================
CREATE TABLE IF NOT EXISTS services (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    shop_id         UUID            NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    name            VARCHAR(100)    NOT NULL,
    description     TEXT,
    duration        INT             NOT NULL,
    price           DECIMAL(10,2)   NOT NULL DEFAULT 0,
    category        VARCHAR(50),
    cover           VARCHAR(500),
    status          SMALLINT        NOT NULL DEFAULT 1,
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_services_shop    ON services(shop_id);
CREATE INDEX IF NOT EXISTS idx_services_status  ON services(status);
CREATE INDEX IF NOT EXISTS idx_services_cat     ON services(category);

COMMENT ON TABLE  services              IS '服务项目表：店铺提供的服务';
COMMENT ON COLUMN services.id            IS '主键ID';
COMMENT ON COLUMN services.shop_id       IS '所属店铺（外键 shops.id）';
COMMENT ON COLUMN services.name          IS '服务名称';
COMMENT ON COLUMN services.description   IS '服务描述';
COMMENT ON COLUMN services.duration      IS '服务时长（分钟）';
COMMENT ON COLUMN services.price         IS '价格（元）';
COMMENT ON COLUMN services.category      IS '服务类别（如"剪发""烫发"）';
COMMENT ON COLUMN services.cover         IS '服务封面图 URL';
COMMENT ON COLUMN services.status        IS '状态：1上架 2下架';
COMMENT ON COLUMN services.created_at    IS '创建时间';
COMMENT ON COLUMN services.updated_at    IS '更新时间';

-- ================================================================
-- 商家-服务关联表 merchant_services（多对多）
-- ================================================================
CREATE TABLE IF NOT EXISTS merchant_services (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id     UUID            NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    service_id      UUID            NOT NULL REFERENCES services(id)  ON DELETE CASCADE,
    price           DECIMAL(10,2),
    status          SMALLINT        NOT NULL DEFAULT 1,
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    UNIQUE (merchant_id, service_id)
);

CREATE INDEX IF NOT EXISTS idx_merchant_services_merchant ON merchant_services(merchant_id);
CREATE INDEX IF NOT EXISTS idx_merchant_services_service  ON merchant_services(service_id);

COMMENT ON TABLE  merchant_services                  IS '商家-服务多对多关联表';
COMMENT ON COLUMN merchant_services.id               IS '主键ID';
COMMENT ON COLUMN merchant_services.merchant_id      IS '商家ID（外键 merchants.id）';
COMMENT ON COLUMN merchant_services.service_id       IS '服务ID（外键 services.id）';
COMMENT ON COLUMN merchant_services.price            IS '商家定制价（可覆盖 services.price）';
COMMENT ON COLUMN merchant_services.status           IS '状态：1提供 2停供';
COMMENT ON COLUMN merchant_services.created_at       IS '创建时间';
COMMENT ON COLUMN merchant_services.updated_at       IS '更新时间';

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
    UNIQUE (merchant_id, work_date, start_time)
);

CREATE INDEX IF NOT EXISTS idx_schedules_merchant  ON schedules(merchant_id);
CREATE INDEX IF NOT EXISTS idx_schedules_date      ON schedules(work_date);
CREATE INDEX IF NOT EXISTS idx_schedules_available ON schedules(is_available);

COMMENT ON TABLE  schedules                IS '排班时段表：商家发布的可用时段';
COMMENT ON COLUMN schedules.id              IS '主键ID';
COMMENT ON COLUMN schedules.merchant_id     IS '商家ID（外键 merchants.id）';
COMMENT ON COLUMN schedules.work_date       IS '工作日期';
COMMENT ON COLUMN schedules.start_time      IS '开始时间（HH:mm）';
COMMENT ON COLUMN schedules.end_time        IS '结束时间（HH:mm）';
COMMENT ON COLUMN schedules.is_available    IS '是否可预约';
COMMENT ON COLUMN schedules.created_at      IS '创建时间';
COMMENT ON COLUMN schedules.updated_at      IS '更新时间';

-- ================================================================
-- 预约表 bookings
-- ================================================================
CREATE TABLE IF NOT EXISTS bookings (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    order_no            VARCHAR(32)     NOT NULL UNIQUE,
    customer_id         UUID            NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    merchant_id         UUID            NOT NULL REFERENCES merchants(id) ON DELETE RESTRICT,
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

CREATE INDEX IF NOT EXISTS idx_bookings_order_no ON bookings(order_no);
CREATE INDEX IF NOT EXISTS idx_bookings_customer ON bookings(customer_id);
CREATE INDEX IF NOT EXISTS idx_bookings_merchant ON bookings(merchant_id);
CREATE INDEX IF NOT EXISTS idx_bookings_shop     ON bookings(shop_id);
CREATE INDEX IF NOT EXISTS idx_bookings_date     ON bookings(appointment_date);
CREATE INDEX IF NOT EXISTS idx_bookings_status   ON bookings(status);

COMMENT ON TABLE  bookings                  IS '预约表：用户预约理发服务';
COMMENT ON COLUMN bookings.id                IS '主键ID';
COMMENT ON COLUMN bookings.order_no          IS '订单号（业务唯一）';
COMMENT ON COLUMN bookings.customer_id       IS '顾客ID（外键 users.id）';
COMMENT ON COLUMN bookings.merchant_id       IS '商家ID（外键 merchants.id）';
COMMENT ON COLUMN bookings.shop_id           IS '店铺ID（外键 shops.id，可选）';
COMMENT ON COLUMN bookings.service_id        IS '服务ID（外键 services.id）';
COMMENT ON COLUMN bookings.schedule_id       IS '排班ID（外键 schedules.id）';
COMMENT ON COLUMN bookings.appointment_date  IS '预约日期';
COMMENT ON COLUMN bookings.appointment_time  IS '预约时间（HH:mm）';
COMMENT ON COLUMN bookings.duration          IS '服务时长（分钟）';
COMMENT ON COLUMN bookings.price             IS '实付金额（元）';
COMMENT ON COLUMN bookings.status            IS '状态：1待确认 2已确认 3服务中 4已完成 5已取消 6已拒绝';
COMMENT ON COLUMN bookings.cancel_reason     IS '取消原因';
COMMENT ON COLUMN bookings.cancel_time       IS '取消时间';
COMMENT ON COLUMN bookings.remark            IS '顾客备注';
COMMENT ON COLUMN bookings.internal_note     IS '商家内部备注';
COMMENT ON COLUMN bookings.confirm_time      IS '确认时间';
COMMENT ON COLUMN bookings.start_time        IS '开始服务时间';
COMMENT ON COLUMN bookings.finish_time       IS '完成时间';
COMMENT ON COLUMN bookings.created_at        IS '创建时间';
COMMENT ON COLUMN bookings.updated_at        IS '更新时间';

-- ================================================================
-- 评价表 reviews
-- ================================================================
CREATE TABLE IF NOT EXISTS reviews (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id      UUID            NOT NULL UNIQUE REFERENCES bookings(id) ON DELETE CASCADE,
    customer_id     UUID            NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    merchant_id     UUID            NOT NULL REFERENCES merchants(id) ON DELETE RESTRICT,
    rating          SMALLINT        NOT NULL,
    content         TEXT,
    images          JSONB,
    is_anonymous    BOOLEAN         DEFAULT FALSE,
    merchant_reply  TEXT,
    reply_time      TIMESTAMP,
    status          SMALLINT        NOT NULL DEFAULT 1,
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reviews_booking  ON reviews(booking_id);
CREATE INDEX IF NOT EXISTS idx_reviews_customer ON reviews(customer_id);
CREATE INDEX IF NOT EXISTS idx_reviews_merchant ON reviews(merchant_id);
CREATE INDEX IF NOT EXISTS idx_reviews_rating   ON reviews(rating);

COMMENT ON TABLE  reviews                IS '评价表：顾客对预约服务的评价';
COMMENT ON COLUMN reviews.id              IS '主键ID';
COMMENT ON COLUMN reviews.booking_id      IS '预约ID（外键 bookings.id）';
COMMENT ON COLUMN reviews.customer_id     IS '顾客ID（外键 users.id）';
COMMENT ON COLUMN reviews.merchant_id     IS '商家ID（外键 merchants.id）';
COMMENT ON COLUMN reviews.rating          IS '评分 1-5';
COMMENT ON COLUMN reviews.content         IS '评价内容';
COMMENT ON COLUMN reviews.images          IS '评价图片 JSONB 数组';
COMMENT ON COLUMN reviews.is_anonymous    IS '是否匿名';
COMMENT ON COLUMN reviews.merchant_reply  IS '商家回复内容';
COMMENT ON COLUMN reviews.reply_time      IS '回复时间';
COMMENT ON COLUMN reviews.status          IS '状态：1显示 2隐藏';
COMMENT ON COLUMN reviews.created_at      IS '创建时间';
COMMENT ON COLUMN reviews.updated_at      IS '更新时间';

-- ================================================================
-- 收藏表 favorites
-- ================================================================
CREATE TABLE IF NOT EXISTS favorites (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID            NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    merchant_id     UUID            NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, merchant_id)
);

CREATE INDEX IF NOT EXISTS idx_favorites_user     ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_merchant ON favorites(merchant_id);

COMMENT ON TABLE  favorites             IS '收藏表：用户收藏商家';
COMMENT ON COLUMN favorites.id           IS '主键ID';
COMMENT ON COLUMN favorites.user_id      IS '用户ID（外键 users.id）';
COMMENT ON COLUMN favorites.merchant_id  IS '商家ID（外键 merchants.id）';
COMMENT ON COLUMN favorites.created_at   IS '收藏时间';

-- ================================================================
-- 商家入驻申请表 merchant_applies
-- ================================================================
CREATE TABLE IF NOT EXISTS merchant_applies (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID            NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    merchant_type       SMALLINT        NOT NULL,
    id_card             VARCHAR(20),
    id_card_front       VARCHAR(500),
    id_card_back        VARCHAR(500),
    business_license    VARCHAR(500),
    company_name        VARCHAR(200),
    contact_name        VARCHAR(50),
    contact_phone       VARCHAR(20),
    shop_name           VARCHAR(100),
    shop_address        VARCHAR(500),
    location_id         UUID            REFERENCES locations(id),
    longitude           DECIMAL(10,7),
    latitude            DECIMAL(10,7),
    introduction        TEXT,
    audit_status        SMALLINT        NOT NULL DEFAULT 2,
    audit_remark        VARCHAR(500),
    audit_time          TIMESTAMP,
    auditor_id          UUID            REFERENCES users(id),
    created_at          TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_merchant_applies_user   ON merchant_applies(user_id);
CREATE INDEX IF NOT EXISTS idx_merchant_applies_audit  ON merchant_applies(audit_status);
CREATE INDEX IF NOT EXISTS idx_merchant_applies_loc    ON merchant_applies(location_id);

COMMENT ON TABLE  merchant_applies                IS '商家入驻申请表';
COMMENT ON COLUMN merchant_applies.id              IS '主键ID';
COMMENT ON COLUMN merchant_applies.user_id         IS '申请人ID（外键 users.id）';
COMMENT ON COLUMN merchant_applies.merchant_type   IS '商家类型：1个人 2个体户 3公司';
COMMENT ON COLUMN merchant_applies.id_card         IS '身份证号';
COMMENT ON COLUMN merchant_applies.id_card_front   IS '身份证正面照片URL';
COMMENT ON COLUMN merchant_applies.id_card_back    IS '身份证反面照片URL';
COMMENT ON COLUMN merchant_applies.business_license IS '营业执照照片URL';
COMMENT ON COLUMN merchant_applies.company_name    IS '公司名称';
COMMENT ON COLUMN merchant_applies.contact_name    IS '联系人姓名';
COMMENT ON COLUMN merchant_applies.contact_phone   IS '联系电话';
COMMENT ON COLUMN merchant_applies.shop_name       IS '店铺名称';
COMMENT ON COLUMN merchant_applies.shop_address    IS '店铺地址';
COMMENT ON COLUMN merchant_applies.location_id     IS '所属地区（外键 locations.id）';
COMMENT ON COLUMN merchant_applies.longitude       IS '经度';
COMMENT ON COLUMN merchant_applies.latitude        IS '纬度';
COMMENT ON COLUMN merchant_applies.introduction    IS '个人/店铺介绍';
COMMENT ON COLUMN merchant_applies.audit_status    IS '审核状态：1通过 2待审核 3拒绝';
COMMENT ON COLUMN merchant_applies.audit_remark    IS '审核备注';
COMMENT ON COLUMN merchant_applies.audit_time      IS '审核时间';
COMMENT ON COLUMN merchant_applies.auditor_id      IS '审核人ID（外键 users.id）';
COMMENT ON COLUMN merchant_applies.created_at      IS '创建时间';
COMMENT ON COLUMN merchant_applies.updated_at      IS '更新时间';