package horse

import (
	"database/sql"
)

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
	Schemas() map[string]Schema

	// ExpectedType takes a target tpye and returns the expected type, some
	// databases alias type names.
	ExpectedType(target string) (string, error)
}
