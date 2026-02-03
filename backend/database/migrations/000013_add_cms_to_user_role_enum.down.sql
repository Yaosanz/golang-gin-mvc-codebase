START TRANSACTION;

-- =====================================================
-- REMOVE CMS FROM USER ROLE ENUM
-- =====================================================
-- Note: PostgreSQL doesn't support removing enum values directly
-- This is a placeholder for potential future rollback needs
-- In practice, enum values cannot be removed once added

COMMIT;
