START TRANSACTION;

-- Drop unique
ALTER TABLE role_permissions DROP CONSTRAINT IF EXISTS unique_role_permission;

-- Drop primary key
ALTER TABLE role_permissions DROP CONSTRAINT IF EXISTS role_permissions_pkey;

-- Add back role column
ALTER TABLE role_permissions ADD COLUMN role user_role;

-- Add back composite primary key
ALTER TABLE role_permissions ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (role, permission_id);

-- Drop id column
ALTER TABLE role_permissions DROP COLUMN IF EXISTS id;

COMMIT;
