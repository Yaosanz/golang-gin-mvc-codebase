START TRANSACTION;

CREATE TABLE shortenlink (
    id UUID PRIMARY KEY,

    original_url TEXT NOT NULL,
    short_code VARCHAR(10) NOT NULL UNIQUE,

    user_id UUID NOT NULL,
    created_by UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

ALTER TABLE shortenlink
ADD CONSTRAINT fk_shortenlink_user
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE;

ALTER TABLE shortenlink
ADD CONSTRAINT fk_shortenlink_created_by
FOREIGN KEY (created_by)
REFERENCES users(id)
ON DELETE RESTRICT;

CREATE INDEX idx_shortenlink_code ON shortenlink(short_code);
CREATE INDEX idx_shortenlink_deleted_at ON shortenlink(deleted_at);

COMMIT;
