-- +goose Up
CREATE INDEX IF NOT EXISTS idx_features_updated_at
  ON features (updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_features_enabled
  ON features (enabled);

CREATE INDEX IF NOT EXISTS idx_features_enabled_updated_at
  ON features (enabled, updated_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_features_enabled_updated_at;
DROP INDEX IF EXISTS idx_features_enabled;
DROP INDEX IF EXISTS idx_features_updated_at;