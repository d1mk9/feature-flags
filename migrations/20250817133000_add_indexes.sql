-- +goose Up
-- Индексы "для жизни": выборка включённых и последние обновлённые
CREATE INDEX IF NOT EXISTS idx_features_enabled ON features (enabled);
CREATE INDEX IF NOT EXISTS idx_features_updated_at ON features (updated_at);

-- +goose Down
DROP INDEX IF EXISTS idx_features_updated_at;
DROP INDEX IF EXISTS idx_features_enabled;