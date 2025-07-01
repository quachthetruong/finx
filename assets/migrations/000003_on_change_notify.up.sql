CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- Ensure the uuid-ossp extension is available

CREATE OR REPLACE FUNCTION notify_trigger() RETURNS TRIGGER AS
$trigger$
DECLARE
    new_row   JSONB;
    old_row   JSONB;
    payload   TEXT;
    notify_id UUID;
BEGIN
    notify_id := uuid_generate_v4();
    new_row := row_to_json(NEW);
    old_row := row_to_json(OLD);

    payload := json_build_object(
        'id', notify_id,
        'timestamp', CURRENT_TIMESTAMP,
        'action', LOWER(TG_OP),
        'db_schema', TG_TABLE_SCHEMA,
        'table', TG_TABLE_NAME,
        'record', new_row,
        'old', old_row
        )::TEXT;

    -- Notify the channel
    PERFORM pg_notify('db_event', payload);

    RETURN NULL;
END;
$trigger$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION create_notify_trigger(table_name text) RETURNS void AS
$$
BEGIN
    EXECUTE 'CREATE TRIGGER ' || table_name || '_notify AFTER INSERT OR UPDATE OR DELETE ON ' || table_name ||
            ' FOR EACH ROW EXECUTE PROCEDURE notify_trigger()';
END;
$$ LANGUAGE plpgsql;
