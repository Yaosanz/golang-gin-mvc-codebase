-- +migrate Up
-- Migrate existing user roles to the new user_roles junction table
INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
SELECT
    u.id as user_id,
    u.role_id as role_id,
    NOW() as created_at,
    NOW() as updated_at
FROM users u
WHERE u.role_id IS NOT NULL
  AND u.deleted_at IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM user_roles ur
      WHERE ur.user_id = u.id AND ur.role_id = u.role_id
  );

-- +migrate StatementEnd
