ALTER TABLE role_permissions
DROP CONSTRAINT IF EXISTS fk_role_permissions_role;

ALTER TABLE role_permissions
DROP COLUMN IF EXISTS role_id;
