-- Drop the trigger
DROP TRIGGER IF EXISTS trigger_update_sales_overview ON sales;

-- Drop the function
DROP FUNCTION IF EXISTS update_sales_overview();
DROP FUNCTION IF EXISTS reduce_sales_overview(bigint, bigint, bigint);


-- Drop the sales_overview table
DROP TABLE IF EXISTS sales_overview;
