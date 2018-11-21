package horse

import (
	"database/sql"
)

// DatabaseType what database we are regerring to.
type DatabaseType string

const (
	// Postgresql the postgresql database.
	Postgresql DatabaseType = "postgres"
)

// Descriptor interface, describes a database element.
type Descriptor interface {

	// Schema return the named schema.
	Schema(*sql.DB, string) (Element, error)

	// Table return the named table.
	Table(*sql.DB, string, string) (Element, error)

	// Column return the named column.
	Column(*sql.DB, string, string, string) (Element, error)

	// Definition creates a definition from the descriptor.
	Definition(*sql.DB, ...string) (*Definition, error)
}

// Diff the difference between to definitions.
type Diff interface {
}

// Element a datbase element.
type Element interface {

	// Create the definition from the element.
	Definition() (interface{}, error)

	String() string
}
