-- +goose Up
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS status VARCHAR(40) NOT NULL;

-- +goose Down
ALTER TABLE jobs DROP COLUMN status;