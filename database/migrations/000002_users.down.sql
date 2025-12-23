START TRANSACTION;

-- Drop unique constraints first (best practice, clearer than relying on CASCADE)
DROP INDEX IF EXISTS uq_users_username;
DROP INDEX IF EXISTS uq_users_email;

-- Drop indexes explicitly (best practice, clearer than relying on CASCADE)
DROP INDEX IF EXISTS idx_users_name;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_deleted_at;

-- Drop the table safely
DROP TABLE IF EXISTS users;

COMMIT;