START TRANSACTION;

-- =====================================================
-- DROP INDEXES
-- =====================================================
DROP INDEX IF EXISTS uq_users_username;
DROP INDEX IF EXISTS uq_users_email;

DROP TYPE IF EXISTS user_role;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_deleted_at;

-- =====================================================
-- DROP TABLES (CASCADE removes indexes & FK automatically)
-- =====================================================
DROP TABLE IF EXISTS users CASCADE;

-- =====================================================
-- DROP ENUM TYPE
-- =====================================================
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'user_role'
    ) THEN
        DROP TYPE user_role;
    END IF;
END$$;

-- =====================================================
-- DROP EXTENSIONS (OPTIONAL BUT CLEAN)
-- =====================================================
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS citext;

COMMIT;
