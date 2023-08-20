-- Create the sales_overview table
CREATE TABLE "sales_overview" (
    "id" bigserial PRIMARY KEY,
    "number_of_sales" bigint NOT NULL DEFAULT 0,
    "sales_percentage" NUMERIC(6, 4) NOT NULL DEFAULT 0,
    "revenue" NUMERIC(10, 2) NOT NULL DEFAULT 0,
    "item_id" bigint NOT NULL,
    "store_id" bigint NOT NULL,
    FOREIGN KEY ("item_id") REFERENCES "items" ("id"),
    FOREIGN KEY ("store_id") REFERENCES "stores" ("id")
);

-- Create the update_sales_overview function
CREATE OR REPLACE FUNCTION update_sales_overview()
RETURNS TRIGGER AS $$
DECLARE
    orderQty int;
    itemPrice NUMERIC(10, 2);
    supplyQuantity NUMERIC;
BEGIN
    -- Fetch the order quantity and item details for the sale
    SELECT o.order_quantity, i.price, i.supply_quantity 
    INTO orderQty, itemPrice, supplyQuantity 
    FROM orders o
    JOIN items i ON o.item_id = i.id
    WHERE o.id = NEW.order_id;

    -- Check if the item and store combination already exists in sales_overview
    IF EXISTS (SELECT 1 FROM sales_overview WHERE item_id = NEW.item_id AND store_id = NEW.store_id) THEN
        -- Update the existing record
        UPDATE sales_overview
        SET 
            number_of_sales = number_of_sales + orderQty,
            sales_percentage = ((number_of_sales + orderQty) / supplyQuantity) * 100,
            revenue = (number_of_sales + orderQty) * itemPrice
        WHERE item_id = NEW.item_id AND store_id = NEW.store_id;
    ELSE
        -- Insert a new record
        INSERT INTO sales_overview (number_of_sales, sales_percentage, revenue, item_id, store_id)
        VALUES (
            orderQty,
            (orderQty / supplyQuantity) * 100,
            orderQty * itemPrice,
            NEW.item_id,
            NEW.store_id
        );
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger to call the update_sales_overview function after a sale is inserted
CREATE TRIGGER trigger_update_sales_overview
AFTER INSERT ON sales
FOR EACH ROW
EXECUTE FUNCTION update_sales_overview();

-- Create the reduce_sales_overview function
CREATE OR REPLACE FUNCTION reduce_sales_overview(item_id_arg bigint, store_id_arg bigint, order_id_arg bigint)
RETURNS void AS $$
DECLARE
    orderQty int;
    itemPrice NUMERIC(10, 2);
    supplyQuantity NUMERIC;
BEGIN
    -- Fetch the order quantity and item details for the sale
    SELECT o.order_quantity, i.price, i.supply_quantity 
    INTO orderQty, itemPrice, supplyQuantity 
    FROM orders o
    JOIN items i ON i.id = item_id_arg
    WHERE o.id = order_id_arg;

    -- Reduce the number of sales by the order quantity
    UPDATE sales_overview
    SET 
        number_of_sales = GREATEST(number_of_sales - orderQty, 0), -- Ensure it doesn't go below 0
        sales_percentage = (GREATEST(number_of_sales - orderQty, 0) / supplyQuantity) * 100,
        revenue = GREATEST(number_of_sales - orderQty, 0) * itemPrice
    WHERE item_id = item_id_arg AND store_id = store_id_arg;
END;
$$ LANGUAGE plpgsql;


