START TRANSACTION;

-- Add id column
ALTER TABLE role_permissions ADD COLUMN id UUID DEFAULT gen_random_uuid();

-- Drop existing primary key
ALTER TABLE role_permissions DROP CONSTRAINT role_permissions_pkey;

-- Add new primary key on id
ALTER TABLE role_permissions ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (id);

-- Drop role column
ALTER TABLE role_permissions DROP COLUMN role;

-- Add unique constraint on role_id and permission_id
ALTER TABLE role_permissions ADD CONSTRAINT unique_role_permission UNIQUE (role_id, permission_id);

COMMIT;
