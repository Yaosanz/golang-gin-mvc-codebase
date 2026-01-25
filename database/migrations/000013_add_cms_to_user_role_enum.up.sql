START TRANSACTION;

-- =====================================================
-- ADD CMS TO USER ROLE ENUM
-- =====================================================
DO $$
BEGIN
    -- Add 'cms' to the user_role enum if it doesn't exist
    IF NOT EXISTS (
        SELECT 1 FROM pg_enum e
        JOIN pg_type t ON e.enumtypid = t.oid
        WHERE t.typname = 'user_role' AND e.enumlabel = 'cms'
    ) THEN
        ALTER TYPE user_role ADD VALUE 'cms';
    END IF;
END$$;

COMMIT;
