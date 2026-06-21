-- ================================================================
-- 000012_replace_v1_banner_copy.down.sql
-- 回滚：恢复 000009 的 V1 原文案
-- ================================================================

UPDATE banners
   SET title    = '公路理发师',
       subtitle = '让理发更便利',
       text     = 'V1 · 2026-06-20'
 WHERE title    = 'Back to Basics'
   AND subtitle = '返璞归真'
   AND text     = '只做预约理发一件事'
   AND sort_order = 50;
