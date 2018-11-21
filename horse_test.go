package horse

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var db *sql.DB

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
	good := `{ "schemas": { "public": { "name": "public", "tables": { "test1": {
		"name": "test1",
		"columns": {
			"first": {
				"name": "first",
				"type": "text"
			},
			"second": {
				"name": "second",
				"type": "decimal"
			},
			"third": {
				"name": "third",
				"type": "integer"
			}
		}
	}}}}}`

	_, err := NewDefinitionFromJSON(good)
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
