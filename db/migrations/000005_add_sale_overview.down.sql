-- Drop the trigger
DROP TRIGGER IF EXISTS trigger_update_sale_overview ON sales;

-- Drop the function
DROP FUNCTION IF EXISTS update_sale_overview();
DROP FUNCTION IF EXISTS reduce_sale(bigint, bigint);


-- Drop the sale_overview table
DROP TABLE IF EXISTS sale_overview;
