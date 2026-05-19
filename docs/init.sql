-- 创建数据库
CREATE DATABASE admin_db;

-- 切换到数据库
\c admin_db;

-- 启用 UUID 扩展（可选）
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 说明：
-- 数据表会通过 GORM AutoMigrate 自动创建
-- 包括：admins, roles, permissions, admin_roles, role_permissions, audit_logs, login_attempts

-- 如果需要手动初始化基础数据，可以取消以下注释：

-- 插入超级管理员角色
-- INSERT INTO roles (name, code, description, status, sort, created_at, updated_at)
-- VALUES ('超级管理员', 'super_admin', '拥有系统所有权限', 1, 0, NOW(), NOW());

-- 插入基础权限数据
-- INSERT INTO permissions (name, code, type, path, method, parent_id, sort, status, created_at, updated_at)
-- VALUES 
-- ('管理员列表', 'admin:list', 'api', '/api/v1/admins', 'GET', 0, 1, 1, NOW(), NOW()),
-- ('创建管理员', 'admin:create', 'api', '/api/v1/admins', 'POST', 0, 2, 1, NOW(), NOW()),
-- ('更新管理员', 'admin:update', 'api', '/api/v1/admins/:id', 'PUT', 0, 3, 1, NOW(), NOW()),
-- ('删除管理员', 'admin:delete', 'api', '/api/v1/admins/:id', 'DELETE', 0, 4, 1, NOW(), NOW()),
-- ('角色列表', 'role:list', 'api', '/api/v1/roles', 'GET', 0, 5, 1, NOW(), NOW()),
-- ('创建角色', 'role:create', 'api', '/api/v1/roles', 'POST', 0, 6, 1, NOW(), NOW()),
-- ('更新角色', 'role:update', 'api', '/api/v1/roles/:id', 'PUT', 0, 7, 1, NOW(), NOW()),
-- ('删除角色', 'role:delete', 'api', '/api/v1/roles/:id', 'DELETE', 0, 8, 1, NOW(), NOW()),
-- ('分配权限', 'role:assign-permissions', 'api', '/api/v1/roles/:id/permissions', 'POST', 0, 9, 1, NOW(), NOW()),
-- ('权限列表', 'permission:list', 'api', '/api/v1/permissions', 'GET', 0, 10, 1, NOW(), NOW()),
-- ('审计日志', 'audit-log:list', 'api', '/api/v1/audit-logs', 'GET', 0, 11, 1, NOW(), NOW());
