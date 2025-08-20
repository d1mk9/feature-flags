-- +goose Up
CREATE TABLE IF NOT EXISTS features (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    key         TEXT        NOT NULL UNIQUE,
    description TEXT,
    enabled     BOOLEAN     NOT NULL DEFAULT false,
    value       JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS features;