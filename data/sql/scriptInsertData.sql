/*
	script generates data
*/

TRUNCATE TABLE letter, mark, glossary CASCADE;

-- Сбросьте значения всех последовательностей
DO $$DECLARE r RECORD;
BEGIN
  FOR r IN (SELECT sequence_name FROM information_schema.sequences WHERE sequence_schema = 'public') LOOP
    EXECUTE 'SELECT setval(''' || r.sequence_name || ''', 1, false)';
  END LOOP;
END$$;

--letter: insert
DO $$
DECLARE
	item_array TEXT[] := ARRAY['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'];
	item TEXT;
	i INT;
BEGIN

	--FOREACH item IN ARRAY item_array
    --LOOP
    --    INSERT INTO letter (name) VALUES (item);
    --END LOOP;

	FOREACH item IN ARRAY item_array
    LOOP
        PERFORM * FROM f_letter_insert(('{"name": "'||item||'"}')::jsonb);
    END LOOP;

END;
$$;

--mark: insert
DO $$
DECLARE
	item_array TEXT[] := ARRAY['Word', 'Phrasal verbs'];
	item TEXT;
BEGIN

	--FOREACH item IN ARRAY item_array
    --LOOP
    --    INSERT INTO mark (name) VALUES (item);
    --END LOOP;

	FOREACH item IN ARRAY item_array
    LOOP
        PERFORM * FROM f_mark_insert(('{"name": "'||item||'"}')::jsonb);
    END LOOP;

END;
$$;


--glossary: insert
DO $$
DECLARE
	test_data JSONB;
BEGIN
/*
	PERFORM * FROM f_glossary_insert(('{
				    "mark_name": "Phrasal verbs",
			        "letter_name": "A",
			        "word": "Тестовое слово",
			        "transcription": "Тестовая транскрипция",
			        "translation": "Test word",
			        "example": [{"ex1": "Пример 1", "ex2": "Перевод 1"}, {"ex1": "Пример 2", "ex2": "Перевод 2"}],
			        "dt_add": ' || (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT || ',
			        "dt_add_txt": "' || to_char(now(), 'DD.MM.YYYY HH24:MI:SS.MS') || '",
			        "enable": true
				}')::jsonb);*/
/*
	PERFORM * FROM f_glossary_insert_array(('
				[{
				    "mark_name": "Phrasal verbs",
			        "letter_name": "A",
			        "word": "Тестовое слово1",
			        "transcription": "Тестовая транскрипция1",
			        "translation": "Test word1",
			        "example": [{"ex1": "Пример 1", "ex2": "Перевод 1"}, {"ex1": "Пример 2", "ex2": "Перевод 2"}],
			        "dt_add": ' || (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT || ',
			        "dt_add_txt": "' || to_char(now(), 'DD.MM.YYYY HH24:MI:SS.MS') || '",
			        "enable": true
				},
				{
				    "mark_name": "Phrasal verbs",
			        "letter_name": "B",
			        "word": "Тестовое слово2",
			        "transcription": "Тестовая транскрипция2",
			        "translation": "Test word2",
			        "example": [{"ex1": "Пример 1", "ex2": "Перевод 1"}, {"ex1": "Пример 2", "ex2": "Перевод 2"}],
			        "dt_add": ' || (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT || ',
			        "dt_add_txt": "' || to_char(now(), 'DD.MM.YYYY HH24:MI:SS.MS') || '",
			        "enable": true
				}]
		')::jsonb);
*/
END;
$$;




