CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE permissions (
    id UUID PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE roles (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id),
    permission_id UUID NOT NULL REFERENCES permissions(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    name TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id),
    role_id UUID NOT NULL REFERENCES roles(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE products (
    id UUID PRIMARY KEY,
    sku TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE inventories (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity BIGINT NOT NULL,
    location TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    actor_id UUID,
    action TEXT NOT NULL,
    resource TEXT NOT NULL,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL
);

-- Seed permissions
INSERT INTO permissions (id, code, description, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'user.read', 'Read users', now(), now()),
    (uuid_generate_v4(), 'user.write', 'Create/Update users', now(), now()),
    (uuid_generate_v4(), 'role.read', 'Read roles', now(), now()),
    (uuid_generate_v4(), 'role.write', 'Create/Update roles', now(), now()),
    (uuid_generate_v4(), 'permission.read', 'Read permissions', now(), now()),
    (uuid_generate_v4(), 'permission.write', 'Create/Update permissions', now(), now()),
    (uuid_generate_v4(), 'product.read', 'Read products', now(), now()),
    (uuid_generate_v4(), 'product.write', 'Create/Update products', now(), now()),
    (uuid_generate_v4(), 'inventory.read', 'Read inventory', now(), now()),
    (uuid_generate_v4(), 'inventory.write', 'Adjust inventory', now(), now()),
    (uuid_generate_v4(), 'audit.read', 'Read audit logs', now(), now());

-- Seed roles
INSERT INTO roles (id, name, description, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'admin', 'Administrator', now(), now()),
    (uuid_generate_v4(), 'user', 'Default user', now(), now());

-- Map admin to all permissions; user to read-only basic ones
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r CROSS JOIN permissions p WHERE r.name = 'admin';

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN ('product.read', 'inventory.read')
WHERE r.name = 'user';
