CREATE OR REPLACE FUNCTION fn_distinct_access_levels()
RETURNS TRIGGER AS $$
BEGIN
    -- Ensure distinct access levels
    NEW.access_levels := ARRAY(
        SELECT DISTINCT unnest(NEW.access_levels)
    );

    -- If access_levels is empty, delete the row
    IF array_length(NEW.access_levels, 1) IS NULL THEN
        DELETE FROM store_owners WHERE user_id = NEW.user_id AND store_id = NEW.store_id;
        RETURN NULL; -- Important to return NULL for DELETE operation in BEFORE trigger
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_distinct_access_levels
BEFORE INSERT OR UPDATE ON store_owners
FOR EACH ROW
EXECUTE FUNCTION fn_distinct_access_levels();
