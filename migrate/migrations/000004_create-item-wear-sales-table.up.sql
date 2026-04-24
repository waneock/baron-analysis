CREATE TABLE IF NOT EXISTS item_wear_sales (
    id BIGSERIAL PRIMARY KEY,
    item_wear_id BIGINT NOT NULL REFERENCES item_wears(id) ON DELETE CASCADE,
    price NUMERIC(12, 2) NOT NULL,
    wear_value NUMERIC(10, 8),
    sold_on DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (item_wear_id, sold_on, price, wear_value)
);