-- +goose Up
CREATE TABLE products (
    seller_id BIGINT      NOT NULL,
    offer_id  BIGINT      NOT NULL, 
    name      VARCHAR(100) NOT NULL,
    price     BIGINT      NOT NULL,
    quantity  BIGINT      NOT NULL,
    PRIMARY KEY(seller_id, offer_id),
    CONSTRAINT no_duplicates UNIQUE(seller_id, offer_id)
);

-- +goose Down
DROP TABLE products;