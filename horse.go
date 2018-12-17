package horse

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"os"
)

// OperationsToMatch takes two definitions and returns the required
// operations to make the source definition match the target.
func OperationsToMatch(source, target Definition) ([]Operation, error) {
	return matchDef(source, target)
}

// NewDatabase returns a Database for a type.
func NewDatabase(d DatabaseType, db *sql.DB) (Database, error) {
	switch d {
	case Postgresql:
		return newPostgresqlDatabase(db)
	}

	return nil, ErrUnknownDatabase
}

// NewDefinitionFromJSON from JSON create a definition.
func NewDefinitionFromJSON(definition string) (Definition, error) {
	b := bytes.NewBuffer([]byte(definition))
	return newJSONDefinition(b)
}

// NewDefinitionFromFile loads a definition from a file.
func NewDefinitionFromFile(filename string) (Definition, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return newJSONDefinition(f)
}

func newJSONDefinition(r io.Reader) (Definition, error) {
	dec := json.NewDecoder(r)
	var d jsonDefinition
	if err := dec.Decode(&d); err != nil {
		return nil, err
	}

	if len(d.Schemas()) == 0 {
		return nil, ErrEmptyDefinition
	}

	return d, nil
}
