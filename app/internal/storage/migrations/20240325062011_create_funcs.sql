-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_store_item_amount_after_create_reservation()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE store_item
    SET amount = amount - NEW.amount
    WHERE store_id = NEW.store_id AND item_id = NEW.item_id;
    RETURN NEW;
END;
$$
 LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_store_item_after_create_reservation_trigger
AFTER INSERT ON reservations
FOR EACH ROW
EXECUTE FUNCTION update_store_item_amount_after_create_reservation();


CREATE OR REPLACE FUNCTION update_store_item_amount_after_delete_reservation()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE store_item
    SET amount = amount + NEW.amount
    WHERE store_id = NEW.store_id AND item_id = NEW.item_id AND NEW.status = 'deleted';
    RETURN NEW;
END;
$$
 LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_store_item_after_delete_reservation_trigger
AFTER UPDATE ON reservations
FOR EACH ROW
EXECUTE FUNCTION update_store_item_amount_after_delete_reservation();

CREATE OR REPLACE FUNCTION create_reservation_one_store(in_uuid VARCHAR, in_client_id VARCHAR, in_store_id VARCHAR, in_item_id VARCHAR, in_amount INTEGER)
RETURNS INTEGER AS 
$$
DECLARE 
	Amount_in_store INTEGER;
  Exec_code INTEGER;
BEGIN  
	SELECT amount INTO Amount_in_store
  FROM store_item
  WHERE store_id = in_store_id AND item_id = in_item_id;
  IF Amount_in_store > in_amount THEN 
  	INSERT INTO reservations ("uuid", "client_id", "store_id", "item_id", "amount", "status") 
    VALUES (in_uuid, in_client_id, in_store_id, in_item_id, in_amount, 'created');
  	Exec_code := 1;
  ELSE
  	Exec_code := 0;
  END IF;
  RETURN Exec_code;
END;
$$
 LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_store_item_after_create_reservation_trigger on reservations;
DROP FUNCTION update_store_item_amount_after_create_reservation();
DROP TRIGGER update_store_item_after_delete_reservation_trigger on reservations;
DROP FUNCTION update_store_item_amount_after_delete_reservation();
DROP FUNCTION create_reservation_one_store(in_uuid VARCHAR, in_client_id VARCHAR, in_store_id VARCHAR, in_item_id VARCHAR, in_amount INTEGER);
-- +goose StatementEnd

