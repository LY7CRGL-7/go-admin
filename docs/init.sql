-- ========================================
-- Kratos Admin Template - 数据库初始化
-- GORM AutoMigrate 会自动建表
-- 此文件仅提供可选的初始数据
-- ========================================

CREATE DATABASE IF NOT EXISTS admin_db;

-- PostgreSQL: 启用 UUID 扩展（可选）
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==================== 初始权限数据（可选） ====================
-- 取消注释以插入默认权限

-- INSERT INTO permissions (name, code, type, path, method, parent_id, sort, status, created_at, updated_at)
-- VALUES
-- ('用户管理',   'user:list',     'menu',   '/users',               'GET',    0, 1,  1, NOW(), NOW()),
-- ('创建用户',   'user:create',   'button', '/users',               'POST',   0, 2,  1, NOW(), NOW()),
-- ('更新用户',   'user:update',   'button', '/users/:id',           'PUT',    0, 3,  1, NOW(), NOW()),
-- ('删除用户',   'user:delete',   'button', '/users/:id',           'DELETE', 0, 4,  1, NOW(), NOW()),
-- ('角色管理',   'role:list',     'menu',   '/roles',               'GET',    0, 5,  1, NOW(), NOW()),
-- ('创建角色',   'role:create',   'button', '/roles',               'POST',   0, 6,  1, NOW(), NOW()),
-- ('更新角色',   'role:update',   'button', '/roles/:id',           'PUT',    0, 7,  1, NOW(), NOW()),
-- ('删除角色',   'role:delete',   'button', '/roles/:id',           'DELETE', 0, 8,  1, NOW(), NOW()),
-- ('分配权限',   'role:assign',   'button', '/roles/:id/perms',     'POST',   0, 9,  1, NOW(), NOW()),
-- ('权限列表',   'perm:list',     'menu',   '/permissions',         'GET',    0, 10, 1, NOW(), NOW()),
-- ('租户管理',   'tenant:list',   'menu',   '/tenants',             'GET',    0, 11, 1, NOW(), NOW()),
-- ('审计日志',   'audit:list',    'menu',   '/audit-logs',          'GET',    0, 12, 1, NOW(), NOW());
