package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ExampleColumn struct {
	Ex1 string `json:"ex1"`
	Ex2 string `json:"ex2"`
}

type Glossary struct {
	Letter        string          `json:"letter"`
	Word          string          `json:"word"`
	Transcription string          `json:"transcription"`
	Translation   string          `json:"translation"`
	Example       []ExampleColumn `json:"example"`
}

type GlossaryLoad struct {
	TableName    string
	TableComment string
	Data         []Glossary
}

func main() {

	connString := os.Getenv("PG_URL_DBLANGUAGE")

	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Чтение JSON-файла
	filePath := "glossary.json"
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var data GlossaryLoad
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		log.Fatal(err)
	}

	// загружаем данные в БД
	var exists bool
	tableName := fmt.Sprintf("catalog_%v", data.TableName)
	tableComment := data.TableComment

	// 1. проверяем наличие таблицы каталога слов
	idTable, err := createOrUpdateTable(pool, tableName, tableComment)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v, %#v, %#v, %#v\n", err, exists, tableName, tableComment)

	// 2. вставляем новые слова в glossary и заполняем соответствующий каталог слов catalog_xxx
	for i, data := range data.Data {
		err = addWordIntoGlossaryTable(pool, tableName, tableComment, data, idTable)
		if err != nil {
			log.Println(i, data.Word, err)
		} else {
			log.Println(i, data.Word, "ok!")
		}

	}

}

// 1. проверяем наличие таблицы каталога слов
func createOrUpdateTable(pool *pgxpool.Pool, tableName, tableComment string) (int, error) {

	// Первый запрос: вставка данных в таблицу glossary
	var newID int
	queryGlossary := `
		INSERT INTO catalogs (name, comment)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id;
	`
	err := pool.QueryRow(context.Background(), queryGlossary, tableName, tableComment).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil

}

// // 2. вставляем новые слова в glossary и заполняем соответствующий каталог слов catalog_xxx
func addWordIntoGlossaryTable(pool *pgxpool.Pool, tableName, tableComment string, data Glossary, idTable int) error {
	// Сериализация поля Example в JSON
	jsonData, err := json.Marshal(data.Example)
	if err != nil {
		return err
	}

	// Первый запрос: вставка данных в таблицу glossary
	var newID int64
	queryGlossary := `
		INSERT INTO glossary (letter, word, transcription, translation, example)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`
	err = pool.QueryRow(context.Background(), queryGlossary, data.Letter, data.Word, data.Transcription, data.Translation, jsonData).Scan(&newID)
	if err != nil {
		return err
	}

	// Второй запрос: вставка glossary_id в указанную таблицу

	queryGlossary = `
		INSERT INTO catalogsglos (glossary_id, catalogs_id)
		VALUES ($1, $2);
	`
	_, err = pool.Exec(context.Background(), queryGlossary, newID, idTable)
	if err != nil {
		return err
	}

	return nil
}
