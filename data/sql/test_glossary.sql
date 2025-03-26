/*
	script to create database structure
*/

DROP FUNCTION IF EXISTS f_glossary_view_random;

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
