package main

import (
	"dictionary/pkg/storage/postgres"
	"log"
	"os"
)

func main() {

	var pool *postgres.Storage
	var err error
	connString := os.Getenv("PG_URL_DBLANGUAGE")

	pool, err = postgres.New(connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Выполнение SQL-скрипта
	// Удаляем все таблицы, функции, последовательности
	if err := pool.ExecuteSQLFromFile("../sql/deleteAllFunction.sql"); err != nil {
		log.Fatalf("Failed to execute SQL deleteAllFunction.sql: %v", err)
	}

	// Создаем структуру бд
	if err := pool.ExecuteSQLFromFile("../sql/script_1_glossary.sql"); err != nil {
		log.Fatalf("Failed to execute SQL script_1_glossary.sql: %v", err)
	}

	if err := pool.ExecuteSQLFromFile("../sql/script_2_catalogs.sql"); err != nil {
		log.Fatalf("Failed to execute SQL script_2_catalogs.sql: %v", err)
	}

	if err := pool.ExecuteSQLFromFile("../sql/script_3_users.sql"); err != nil {
		log.Fatalf("Failed to execute SQL script_3_users.sql: %v", err)
	}

	if err := pool.ExecuteSQLFromFile("../sql/script_4_records.sql"); err != nil {
		log.Fatalf("Failed to execute SQL script_4_records.sql: %v", err)
	}

	if err := pool.ExecuteSQLFromFile("../sql/script_5_catalogsglos.sql"); err != nil {
		log.Fatalf("Failed to execute SQL script_5_catalogsglos.sql: %v", err)
	}
	log.Println("All is done!!!")

}
