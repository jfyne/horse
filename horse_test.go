package horse

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var db *sql.DB

var testDef = `{ "schemas": { "public": { "name": "public", "tables": { "test1": {
		"name": "test1",
		"columns": {
			"first": {
				"name": "first",
				"type": "text",
				"nullable": true
			},
			"second": {
				"name": "second",
				"type": "decimal",
				"nullable": true
			},
			"third": {
				"name": "third",
				"type": "integer",
				"nullable": true
			}
		}
	}}}}}`

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
	if _, err := db.Exec("create table if not exists test1 (first text, second decimal, third integer)"); err != nil {
		return err
	}
	return nil
}

func teardown() error {
	if _, err := db.Exec("drop table if exists test1"); err != nil {
		return err
	}
	return db.Close()
}

func TestJSON(t *testing.T) {
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

func TestDatabase(t *testing.T) {
	database, err := NewDatabase(Postgresql)
	if err != nil {
		t.Error(err)
		return
	}

	dbDefinition, err := database.Definition(db, "public")
	if err != nil {
		t.Error(err)
		return
	}

	jsonDefinition, err := NewDefinitionFromJSON(testDef)
	if err != nil {
		t.Error(err)
		return
	}

	ops, err := OperationsToMatch(dbDefinition, jsonDefinition)
	if err != nil {
		t.Error(err)
		return
	}

	if len(ops) != 0 {
		migs, _ := database.Migrations(db, ops)
		for _, m := range migs {
			t.Error(m)
		}
		t.Error("Definitions should match")
		return
	}
}
