/*
	script to create database structure
*/

DROP TABLE IF EXISTS records CASCADE;


CREATE TABLE records (
    id BIGSERIAL PRIMARY KEY,
	users_id BIGINT REFERENCES users(id) DEFAULT 0,
	glossary_id BIGINT REFERENCES glossary(id) DEFAULT 0,
	catalogs_id BIGINT NOT NULL DEFAULT 0,
    attempt INT NOT NULL DEFAULT 0,
	guess BOOL NOT NULL DEFAULT FALSE,
	dt_add BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
	UNIQUE (users_id, glossary_id, catalogs_id)
);
COMMENT ON TABLE records IS 'Table for list of guess words';


DROP FUNCTION IF EXISTS f_records_insert, f_records_update, f_records_delete, f_records_view;
--=======================
--1::: table: records
--=======================
--insert
CREATE FUNCTION f_records_insert(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	INSERT INTO records (
		users_id,
		glossary_id,
		catalogs_id,
		attempt,
		guess,
		dt_add
	)    
	VALUES (
		(json_data ->> 'users_id')::BIGINT,
		(json_data ->> 'glossary_id')::BIGINT,
		(json_data ->> 'catalogs_id')::BIGINT,
		(json_data ->> 'attempt')::INT,
		(json_data ->> 'guess')::BOOL,
		(EXTRACT(EPOCH FROM now()) * 1000)::BIGINT
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
CREATE FUNCTION f_records_update(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
	v_id BIGINT;
	v_users_id BIGINT;
	v_glossary_id BIGINT;
	v_catalogs_id BIGINT;
	v_exists BOOLEAN;
BEGIN
	 -- Получаем id из JSON
    v_id := COALESCE((json_data ->> 'id')::BIGINT, 0);
	v_users_id := COALESCE((json_data ->> 'users_id')::BIGINT, 0);
	v_glossary_id := COALESCE((json_data ->> 'glossary_id')::BIGINT, 0);
	v_catalogs_id := COALESCE((json_data ->> 'catalogs_id')::BIGINT, 0);

    -- Проверяем, существует ли запись с таким id
    -- SELECT EXISTS (SELECT 1 FROM records WHERE id = v_id) INTO v_exists;
	SELECT EXISTS (SELECT 1 FROM records WHERE users_id = v_users_id AND  glossary_id = v_glossary_id AND  catalogs_id = v_catalogs_id) INTO v_exists;

    IF v_exists THEN
        -- Если запись существует, обновляем ее
        UPDATE records SET
			users_id = (json_data ->> 'users_id')::BIGINT,
			glossary_id = (json_data ->> 'glossary_id')::BIGINT,
			catalogs_id = (json_data ->> 'catalogs_id')::BIGINT,
			attempt = (json_data ->> 'attempt')::INT,
			guess = (json_data ->> 'guess')::BOOL,
			dt_add = (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT

		WHERE users_id = v_users_id AND  glossary_id = v_glossary_id AND  catalogs_id = v_catalogs_id --id = v_id 
		RETURNING id INTO new_id; 
    ELSE
        -- Если запись не существует, вставляем новую
        INSERT INTO records (
			users_id,
			glossary_id,
			catalogs_id,
			attempt,
			guess,
			dt_add
		)    
		VALUES (
			(json_data ->> 'users_id')::BIGINT,
			(json_data ->> 'glossary_id')::BIGINT,
			(json_data ->> 'catalogs_id')::BIGINT,
			(json_data ->> 'attempt')::INT,
			(json_data ->> 'guess')::BOOL,
			(EXTRACT(EPOCH FROM now()) * 1000)::BIGINT
		) 
		RETURNING id INTO new_id;
    END IF;

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
CREATE FUNCTION f_records_delete(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	DELETE FROM records WHERE id = (json_data ->> 'id')::BIGINT; 

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
CREATE FUNCTION f_records_view(
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
		SELECT records.id,
			   records.users_id,  
			   records.glossary_id,
			   records.catalogs_id,
			   records.attempt,
			   records.guess,
			   records.dt_add
		FROM records
		WHERE
			(par_id = 0 OR records.id = par_id)
		ORDER BY records.name;

END;
$$ LANGUAGE plpgsql;

/*

SELECT * FROM f_records_view('{
    "id": 0
     
}'::jsonb);

*/
