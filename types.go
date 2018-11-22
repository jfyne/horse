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

// Element a database element.
type Element interface {

	// Create the definition from the element.
	Definition() (interface{}, error)

	String() string
}
