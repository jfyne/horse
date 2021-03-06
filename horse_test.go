package horse

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var db *sql.DB

var testDef = `{ "schemas": [ { "name": "test", "tables": [ {
		"name": "test1",
		"columns": [
			{
				"name": "first",
				"type": "text",
				"nullable": true
			},
			{
				"name": "second",
				"type": "decimal",
				"nullable": true
			},
			{
				"name": "third",
				"type": "integer",
				"nullable": true
			}
		]
	}]}]}`

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		panic(err)
	}

	defer teardown()
	os.Exit(m.Run())
}

func setup() error {
	testDB, err := sql.Open("postgres", "user=postgres password=max dbname=horse host=localhost sslmode=disable")
	if err != nil {
		return err
	}
	db = testDB
	return nil
}

func teardown() error {
	return db.Close()
}

func TestJSONDefintionGenerationWorks(t *testing.T) {
	_, err := NewDefinitionFromJSON(testDef)
	if err != nil {
		t.Error(err)
		return
	}

	bad := `{"aaa": "bbb", "ccc": {}}`

	if _, err := NewDefinitionFromJSON(bad); err == nil {
		t.Error("Bad JSON not flagged")
		return
	}
}

func TestDatabaseDefinitionGenerationWorks(t *testing.T) {
	database, err := NewDatabase(Postgresql, db)
	if err != nil {
		t.Error(err)
		return
	}

	if _, err := database.Definition("public"); err != nil {
		t.Error(err)
		return
	}
}
