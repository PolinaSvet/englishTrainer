/*
	script to create database structure
*/

--DROP TABLE IF EXISTS catalogs CASCADE;

-- delete all tables 'catalog_%'
DO $$ 
DECLARE
    r RECORD;
BEGIN
    FOR r IN 
        SELECT tablename 
        FROM pg_catalog.pg_tables 
        WHERE tablename LIKE 'catalog_%' 
        AND schemaname = 'public'  -- Или укажите схему, в которой вы ищете таблицы
    LOOP
        EXECUTE format('DROP TABLE IF EXISTS %I CASCADE', r.tablename);
    END LOOP;
END $$;


CREATE TABLE catalogs (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
	comment TEXT NOT NULL,
	dt_add BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
	enable BOOL NOT NULL DEFAULT TRUE
);
COMMENT ON TABLE catalogs IS 'Table for list of tables';





DROP FUNCTION IF EXISTS f_catalogs_insert, f_catalogs_update, f_catalogs_delete, f_catalogs_view;
--=======================
--1::: table: catalogs
--=======================
--insert
CREATE FUNCTION f_catalogs_insert(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	INSERT INTO catalogs (
		name,
		comment,
		dt_add,
		enable
	)    
	VALUES (
		(json_data ->> 'name')::TEXT,
		(json_data ->> 'comment')::TEXT
		(EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
		(json_data ->> 'enable')::BOOL
	) 
	RETURNING id INTO new_id;
	
	IF new_id IS NULL THEN
		RAISE EXCEPTION 'Parameter value cannot be null. ';
	END IF;

	SELECT json_build_object('id',new_id,'err','') INTO json_result;
  	RETURN json_result;

EXCEPTION
    WHEN others THEN
		GET STACKED DIAGNOSTICS err_context = PG_EXCEPTION_CONTEXT;
    	GET STACKED DIAGNOSTICS err_mess = MESSAGE_TEXT;

        SELECT json_build_object('id',null,'err',err_mess||err_context) INTO json_result;
  		RETURN json_result;  
END;
$$ LANGUAGE plpgsql;

--update
CREATE FUNCTION f_catalogs_update(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	UPDATE catalogs SET
		name = (json_data ->> 'name')::TEXT,
		comment = (json_data ->> 'comment')::TEXT,
		dt_add = (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
    	enable = (json_data ->> 'enable')::BOOL
	WHERE id = (json_data ->> 'id')::BIGINT 
	RETURNING id INTO new_id;
	
	IF new_id IS NULL THEN
		RAISE EXCEPTION 'Parameter value cannot be null. ';
	END IF;

	SELECT json_build_object('id',new_id,'err','') INTO json_result;
  	RETURN json_result;

EXCEPTION
    WHEN others THEN
		GET STACKED DIAGNOSTICS err_context = PG_EXCEPTION_CONTEXT;
    	GET STACKED DIAGNOSTICS err_mess = MESSAGE_TEXT;

        SELECT json_build_object('id',null,'err',err_mess||err_context) INTO json_result;
  		RETURN json_result;  
END;
$$ LANGUAGE plpgsql;


--delete
CREATE FUNCTION f_catalogs_delete(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	DELETE FROM catalogs WHERE id = (json_data ->> 'id')::BIGINT; 

	SELECT json_build_object('id',(json_data ->> 'id')::BIGINT,'err','') INTO json_result;
  	RETURN json_result;

EXCEPTION
    WHEN others THEN
		GET STACKED DIAGNOSTICS err_context = PG_EXCEPTION_CONTEXT;
    	GET STACKED DIAGNOSTICS err_mess = MESSAGE_TEXT;

        SELECT json_build_object('id',null,'err',err_mess||err_context) INTO json_result;
  		RETURN json_result;  
END;
$$ LANGUAGE plpgsql;

--select view
CREATE FUNCTION f_catalogs_view(
		json_data jsonb
) 
RETURNS TABLE (
	id BIGINT,
    name TEXT,
	comment TEXT,
	dt_add BIGINT,
	enable BOOL
) AS $$
DECLARE
  	par_id BIGINT = 0;
BEGIN

	IF (json_data ->> 'id') IS NOT NULL THEN
		par_id = (json_data ->> 'id')::BIGINT;
	END IF;

	RETURN QUERY
		SELECT catalogs.id,
			   catalogs.name,  
			    catalogs.comment,
			    catalogs.dt_add,
				catalogs.enable
		FROM catalogs
		WHERE
			(par_id = 0 OR catalogs.id = par_id)
		ORDER BY catalogs.name;
	
END;
$$ LANGUAGE plpgsql;

/*

SELECT * FROM f_catalogs_view('{
    "id": 0
     
}'::jsonb);

*/
