package horse

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

// NewDefinitionFromJSON from JSON create a definition.
func NewDefinitionFromJSON(definition string) (Definition, error) {
	b := bytes.NewBuffer([]byte(definition))
	return newDefinition(b)
}

// NewDefinitionFromFile loads a definition from a file.
func NewDefinitionFromFile(filename string) (Definition, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return newDefinition(f)
}

func newDefinition(r io.Reader) (Definition, error) {
	dec := json.NewDecoder(r)
	var d StdDefinition
	if err := dec.Decode(&d); err != nil {
		return nil, err
	}

	if len(d.Schemas()) == 0 {
		return nil, ErrEmptyDefinition
	}

	return d, nil
}

// NewDescriptor returns a Descriptor for a type of database.
func NewDescriptor(d DatabaseType) (Descriptor, error) {
	switch d {
	case Postgresql:
		return newPostgresqlDescriptor()
	}

	return nil, ErrUnknownDatabase
}
