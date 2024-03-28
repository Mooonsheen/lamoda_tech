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
  IF Amount_in_store >= in_amount THEN 
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

 CREATE OR REPLACE FUNCTION delete_reservation_one_store(in_client_id VARCHAR, in_uuid VARCHAR)
RETURNS INTEGER AS
$$
DECLARE
	Exec_code INTEGER;
  rows_updated INTEGER;
BEGIN
	UPDATE reservations
  SET status = 'deleted'
  WHERE client_id = in_client_id AND uuid = in_uuid AND status = 'created';
  
  GET DIAGNOSTICS rows_updated = ROW_COUNT;
  
  IF rows_updated > 0 THEN
        RETURN 1;
    ELSE
        RETURN 0;
    END IF;
END;
$$
	LANGUAGE plpgsql;

  CREATE OR REPLACE FUNCTION apply_reservation_one_store(in_client_id VARCHAR, in_uuid VARCHAR)
RETURNS INTEGER AS
$$
DECLARE
	Exec_code INTEGER;
  rows_updated INTEGER;
BEGIN
	UPDATE reservations
  SET status = 'applied'
  WHERE client_id = in_client_id AND uuid = in_uuid AND status = 'created';
  
  GET DIAGNOSTICS rows_updated = ROW_COUNT;
  
  IF rows_updated > 0 THEN
        RETURN 1;
    ELSE
        RETURN 0;
    END IF;
END;
$$
	LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION check_store_availability(in_store_id VARCHAR)
RETURNS BOOLEAN AS
$$
DECLARE
	is_available BOOLEAN;
BEGIN
	SELECT availability INTO is_available
	FROM stores
	WHERE id = in_store_id;
  RETURN is_available;
END;
$$
	LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_store_remains(in_store_id VARCHAR, sort VARCHAR, pages INTEGER)
RETURNS TABLE(out_item_id VARCHAR, out_amount INTEGER) AS
$$
DECLARE
	item_count_in_one_page INTEGER;
  current_limit Integer;
BEGIN
	item_count_in_one_page = 5;
  current_limit := item_count_in_one_page * pages;
	IF sort = 'asc' THEN
		RETURN QUERY SELECT item_id, amount 
    							FROM store_item 
                  WHERE store_id = in_store_id
                  ORDER BY amount ASC
                  LIMIT current_limit;
  ELSE
  	RETURN QUERY SELECT item_id, amount 
    							FROM store_item 
                  WHERE store_id = in_store_id
                  ORDER BY amount DESC
                  LIMIT current_limit;
  END IF;
END;
$$
	LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION create_reservation_many_stores(in_uuid VARCHAR, in_client_id VARCHAR, in_store_ids TEXT[], in_item_id VARCHAR, in_amount INTEGER)
RETURNS INTEGER AS 
$$
DECLARE
  Exec_code INTEGER;
  Have_now INTEGER;
  Store_counter INTEGER;
  Current_store_id VARCHAR;
  Amount_in_current_store INTEGER;
 	Status_from_create_reservation_one_store INTEGER;
  loop_counter INTEGER := 1;
  Amount_to_get_remains INTEGER;
  
BEGIN
CREATE TEMP TABLE required_stores_item (store_id, amount)
AS
(SELECT store_id, amount
FROM store_item
WHERE store_id = ANY(in_store_ids) AND item_id = in_item_id AND amount != 0
ORDER BY store_item.amount ASC);

SELECT SUM(amount) INTO Have_now
FROM required_stores_item;
IF in_amount > Have_now THEN
  Exec_code := 0;
  DROP TABLE required_stores_item;
  RETURN Exec_code;
END IF;

Have_now = 0;
Store_counter := 1;

WHILE Have_now < in_amount LOOP
	
  SELECT SUM(amount) INTO Have_now
  FROM 
  (SELECT amount 
   FROM required_stores_item
  LIMIT Store_counter);
  
  IF Have_now < in_amount THEN
    Store_counter = Store_counter + 1;
    CONTINUE;
  ELSE
    IF Store_counter = 1 THEN
    	SELECT store_id INTO Current_store_id
      FROM required_stores_item
      LIMIT Store_counter;
      
    	SELECT create_reservation_one_store(in_uuid, in_client_id, Current_store_id, in_item_id, in_amount)
      INTO Status_from_create_reservation_one_store;
      
      IF Status_from_create_reservation_one_store = 1 THEN
      	Exec_code := 1;
        DROP TABLE required_stores_item;
      	RETURN Exec_code;
      ELSE
      	Exec_code := 2;
        DROP TABLE required_stores_item;
      	RETURN Exec_code;
      END IF;
    ELSE
    	WHILE loop_counter < Store_counter LOOP
        SELECT store_id INTO Current_store_id
      	FROM required_stores_item
        OFFSET loop_counter - 1
      	LIMIT 1;
        
        SELECT amount INTO Amount_in_current_store
      	FROM store_item
        WHERE store_id = Current_store_id AND item_id = in_item_id;
        
        SELECT create_reservation_one_store(in_uuid, in_client_id, Current_store_id, in_item_id, Amount_in_current_store)
      	INTO Status_from_create_reservation_one_store;
        
        IF Status_from_create_reservation_one_store = 1 THEN
      		Exec_code := 1;
      	ELSE
      		Exec_code := 3;
          DROP TABLE required_stores_item;
      		RETURN Exec_code;
      	END IF;
        loop_counter = loop_counter + 1;
    	END LOOP;
      
    	Amount_to_get_remains = Have_now - in_amount;
      SELECT store_id INTO Current_store_id
      FROM required_stores_item
      OFFSET Store_counter - 1
      LIMIT 1;
      
      SELECT amount INTO Amount_in_current_store
      FROM store_item
      WHERE store_id = Current_store_id AND item_id = in_item_id;
      
      SELECT create_reservation_one_store(in_uuid, in_client_id, Current_store_id, in_item_id, (Amount_in_current_store - Amount_to_get_remains))
      INTO Status_from_create_reservation_one_store;
        
      IF Status_from_create_reservation_one_store = 1 THEN
      		Exec_code := 1;
      	ELSE
      		Exec_code := 4;
          DROP TABLE required_stores_item;
      		RETURN Exec_code;
      	END IF;
      END IF;
      DROP TABLE required_stores_item;
    RETURN Exec_code;
  END IF;
  END LOOP;
END;
$$
  LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_all_available_stores()
RETURNS TABLE(store_id VARCHAR) AS
$$
BEGIN
	RETURN QUERY SELECT id
    							FROM stores
                  WHERE availability = true;
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
DROP FUNCTION delete_reservation_one_store(in_client_id VARCHAR, in_uuid VARCHAR);
DROP FUNCTION apply_reservation_one_store(in_client_id VARCHAR, in_uuid VARCHAR);
DROP FUNCTION check_store_availability(in_store_id VARCHAR);
DROP FUNCTION get_store_remains(in_store_id VARCHAR, sort VARCHAR, pages INTEGER);
DROP FUNCTION create_reservation_many_stores(in_uuid VARCHAR, in_client_id VARCHAR, in_store_ids TEXT[], in_item_id VARCHAR, in_amount INTEGER);
DROP FUNCTION get_all_available_stores();
-- +goose StatementEnd

