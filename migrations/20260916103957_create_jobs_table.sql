-- +goose Up
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    empresa VARCHAR(100) NOT NULL,
    cargo VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()

);

-- +goose Down
DROP TABLE IF EXISTS jobs;
