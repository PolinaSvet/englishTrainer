/*
	script to create database structure
*/

DROP TABLE IF EXISTS catalogsglos CASCADE;


CREATE TABLE catalogsglos (
    id BIGSERIAL PRIMARY KEY,
	glossary_id BIGINT REFERENCES glossary(id) DEFAULT 0,
	catalogs_id BIGINT REFERENCES catalogs(id) DEFAULT 0,
	dt_add BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT
);
COMMENT ON TABLE catalogsglos IS 'Table for list of filter for words';


DROP FUNCTION IF EXISTS f_catalogsglos_insert, f_catalogsglos_update, f_catalogsglos_delete, f_catalogsglos_view;
--=======================
--1::: table: catalogsglos
--=======================
--insert
CREATE FUNCTION f_catalogsglos_insert(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	INSERT INTO catalogsglos (
		glossary_id,
		catalogs_id,
		dt_add
	)    
	VALUES (
		(json_data ->> 'glossary_id')::BIGINT,
		(json_data ->> 'catalogs_id')::BIGINT,
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
CREATE FUNCTION f_catalogsglos_update(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	UPDATE catalogsglos SET
		glossary_id = (json_data ->> 'glossary_id')::BIGINT,
		catalogs_id = (json_data ->> 'catalogs_id')::BIGINT,
		dt_add = (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT

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
CREATE FUNCTION f_catalogsglos_delete(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	DELETE FROM catalogsglos WHERE id = (json_data ->> 'id')::BIGINT; 

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
CREATE FUNCTION f_catalogsglos_view(
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
		SELECT catalogsglos.id,
			   catalogsglos.glossary_id,
			   catalogsglos.catalogs_id,
			   catalogsglos.dt_add
		FROM catalogsglos
		WHERE
			(par_id = 0 OR catalogsglos.id = par_id)
		ORDER BY catalogsglos.name;

END;
$$ LANGUAGE plpgsql;

