package postgres

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	strConnection = "PG_URL_DBLANGUAGE"
)

func TestPostgres_New(t *testing.T) {
	t.Run("Postgres.New: Good (valid connection)", func(t *testing.T) {

		//os.Setenv("PG_URL_DBLANGUAGE", "postgres://postgres:root@localhost:5432/dbLanguage")

		connstr := os.Getenv(strConnection)
		assert.NotNil(t, connstr)

		store, err := New(connstr)
		assert.NoError(t, err)
		assert.NotNil(t, store)
		store.Close()
	})

	t.Run("Postgres.New: Error (invalid connection)", func(t *testing.T) {
		store, err := New("invalid_conn_string")
		assert.Error(t, err)
		assert.Nil(t, store)
	})
}

// проверяем таблицу Mark
func TestMark_AllFunctions(t *testing.T) {

	connstr := os.Getenv(strConnection)
	assert.NotNil(t, connstr)

	store, err := New(connstr)
	assert.NoError(t, err)
	defer store.Close()

	t.Run("Mark: Good (insert_update_view_delete)", func(t *testing.T) {

		//f_mark_insert
		jsonDataMap := make(map[string]interface{})
		jsonDataMap["name"] = "test_data"
		id_insert, err := store.InsertMark(jsonDataMap)
		assert.NoError(t, err)
		assert.NotEmpty(t, id_insert)

		//f_mark_update
		jsonDataMap = make(map[string]interface{})
		jsonDataMap["id"] = id_insert
		jsonDataMap["name"] = "test_data_update"
		id, err := store.UpdateMark(jsonDataMap)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		//f_mark_view
		jsonDataMap = make(map[string]interface{})
		jsonDataMap["id"] = id_insert
		mark, err := store.ViewMark(jsonDataMap)
		assert.NoError(t, err)
		assert.NotEmpty(t, mark)

		//f_mark_delete
		jsonDataMap = make(map[string]interface{})
		jsonDataMap["id"] = id_insert
		id, err = store.DeleteMark(jsonDataMap)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})

}

// проверяем таблицу Letter
func TestLetter_AllFunctions(t *testing.T) {

	connstr := os.Getenv(strConnection)
	assert.NotNil(t, connstr)

	store, err := New(connstr)
	assert.NoError(t, err)
	defer store.Close()

	t.Run("Letter: Good (insert_update_view_delete)", func(t *testing.T) {

		//f_letter_insert
		jsonDataMap := make(map[string]interface{})
		jsonDataMap["name"] = "test_data"
		id_insert, err := store.InsertLetter(jsonDataMap)
		assert.NoError(t, err)
		assert.NotEmpty(t, id_insert)

		//f_letter_update
		jsonDataMap = make(map[string]interface{})
		jsonDataMap["id"] = id_insert
		jsonDataMap["name"] = "test_data_update"
		id, err := store.UpdateLetter(jsonDataMap)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		//f_letter_view
		jsonDataMap = make(map[string]interface{})
		jsonDataMap["id"] = id_insert
		letter, err := store.ViewLetter(jsonDataMap)
		assert.NoError(t, err)
		assert.NotEmpty(t, letter)

		//f_letter_delete
		jsonDataMap = make(map[string]interface{})
		jsonDataMap["id"] = id_insert
		id, err = store.DeleteLetter(jsonDataMap)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})

}

// проверяем таблицу Glossary
func TestGlossary_AllFunctions(t *testing.T) {

	connstr := os.Getenv(strConnection)
	assert.NotNil(t, connstr)

	store, err := New(connstr)
	assert.NoError(t, err)
	defer store.Close()

	t.Run("Glossary: Good (insert_insertarray_update_view_delete)", func(t *testing.T) {

		//f_glossary_insert
		jsonDataMap := map[string]interface{}{
			"mark_name":     "Phrasal verbs",
			"letter_name":   "A",
			"word":          "Test_word_zzz",
			"transcription": "[əˈkaʊnt fɔː]",
			"translation":   "Давать объяснение, составлять долю от ч.-л.",
			"example": []map[string]interface{}{
				{"ex1": "Пример 1", "ex2": "Перевод 1"},
				{"ex1": "Пример 2", "ex2": "Перевод 2"},
			},
			"dt_add":     "1740652272547",
			"dt_add_txt": "DD.MM.YYYY HH24:MI:SS.MS",
			"enable":     true,
		}
		id_insert, err := store.InsertGlossary(jsonDataMap)
		assert.NoError(t, err)
		assert.NotEmpty(t, id_insert)

		//f_glossary_update
		jsonDataMap["id"] = id_insert
		jsonDataMap["word"] = "Account ForXXX"
		id, err := store.UpdateGlossary(jsonDataMap)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		// f_glossary_view_random
		jsonDataMap = make(map[string]interface{})
		jsonDataMap["limit"] = 5
		glossary, err := store.ViewRandomGlossary(jsonDataMap)
		assert.NoError(t, err)
		assert.NotEmpty(t, glossary)

		//f_glossary_view
		jsonDataMap = make(map[string]interface{})
		jsonDataMap["id"] = id_insert
		glossary, err = store.ViewGlossary(jsonDataMap)
		assert.NoError(t, err)
		assert.NotEmpty(t, glossary)

		//f_glossary_delete
		jsonDataMap = make(map[string]interface{})
		jsonDataMap["id"] = id_insert
		id, err = store.DeleteGlossary(jsonDataMap)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})

}
