-- +migrate Down
-- Remove migrated user roles (this is a data migration, so we don't actually remove data in down)
-- The data will remain in the user_roles table for safety

-- +migrate StatementEnd
