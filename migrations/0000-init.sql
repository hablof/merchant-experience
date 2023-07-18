-- +goose Up
CREATE TABLE products (
    seller_id BIGINT      NOT NULL,
    offer_id  BIGINT      NOT NULL, 
    name      VARCHAR(32) NOT NULL,
    price     BIGINT      NOT NULL,
    quantity  BIGINT      NOT NULL
);

-- +goose Down
DROP TABLE products;