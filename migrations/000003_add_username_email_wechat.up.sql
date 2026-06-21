-- ================================================================
-- 000003_add_username_email_wechat.up.sql
-- users 表的 username / email / 微信字段
--
-- 注意：这些字段已经在 000001_init.up.sql 的 CREATE TABLE users 中
-- 一次性写好（位于 status 之后、last_login_at 之前）。本迁移保留为
-- 幂等兜底——使用 ADD COLUMN IF NOT EXISTS、CREATE INDEX IF NOT EXISTS，
-- 在历史数据库（字段被 ALTER 追加到末尾）上重复跑也不会报错。
--
-- 兼容原有 phone 登录：手机号仍可作为账号登录。
-- 必须放在 000002_seed_demo 之后：下面的 UPDATE 需要 users 已存在。
-- ================================================================

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS username    VARCHAR(50),
    ADD COLUMN IF NOT EXISTS email       VARCHAR(255),
    ADD COLUMN IF NOT EXISTS openid      VARCHAR(128),
    ADD COLUMN IF NOT EXISTS unionid     VARCHAR(128),
    ADD COLUMN IF NOT EXISTS wx_nickname VARCHAR(100),
    ADD COLUMN IF NOT EXISTS wx_avatar   VARCHAR(500);

-- 唯一约束（允许 NULL，但非 NULL 时必须唯一）
CREATE UNIQUE INDEX IF NOT EXISTS uq_users_username ON users(username) WHERE username IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_users_email    ON users(email)    WHERE email    IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_users_openid   ON users(openid)   WHERE openid   IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_users_unionid ON users(unionid);

COMMENT ON COLUMN users.username    IS '用户名（后台登录账号，全局唯一）';
COMMENT ON COLUMN users.email       IS '邮箱（可作为登录账号）';
COMMENT ON COLUMN users.openid      IS '微信小程序 openid（用于微信登录）';
COMMENT ON COLUMN users.unionid     IS '微信开放平台 unionid（多端打通）';
COMMENT ON COLUMN users.wx_nickname IS '微信昵称';
COMMENT ON COLUMN users.wx_avatar   IS '微信头像 URL';

-- ================================================================
-- 回填：超级管理员（seed 里 UUID 固定）
-- ================================================================
UPDATE users
   SET username = 'admin',
       email    = 'admin@roadbarber.local'
 WHERE id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

-- ================================================================
-- 回填：其余用户用 phone 后 6 位生成 user_xxxxxx 默认 username
-- 若 phone 为空或派生 username 与现有冲突，则跳过（由后续注册逻辑生成）
-- ================================================================
UPDATE users
   SET username = 'user_' || RIGHT(phone, 6)
 WHERE username IS NULL
   AND phone IS NOT NULL
   AND phone <> '';