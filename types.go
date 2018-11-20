package horse

import "database/sql"

// DatabaseType what database we are regerring to.
type DatabaseType string

const (
	// Postgresql the postgresql database.
	Postgresql DatabaseType = "postgres"
)

// Definition interface, defines a database element.
type Definition interface {
	String() string
}

// Descriptor interface, used to describe a database.
type Descriptor interface {
	Schemas(*sql.DB) (Definition, error)
}
