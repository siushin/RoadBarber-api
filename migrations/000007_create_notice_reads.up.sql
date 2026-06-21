-- ================================================================
-- 000007_create_notice_reads.up.sql
-- 通知已读记录表：标记某用户已读某公告。
-- we 项目登录后，/api/home/notices 按 user_id 过滤已读记录。
-- ================================================================

CREATE TABLE IF NOT EXISTS notice_reads (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    notice_id   UUID         NOT NULL REFERENCES notices(id)  ON DELETE CASCADE,
    read_at     TIMESTAMP    NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, notice_id)
);

CREATE INDEX IF NOT EXISTS idx_notice_reads_user  ON notice_reads(user_id);
CREATE INDEX IF NOT EXISTS idx_notice_reads_notice ON notice_reads(notice_id);

COMMENT ON TABLE  notice_reads         IS '通知已读记录表：标记某用户已读某公告';
COMMENT ON COLUMN notice_reads.id      IS '主键ID';
COMMENT ON COLUMN notice_reads.user_id IS '用户ID（外键 users.id）';
COMMENT ON COLUMN notice_reads.notice_id IS '公告ID（外键 notices.id）';
COMMENT ON COLUMN notice_reads.read_at IS '已读时间';