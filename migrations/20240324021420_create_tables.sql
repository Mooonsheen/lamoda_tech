-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stores (
  id           VARCHAR(10) PRIMARY KEY,
  name         VARCHAR(255) NOT NULL,
  availability BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
  id     VARCHAR(10) PRIMARY KEY,
  name   VARCHAR(255) NOT NULL,
  size   VARCHAR(5)
);

CREATE TABLE IF NOT EXISTS store_item (
  store_id     VARCHAR(10) NOT NULL,
  item_id      VARCHAR(10) NOT NULL,
  amount       INTEGER NOT NULL CHECK (amount >= 0),
  CONSTRAINT fk_store_item_1
    FOREIGN KEY(store_id) 
      REFERENCES stores(id),
  CONSTRAINT fk_store_item_2
    FOREIGN KEY(item_id) 
      REFERENCES items(id)
);

CREATE TABLE IF NOT EXISTS reservations (
  id        SERIAL PRIMARY KEY,
  uuid      VARCHAR(60) NOT NULL,
  client_id VARCHAR(10) NOT NULL,
  store_id  VARCHAR(10) NOT NULL,
  item_id   VARCHAR(10) NOT NULL,
  amount    INTEGER NOT NULL,
  status    VARCHAR(10),
  CONSTRAINT fk_reservations_1
    FOREIGN KEY(store_id) 
      REFERENCES stores(id),
  CONSTRAINT fk_reservations_2
    FOREIGN KEY(item_id) 
      REFERENCES items(id)
);

CREATE VIEW item_total_amount AS 
  SELECT item_id, SUM(amount) as total_amount
  FROM store_item
  GROUP BY item_id
  ORDER BY total_amount DESC;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS item_total_amount;
DROP TABLE IF EXISTS store_item;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS stores;
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
