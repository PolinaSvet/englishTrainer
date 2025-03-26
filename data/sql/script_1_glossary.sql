/*
	script to create database structure
*/

-- Сбросьте значения всех последовательностей
/*DO $$DECLARE r RECORD;
BEGIN
  FOR r IN (SELECT sequence_name FROM information_schema.sequences WHERE sequence_schema = 'public') LOOP
    EXECUTE 'SELECT setval(''' || r.sequence_name || ''', 1, false)';
  END LOOP;
END$$;*/

DROP TABLE IF EXISTS glossary CASCADE;


CREATE TABLE glossary (
    id BIGSERIAL PRIMARY KEY,
	letter TEXT NOT NULL,
	word TEXT NOT NULL UNIQUE, 
	transcription TEXT NOT NULL, 
	translation TEXT NOT NULL, 
    dt_add BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT, 
	enable BOOL NOT NULL DEFAULT TRUE,
	example JSONB NOT NULL
);
COMMENT ON TABLE glossary IS 'Table for words studied';


DROP FUNCTION IF EXISTS f_glossary_insert, f_glossary_update, f_glossary_delete, f_glossary_view, f_glossary_view_random;

--=======================
--1::: table: glossary
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
BEGIN

	INSERT INTO glossary (
		letter, 
		word,
		transcription,
		translation,
		example,
		dt_add,
    	enable    
		) 
	VALUES (
		(json_data ->> 'letter')::TEXT,
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
BEGIN

	--RAISE NOTICE 'Текущий текст: %', json_data;

	UPDATE glossary SET 
		letter = (json_data ->> 'letter')::TEXT,
		word = (json_data ->> 'word')::TEXT,
		transcription = (json_data ->> 'transcription')::TEXT,
		translation = (json_data ->> 'translation')::TEXT,
		example = (json_data ->> 'example')::JSONB,
		dt_add = (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
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

	--IF(json_data ->> 'id')::BIGINT = 0 THEN
	--	TRUNCATE TABLE glossary CASCADE;
	--ELSE
	--	DELETE FROM glossary WHERE id = (json_data ->> 'id')::BIGINT; 
	--END IF;

	DELETE FROM glossary WHERE id = (json_data ->> 'id')::BIGINT;

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
    letter TEXT,
    word TEXT, 
    transcription TEXT, 
    translation TEXT, 
    example JSONB,
    dt_add BIGINT,
    enable BOOL,
	catalog_id BIGINT, --dop
	users_id BIGINT,   --dop
    records_id BIGINT, --dop
	attempt INT,       --dop
	guess BOOL         --dop
) AS $$
DECLARE
  	par_id BIGINT = 0;
BEGIN

	IF (json_data ->> 'id') IS NOT NULL THEN
		par_id = (json_data ->> 'id')::BIGINT;
	END IF;

	RETURN QUERY
		SELECT 
			glossary.id as id,
			glossary.letter as letter,
			glossary.word as word, 
			glossary.transcription as transcription,
			glossary.translation as  translation,
			glossary.example as example,
		    glossary.dt_add as dt_add,
			--COALESCE(TO_CHAR(TO_TIMESTAMP(glossary.dt_add), 'DD.MM.YYYY HH24:MI:SS.MS'), '') as dt_add_txt, 
			glossary.enable as enable,
			(0)::BIGINT as catalog_id,
            (0)::BIGINT as users_id,
            (0)::BIGINT as records_id,
            (0)::INT as attempt,
            FALSE as guess
		FROM 
			glossary
		WHERE
			(par_id = 0 OR glossary.id = par_id)
		ORDER BY 
			glossary.id;
	
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION f_glossary_view_random(
    json_data jsonb
) 
RETURNS TABLE (
    id BIGINT,
    letter TEXT,
    word TEXT, 
    transcription TEXT, 
    translation TEXT, 
    example JSONB,
    dt_add BIGINT,
    enable BOOL,
	catalog_id BIGINT, --dop
	users_id BIGINT,   --dop
    records_id BIGINT, --dop
	attempt INT,       --dop
	guess BOOL         --dop
) AS $$
DECLARE
    v_limit INT;
    v_users_id BIGINT;
    v_attempt INT;
    v_guess INT;
    v_catalog_id BIGINT;
BEGIN
    -- Устанавливаем значения параметров с проверкой границ и значений по умолчанию
    v_limit := COALESCE((json_data ->> 'limit')::INT, 5);
    v_limit := GREATEST(1, LEAST(100, v_limit)); -- Ограничиваем от 1 до 100
	v_catalog_id := COALESCE((SELECT catalogs.id FROM catalogs WHERE catalogs.id = (json_data ->> 'catalog_id')::BIGINT), 0);
    
    v_users_id := COALESCE((SELECT users.id FROM users WHERE users.id = (json_data ->> 'users_id')::BIGINT), 0);
    v_attempt := COALESCE((json_data ->> 'attempt')::INT, 0);
    v_guess := COALESCE((json_data ->> 'guess')::INT, 0);
	
    
    IF v_users_id = 0 THEN
            -- Вариант без пользователем - используем catalogsglos
            RETURN QUERY
            SELECT 
                g.id,
                g.letter,
                g.word, 
                g.transcription, 
                g.translation, 
                g.example,
                g.dt_add,
                g.enable,
                v_catalog_id,
                v_users_id,
                (0)::BIGINT,
                0,
                FALSE
            FROM 
                glossary g
            LEFT JOIN 
                catalogsglos cg ON g.id = cg.glossary_id AND cg.catalogs_id = v_catalog_id
            WHERE 
                g.enable = TRUE
                AND (v_catalog_id = 0 OR v_catalog_id = cg.catalogs_id)
            ORDER BY 
                random()
            LIMIT 
                v_limit;
    ELSE
            -- Вариант с пользователем - используем records
            RETURN QUERY
            WITH filtered_records AS (
                SELECT 
                    r.glossary_id,
                    r.attempt,
                    r.guess,
                    r.catalogs_id,
                    r.users_id,
                    r.id,
                    CASE 
                        WHEN v_attempt = 1 THEN rank() OVER (PARTITION BY r.glossary_id ORDER BY r.attempt DESC)
                        WHEN v_attempt = 2 THEN rank() OVER (PARTITION BY r.glossary_id ORDER BY r.attempt ASC)
                        ELSE 1
                    END as attempt_rank
                FROM 
                    records r
                WHERE 
                    r.users_id = v_users_id AND
                    (v_catalog_id = 0 OR r.catalogs_id = v_catalog_id) AND
                    (v_guess = 0 OR
                    (v_guess = 1 AND r.guess = FALSE) OR
                    (v_guess = 2 AND r.guess = TRUE))
            )
            SELECT 
                g.id,
                g.letter,
                g.word, 
                g.transcription, 
                g.translation, 
                g.example,
                g.dt_add,
                g.enable,
                fr.catalogs_id,
                fr.users_id,
                fr.id,
                fr.attempt,
                fr.guess
            FROM 
                glossary g
            JOIN
                filtered_records fr ON g.id = fr.glossary_id AND (v_attempt = 0 OR fr.attempt_rank = 1)
            WHERE 
                g.enable = TRUE
            ORDER BY
                CASE WHEN v_attempt = 0 THEN random() ELSE g.id END
            LIMIT 
                v_limit;
    END IF;	

END;
$$ LANGUAGE plpgsql;


