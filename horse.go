package horse

import (
	"encoding/json"
	"os"
)

// NewDescriptor returns a Descriptor for a type of database.
func NewDescriptor(d DatabaseType) (Descriptor, error) {
	switch d {
	case Postgresql:
		return newPostgresqlDescriptor()
	}

	return nil, ErrUnknownDatabase
}

// NewDefinitionFromJSONFile loads a definition from a file.
func NewDefinitionFromJSONFile(filename string) (*Definition, error) {
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
