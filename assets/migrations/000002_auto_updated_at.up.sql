CREATE OR REPLACE FUNCTION auto_timestamps()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = now();
        RETURN NEW;
    END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION create_updated_at_trigger(table_name text) RETURNS void AS $$
    BEGIN
        EXECUTE 'CREATE TRIGGER ' || table_name || '_updated_at BEFORE UPDATE ON ' || table_name || ' FOR EACH ROW EXECUTE PROCEDURE auto_timestamps()';
    END;
$$ LANGUAGE plpgsql;

