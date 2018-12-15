package horse

import (
	"database/sql"
)

// Schema a database schema element.
type Schema struct {
	Name   string  `json:"name"`
	Tables []Table `json:"tables"`
}

// Table fetch an individual table.
func (s Schema) Table(name string) (Table, error) {
	for _, t := range s.Tables {
		if t.Name == name {
			return t, nil
		}
	}
	return Table{}, ErrNotFound
}

// Table a database table element.
type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

// Column fetch an individual column.
func (t Table) Column(name string) (Column, error) {
	for _, c := range t.Columns {
		if c.Name == name {
			return c, nil
		}
	}
	return Column{}, ErrNotFound
}

// Column a database column element.
type Column struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Length    int64  `json:"length"`
	Precision string `json:"precision"`
	Nullable  bool   `json:"nullable"`
	Default   string `json:"default"`
}

// baseDefinition a description of a state, either desired or current.
type baseDefinition struct {
	StdSchemas []Schema `json:"schemas"`
}

// Schemas get the schemas.
func (s baseDefinition) Schemas() []Schema {
	return s.StdSchemas
}

// Schema fetch an individual schema.
func (s baseDefinition) Schema(name string) (Schema, error) {
	for _, sch := range s.StdSchemas {
		if sch.Name == name {
			return sch, nil
		}
	}
	return Schema{}, ErrNotFound
}

// Action performed by horse.
type Action string

// Horse operation types.
const (
	CreateSchema Action = "create-schema"
	CreateTable  Action = "create-table"
	CreateColumn Action = "create-column"
	AlterColumn  Action = "alter-column"
)

// DatabaseType what database we are regerring to.
type DatabaseType string

// Horse supported datbases.
const (
	Postgresql DatabaseType = "postgres"
)

// Database interface, describes a database.
type Database interface {

	// Schema return the named schema.
	Schema(*sql.DB, string) (Element, error)

	// Table return the named table.
	Table(*sql.DB, string, string) (Element, error)

	// Column return the named column.
	Column(*sql.DB, string, string, string) (Element, error)

	// Definition creates a definition from the database.
	Definition(*sql.DB, ...string) (Definition, error)

	// Migrations takes a slice of Operations and converts them into
	// steps to complete the Operation.
	Migrations(*sql.DB, []Operation) ([]string, error)

	// Migrate takes migrations and in a transaction
	// migrates the database.
	Migrate(*sql.DB, []string) error
}

// Element a database element.
type Element interface {

	// Create the definition from the element.
	Definition() (interface{}, error)

	String() string
}

// Definition a description of a schema and its tables.
type Definition interface {

	// Get the schemas for the definition.
	Schemas() []Schema

	// Schema pull an individual schema.
	Schema(string) (Schema, error)

	// ExpectedType takes a target tpye and returns the expected type, some
	// databases alias type names.
	ExpectedType(target string) (string, error)
}
