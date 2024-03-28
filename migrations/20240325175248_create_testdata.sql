-- +goose Up
-- +goose StatementBegin
INSERT INTO stores (id, name, availability)
VALUES ('1', 'Быково', true), ('2', 'Подольск', true), ('3', 'Химки', true), ('4', 'Балашиха', false);

INSERT INTO items (id, name, size)
VALUES ('1', 'Кроссовки', '44'), ('2', 'Туфли', '38'), ('3', 'Ботинки', '42'), ('4', 'Тапки', '43'),
       ('5', 'Футболка', 'L'), ('6', 'Худи', 'XL'), ('7', 'Свитер', 'L'), ('8', 'Блузка', 'M'),
       ('9', 'Куртка', '48'), ('10', 'Пальто', '52'), ('11', 'Пиджак', '50'), ('12', 'Пуховик', '46');

INSERT INTO store_item (store_id, item_id, amount)
VALUES ('1', '1', 1000), ('1', '2', 1000), ('1', '3', 1000), ('1', '4', 1000), ('1', '5', 1000), 
       ('1', '6', 1000), ('1', '7', 1000), ('2', '6', 1000), ('2', '7', 1000), ('2', '8', 1000), 
       ('2', '9', 1000), ('2', '10', 1000), ('2', '11', 1000), ('2', '12', 1000), ('3', '1', 1000), 
       ('3', '3', 1000), ('3', '5', 1000), ('3', '7', 1000), ('3', '9', 1000), ('3', '11', 1000),
       ('4', '1', 1000), ('4', '2', 1000);

INSERT INTO reservations (uuid, client_id, store_id, item_id, amount, status)
VALUES ('test-uuid-1', '1', '1', '1', 10, 'created'), ('test-uuid-1', '1', '1', '2', 100, 'created'), 
       ('test-uuid-2', '1', '1', '1', 100, 'created'), ('test-uuid-2', '1', '1', '2', 100, 'created'), 
       ('test-uuid-3', '1', '1', '1', 111, 'applied'), ('test-uuid-3', '1', '1', '2', 111, 'applied');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM stores;
DELETE FROM items;
DELETE FROM store_item;
DELETE FROM reservations;
-- +goose StatementEnd
