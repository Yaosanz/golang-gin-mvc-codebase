START TRANSACTION;

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY, -- auto-incrementing primary key
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL, -- timestamp for soft deletion (optional)

    created_by BIGINT NULL, -- user ID who created the record
    updated_by BIGINT NULL, -- user ID who last updated the record
    deleted_by BIGINT NULL, -- user ID who deleted the record (if applicable)

    name CITEXT NOT NULL,
    username CITEXT NOT NULL,
    email CITEXT NOT NULL,
    phone VARCHAR(20) NULL,
    password VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- Useful indexes
CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Unique constraints
CREATE UNIQUE INDEX uq_users_username ON users(username);
CREATE UNIQUE INDEX uq_users_email ON users(email);

COMMIT;