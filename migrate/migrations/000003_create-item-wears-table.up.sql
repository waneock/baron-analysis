CREATE TABLE IF NOT EXISTS item_wears (
    id BIGSERIAL PRIMARY KEY,
    item_id TEXT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    UNIQUE (item_id, name)
);