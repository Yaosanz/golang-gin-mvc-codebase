START TRANSACTION;

-- Drop unique
ALTER TABLE role_permissions DROP CONSTRAINT IF EXISTS unique_role_permission;

-- Drop primary key
ALTER TABLE role_permissions DROP CONSTRAINT IF EXISTS role_permissions_pkey;

-- Add back role column
ALTER TABLE role_permissions ADD COLUMN role user_role;

-- Populate role column from roles table
UPDATE role_permissions SET role = roles.name::user_role FROM roles WHERE role_permissions.role_id = roles.id;

-- Add back composite primary key
ALTER TABLE role_permissions ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (role, permission_id);

-- Drop id column
ALTER TABLE role_permissions DROP COLUMN IF EXISTS id;

COMMIT;
