package horse

import (
	"bytes"
	"encoding/json"
	"os"
)

// NewDefinitionFromJSON from JSON create a definition.
func NewDefinitionFromJSON(definition string) (*Definition, error) {
	b := bytes.NewBuffer([]byte(definition))
	dec := json.NewDecoder(b)
	var d Definition
	if err := dec.Decode(&d); err != nil {
		return nil, err
	}

	return &d, nil
}

// NewDefinitionFromFile loads a definition from a file.
func NewDefinitionFromFile(filename string) (*Definition, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var d Definition
	if err := dec.Decode(&d); err != nil {
		return nil, err
	}

	return &d, nil
}

// NewDescriptor returns a Descriptor for a type of database.
func NewDescriptor(d DatabaseType) (Descriptor, error) {
	switch d {
	case Postgresql:
		return newPostgresqlDescriptor()
	}

	return nil, ErrUnknownDatabase
}
