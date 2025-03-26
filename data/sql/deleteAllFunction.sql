/*DO $$ 
DECLARE
    r RECORD;
BEGIN
    FOR r IN 
        SELECT 
            p.proname, 
            pg_catalog.pg_get_function_result(p.oid) AS result_type,
            pg_catalog.pg_get_function_arguments(p.oid) AS arguments,
            n.nspname
        FROM 
            pg_catalog.pg_proc p
        JOIN 
            pg_catalog.pg_namespace n ON n.oid = p.pronamespace
        WHERE 
            n.nspname NOT IN ('pg_catalog', 'information_schema')  -- Исключаем системные функции
    LOOP
        EXECUTE format('DROP FUNCTION %I.%I(%s)', r.nspname, r.proname, r.arguments);
    END LOOP;
END $$;*/

-- Удаление всех пользовательских таблиц и индексов
DO $$ 
DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
        EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
    END LOOP;
END $$;

-- Удаление всех пользовательских функций
DO $$ 
DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT proname, nspname FROM pg_proc p JOIN pg_namespace n ON p.pronamespace = n.oid 
               WHERE nspname NOT IN ('pg_catalog', 'information_schema')) LOOP
        EXECUTE 'DROP FUNCTION IF EXISTS ' || quote_ident(r.nspname) || '.' || quote_ident(r.proname) || '() CASCADE';
    END LOOP;
END $$;

-- Удаление всех пользовательских индексов
DO $$ 
DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT indexname FROM pg_indexes WHERE schemaname = 'public') LOOP
        EXECUTE 'DROP INDEX IF EXISTS ' || quote_ident(r.indexname) || ' CASCADE';
    END LOOP;
END $$;