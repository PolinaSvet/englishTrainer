package main

import (
	"dictionary/pkg/storage/postgres"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {

	var pool *postgres.Storage
	var err error
	connString := os.Getenv("PG_URL_DBLANGUAGE")
	//connString := "postgres://postgres:root@localhost:5432/dbLanguage"

	pool, err = postgres.New(connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Выполнение SQL-скрипта
	if err := pool.ExecuteSQLFromFile("../sql/scriptCreateTable.sql"); err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	if err := pool.ExecuteSQLFromFile("../sql/scriptCreateFunctions.sql"); err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	if err := pool.ExecuteSQLFromFile("../sql/scriptInsertData.sql"); err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	// Чтение JSON-файла
	filePath := "glossary.json"
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Преобразование JSON в []map[string]interface{}
	var data []map[string]interface{}
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		log.Fatal(err)
	}

	// clear glossary
	jsonDataMap := make(map[string]interface{})
	jsonDataMap["id"] = 0
	id, err := pool.DeleteGlossary(jsonDataMap)
	if err != nil {
		log.Fatal(err, id)
	}

	// insert glossary
	fmt.Println(len(data))
	cntErr := 0
	for i, v := range data {
		id, err := pool.InsertGlossary(v)
		if err != nil || id == 0 {
			fmt.Println(i, v["word"], "err", err)
			cntErr++
		} else {
			fmt.Println(i, v["word"], "add")
		}
	}

	fmt.Println("Add:", (len(data) - cntErr), "Err:", cntErr)

}
