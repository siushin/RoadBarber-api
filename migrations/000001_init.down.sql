-- ================================================================
-- 000001_init.down.sql
-- 回滚：删除 000001_init 创建的所有表
-- 注意：locations 被多个表外键依赖，需先解约束或先删依赖表
-- ================================================================

DROP TABLE IF EXISTS favorites;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS merchant_services;
DROP TABLE IF EXISTS services;

-- merchants / shops / merchant_applies 都引用 locations.id
DROP TABLE IF EXISTS merchants;
DROP TABLE IF EXISTS shops;
DROP TABLE IF EXISTS merchant_applies;

-- merchant_profiles 引用 users.id
DROP TABLE IF EXISTS merchant_profiles;

-- 最后删 users 和 locations（基础表）
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS locations;

-- 不删除扩展（uuid-ossp / pgcrypto）以便其他数据库复用