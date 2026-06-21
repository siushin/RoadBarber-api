-- ================================================================
-- 000003_add_username_email_wechat.down.sql
-- 回滚：删除 username / email / 微信小程序相关字段与索引
-- ================================================================

DROP INDEX IF EXISTS idx_users_unionid;
DROP INDEX IF EXISTS uq_users_openid;
DROP INDEX IF EXISTS uq_users_email;
DROP INDEX IF EXISTS uq_users_username;

ALTER TABLE users
    DROP COLUMN IF EXISTS wx_avatar,
    DROP COLUMN IF EXISTS wx_nickname,
    DROP COLUMN IF EXISTS unionid,
    DROP COLUMN IF EXISTS openid,
    DROP COLUMN IF EXISTS email,
    DROP COLUMN IF EXISTS username;