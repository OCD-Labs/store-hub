-- Drop the trigger
DROP TRIGGER after_insert_sale ON sales;

-- Drop the function
DROP FUNCTION update_sale_overview();

-- Drop the sale_overview table
DROP TABLE sale_overview;