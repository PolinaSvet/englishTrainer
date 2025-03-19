/*SELECT * FROM f_glossary_insert_array(('
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
		')::jsonb);*/

/*
SELECT id, name
	FROM public.letter;

SELECT id, name
	FROM public.mark;
*/
SELECT *
	FROM public.glossary;

/*SELECT * FROM f_glossary_view('{
    "id": 0
     
}'::jsonb);*/
