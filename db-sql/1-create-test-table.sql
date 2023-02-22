CREATE TABLE IF NOT EXISTS test
(
    test_id uuid primary key,
    created_at timestamptz NOT NULL,
    updated_at timestamptz default CURRENT_TIMESTAMP
);