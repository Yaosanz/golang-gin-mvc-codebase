START TRANSACTION;

CREATE TABLE role_permissions (
    role user_role NOT NULL,
    permission_id UUID NOT NULL,

    PRIMARY KEY (role, permission_id),

    CONSTRAINT fk_role_permissions_permission
        FOREIGN KEY (permission_id)
        REFERENCES permissions(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_role_permissions_role ON role_permissions(role);

COMMIT;
