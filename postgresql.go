package horse

import "database/sql"

type postgreqlDescriptor struct {
}

func newPostgresqlDescriptor() (Descriptor, error) {
	d := postgreqlDescriptor{}
	return d, nil
}

func (p postgreqlDescriptor) Schemas(db *sql.DB) (Definition, error) {
	ps := postgreqlSchema{}
	return ps, nil
}

type postgreqlSchema struct {
}

func (p postgreqlSchema) String() string {
	return "schemas"
}
