ALTER TABLE role_permissions
ADD COLUMN IF NOT EXISTS role_id UUID;

ALTER TABLE role_permissions
ADD CONSTRAINT fk_role_permissions_role
FOREIGN KEY (role_id)
REFERENCES roles(id);
