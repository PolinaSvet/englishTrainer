/*
	script to create database structure
*/

DROP TABLE IF EXISTS mark, letter, glossary;

CREATE TABLE mark (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
COMMENT ON TABLE mark IS 'Table for types of words studied';

CREATE TABLE letter (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
COMMENT ON TABLE letter IS 'Table for first letter of words studied';

CREATE TABLE glossary (
    id BIGSERIAL PRIMARY KEY,
	mark_id BIGINT REFERENCES mark(id) DEFAULT 0,
	letter_id BIGINT REFERENCES letter(id) DEFAULT 0,
	word TEXT NOT NULL UNIQUE, 
	transcription TEXT NOT NULL, 
	translation TEXT NOT NULL, 
    dt_add BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT, 
	enable BOOL NOT NULL DEFAULT TRUE,
	example JSONB NOT NULL
);
COMMENT ON TABLE glossary IS 'Table for words studied';



