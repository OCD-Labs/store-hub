-- Create the sale_overview table
CREATE TABLE sale_overview (
    store_id bigint NOT NULL,
    item_id bigint NOT NULL,
    number_of_sales int NOT NULL DEFAULT 0,
    PRIMARY KEY (store_id, item_id),
    FOREIGN KEY (store_id) REFERENCES stores (id),
    FOREIGN KEY (item_id) REFERENCES items (id)
);

-- Create the update_sale_overview function
CREATE OR REPLACE FUNCTION update_sale_overview()
RETURNS TRIGGER AS $$
BEGIN
    -- Increment the number_of_sales for the corresponding store_id and item_id
    UPDATE sale_overview
    SET number_of_sales = number_of_sales + 1
    WHERE store_id = NEW.store_id AND item_id = NEW.item_id;

    -- If the row doesn't exist yet, insert a new row
    IF NOT FOUND THEN
        INSERT INTO sale_overview (store_id, item_id, number_of_sales)
        VALUES (NEW.store_id, NEW.item_id, 1);
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger to call the update_sale_overview function after a sale is inserted
CREATE TRIGGER after_insert_sale
AFTER INSERT ON sales
FOR EACH ROW
EXECUTE FUNCTION update_sale_overview();

-- Create the reduce_sale_count function
CREATE OR REPLACE FUNCTION reduce_sale_count(p_store_id bigint, p_item_id bigint)
RETURNS VOID AS $$
BEGIN
    -- Decrement the number_of_sales for the specified store_id and item_id
    UPDATE sale_overview
    SET number_of_sales = GREATEST(number_of_sales - 1, 0)
    WHERE store_id = p_store_id AND item_id = p_item_id;
    
    -- If the row doesn't exist yet, insert a new row with 0 count
    IF NOT FOUND THEN
        INSERT INTO sale_overview (store_id, item_id, number_of_sales)
        VALUES (p_store_id, p_item_id, 0);
    END IF;
END;
$$ LANGUAGE plpgsql;
