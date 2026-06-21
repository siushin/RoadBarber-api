-- ================================================================
-- 000012_replace_v1_banner_copy.up.sql
-- 替换 000009 V1 历史 banner 的文案
--
-- 000009 写入的 V1 文案是考古（来源：we 17dd6a2）：
--   title='公路理发师' subtitle='让理发更便利' text='V1 · 2026-06-20'
-- 现按需求改为"返璞归真"主题：
--   title='Back to Basics' subtitle='返璞归真' text='只做预约理发一件事'
-- align 仍为 center（V1 原本就是居中），sort_order 50 不变（保持先播顺序）
-- 用 title+subtitle+text 三字段定位，避免误改其它行
-- ================================================================

UPDATE banners
   SET title    = 'Back to Basics',
       subtitle = '返璞归真',
       text     = '只做预约理发一件事'
 WHERE title    = '公路理发师'
   AND subtitle = '让理发更便利'
   AND text     = 'V1 · 2026-06-20'
   AND sort_order = 50;
