/*
	script to create procedures
*/
DROP FUNCTION IF EXISTS f_letter_insert, f_letter_update, f_letter_delete, f_letter_view;
DROP FUNCTION IF EXISTS f_mark_insert, f_mark_update, f_mark_delete, f_mark_view;
DROP FUNCTION IF EXISTS f_glossary_insert, f_glossary_insert_array, f_glossary_update, f_glossary_delete, f_glossary_view, f_glossary_view_random;


--=======================
--1::: table: letter
--=======================
--insert
CREATE FUNCTION f_letter_insert(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	INSERT INTO letter (name) VALUES ((json_data ->> 'name')::TEXT) RETURNING id INTO new_id;
	
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
CREATE FUNCTION f_letter_update(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	UPDATE letter SET name = (json_data ->> 'name')::TEXT WHERE id = (json_data ->> 'id')::BIGINT RETURNING id INTO new_id; 
	
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
CREATE FUNCTION f_letter_delete(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	DELETE FROM letter WHERE id = (json_data ->> 'id')::BIGINT; 

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
CREATE FUNCTION f_letter_view(
		json_data jsonb
) 
RETURNS TABLE (
	id BIGINT,
    name TEXT
) AS $$
DECLARE
  	par_id BIGINT = 0;

BEGIN

	IF (json_data ->> 'id') IS NOT NULL THEN
		par_id = (json_data ->> 'id')::BIGINT;
	END IF;

	RETURN QUERY
		SELECT letter.id,
			   letter.name   
		FROM letter
		WHERE
			(par_id = 0 OR letter.id = par_id)
		ORDER BY letter.id;
	
END;
$$ LANGUAGE plpgsql;


--=======================
--2::: table: mark
--=======================
--insert
CREATE FUNCTION f_mark_insert(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	INSERT INTO mark (name) VALUES ((json_data ->> 'name')::TEXT) RETURNING id INTO new_id;
	
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
CREATE FUNCTION f_mark_update(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	UPDATE mark SET name = (json_data ->> 'name')::TEXT WHERE id = (json_data ->> 'id')::BIGINT RETURNING id INTO new_id; 
	
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
CREATE FUNCTION f_mark_delete(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	DELETE FROM mark WHERE id = (json_data ->> 'id')::BIGINT; 

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
CREATE FUNCTION f_mark_view(
		json_data jsonb
) 
RETURNS TABLE (
	id BIGINT,
    name TEXT
) AS $$
DECLARE
  	par_id BIGINT = 0;

BEGIN

	IF (json_data ->> 'id') IS NOT NULL THEN
		par_id = (json_data ->> 'id')::BIGINT;
	END IF;

	RETURN QUERY
		SELECT mark.id,
			   mark.name   
		FROM mark
		WHERE
			(par_id = 0 OR mark.id = par_id)
		ORDER BY mark.id;
	
END;
$$ LANGUAGE plpgsql;

--=======================
--3::: table: glossary
--=======================
--insert
CREATE FUNCTION f_glossary_insert(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
	mark_id BIGINT;
	letter_id BIGINT;
BEGIN

	SELECT mark.id
	INTO   mark_id
	FROM   mark
	WHERE  ((json_data ->> 'mark_name')::TEXT = mark.name);

	SELECT letter.id
	INTO   letter_id
	FROM   letter
	WHERE  ((json_data ->> 'letter_name')::TEXT = letter.name);

	INSERT INTO glossary (
		mark_id, 
		letter_id, 
		word,
		transcription,
		translation,
		example,
		dt_add,
    	enable    
		) 
	VALUES (
		mark_id, 
		letter_id, 
		(json_data ->> 'word')::TEXT, 
		(json_data ->> 'transcription')::TEXT,
		(json_data ->> 'translation')::TEXT,
		(json_data ->> 'example')::JSONB,
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

--insert array
CREATE FUNCTION f_glossary_insert_array(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
	json_record jsonb;
BEGIN

	IF jsonb_typeof(json_data) != 'array' THEN
        RAISE EXCEPTION 'Input must be a JSON array. ';
    END IF;


	FOR json_record IN SELECT * FROM jsonb_array_elements(json_data)
    LOOP

		json_result := f_glossary_insert(json_record);

		IF (json_result ->> 'id')::BIGINT IS NULL THEN
			RAISE EXCEPTION 'Parameter value cannot be null. ';
		END IF;

	END LOOP;

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
CREATE FUNCTION f_glossary_update(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
	par_mark_id BIGINT;
	par_letter_id BIGINT;
BEGIN

	--RAISE NOTICE 'Текущий текст: %', json_data;

	SELECT mark.id
	INTO   par_mark_id
	FROM   mark
	WHERE  ((json_data ->> 'mark_name')::TEXT = mark.name);

	SELECT letter.id
	INTO   par_letter_id
	FROM   letter
	WHERE  ((json_data ->> 'letter_name')::TEXT = letter.name);

	UPDATE glossary SET 
		mark_id  = par_mark_id,
		letter_id = par_letter_id, 
		word = (json_data ->> 'word')::TEXT,
		transcription = (json_data ->> 'transcription')::TEXT,
		translation = (json_data ->> 'translation')::TEXT,
		example = (json_data ->> 'example')::JSONB,
		dt_add = (json_data ->> 'dt_add')::BIGINT,
    	enable = (json_data ->> 'enable')::BOOL
	WHERE 
		id = (json_data ->> 'id')::BIGINT
	RETURNING id INTO new_id;
	
	IF new_id IS NULL THEN
		RAISE EXCEPTION 'Parameter value cannot be null. A closed task cannot be modified. ';
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
CREATE FUNCTION f_glossary_delete(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	IF(json_data ->> 'id')::BIGINT = 0 THEN
		TRUNCATE TABLE glossary CASCADE;
		--DELETE FROM glossary;
	ELSE
		DELETE FROM glossary WHERE id = (json_data ->> 'id')::BIGINT; 
	END IF;

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
CREATE FUNCTION f_glossary_view(
		json_data jsonb
) 
RETURNS TABLE (
	id BIGINT,
	mark_id BIGINT,
	mark_name TEXT,
	letter_id BIGINT,
	letter_name TEXT,
	word TEXT, 
	transcription TEXT, 
	translation TEXT, 
	example JSONB,
    dt_add BIGINT,
	dt_add_txt TEXT,
	enable BOOL
) AS $$
DECLARE
  	par_id BIGINT = 0;
	par_mark_id BIGINT = 0;
	par_mark_name TEXT = '';
	par_letter_name TEXT = '';
	par_letter_id BIGINT = 0;

BEGIN

	IF (json_data ->> 'id') IS NOT NULL THEN
		par_id = (json_data ->> 'id')::BIGINT;
	END IF;

	IF (json_data ->> 'mark_name') IS NOT NULL THEN
		par_mark_name = (json_data ->> 'mark_name')::BIGINT;

		SELECT mark.id
		INTO   par_mark_id
		FROM   mark
		WHERE  (par_mark_name = mark.name);
	END IF;

	IF (json_data ->> 'letter_name') IS NOT NULL THEN
		par_letter_name = (json_data ->> 'letter_name')::BIGINT;

		SELECT letter.id
		INTO   par_letter_id
		FROM   letter
		WHERE  (par_letter_name = letter.name);
	END IF;

	RETURN QUERY
		SELECT 
			glossary.id as id,
			glossary.mark_id as mark_id,
			COALESCE((
					SELECT mark.name 
					FROM   mark 
					WHERE  mark.id = glossary.mark_id
				), '') as mark_name,
			glossary.letter_id as letter_id,
			COALESCE((
					SELECT letter.name 
					FROM   letter 
					WHERE  letter.id = glossary.letter_id
				), '') as letter_name,

			glossary.word as word, 
			glossary.transcription as transcription,
			glossary.translation as  translation,
			glossary.example as example,
		    glossary.dt_add as dt_add,
			COALESCE(TO_CHAR(TO_TIMESTAMP(glossary.dt_add), 'DD.MM.YYYY HH24:MI:SS.MS'), '') as dt_add_txt, 
			glossary.enable as enable
		FROM 
			glossary
		WHERE
			(par_id = 0 OR glossary.id = par_id) AND
			(par_mark_id = 0 OR glossary.mark_id = par_mark_id) AND 
			(par_letter_id = 0 OR glossary.letter_id = par_letter_id)
		ORDER BY 
			glossary.id;
	
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION f_glossary_view_random(
    json_data jsonb
) 
RETURNS TABLE (
    id BIGINT,
    mark_id BIGINT,
    mark_name TEXT,
    letter_id BIGINT,
    letter_name TEXT,
    word TEXT, 
    transcription TEXT, 
    translation TEXT, 
    example JSONB,
    dt_add BIGINT,
    dt_add_txt TEXT,
    enable BOOL
) AS $$
DECLARE
    par_id BIGINT = 0;
    par_mark_id BIGINT = 0;
    par_mark_name TEXT = '';
    par_letter_name TEXT = '';
    par_letter_id BIGINT = 0;
    par_limit INT = 10; -- Значение по умолчанию
BEGIN
    -- Парсинг параметров из JSON
    IF (json_data ->> 'id') IS NOT NULL THEN
        par_id = (json_data ->> 'id')::BIGINT;
    END IF;

    IF (json_data ->> 'mark_name') IS NOT NULL THEN
        par_mark_name = (json_data ->> 'mark_name')::TEXT;

        SELECT mark.id
        INTO   par_mark_id
        FROM   mark
        WHERE  (par_mark_name = mark.name);
    END IF;

    IF (json_data ->> 'letter_name') IS NOT NULL THEN
        par_letter_name = (json_data ->> 'letter_name')::TEXT;

        SELECT letter.id
        INTO   par_letter_id
        FROM   letter
        WHERE  (par_letter_name = letter.name);
    END IF;

    IF (json_data ->> 'limit') IS NOT NULL THEN
        par_limit = (json_data ->> 'limit')::INT;
    END IF;

    -- Возврат случайных записей
    RETURN QUERY
        SELECT 
            glossary.id as id,
            glossary.mark_id as mark_id,
            COALESCE((
                    SELECT mark.name 
                    FROM   mark 
                    WHERE  mark.id = glossary.mark_id
                ), '') as mark_name,
            glossary.letter_id as letter_id,
            COALESCE((
                    SELECT letter.name 
                    FROM   letter 
                    WHERE  letter.id = glossary.letter_id
                ), '') as letter_name,
            glossary.word as word, 
            glossary.transcription as transcription,
            glossary.translation as  translation,
            glossary.example as example,
            glossary.dt_add as dt_add,
            COALESCE(TO_CHAR(TO_TIMESTAMP(glossary.dt_add), 'DD.MM.YYYY HH24:MI:SS.MS'), '') as dt_add_txt, 
            glossary.enable as enable
        FROM 
            glossary
        WHERE
            (par_id = 0 OR glossary.id = par_id) AND
            (par_mark_id = 0 OR glossary.mark_id = par_mark_id) AND 
            (par_letter_id = 0 OR glossary.letter_id = par_letter_id)
        ORDER BY 
            RANDOM()
        LIMIT 
            par_limit;
END;
$$ LANGUAGE plpgsql;



