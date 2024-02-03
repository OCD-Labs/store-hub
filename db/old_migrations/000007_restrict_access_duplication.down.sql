-- Drop the trigger
DROP TRIGGER trigger_distinct_access_levels ON store_owners;

-- Drop the function
DROP FUNCTION fn_distinct_access_levels();

-- Revert the sales_percentage column
ALTER TABLE sales_overview ALTER COLUMN sales_percentage TYPE NUMERIC(6, 4);
