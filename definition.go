package horse

// Column a database column element.
type Column struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
	Default  string `json:"default"`
}

// Table a database table element.
type Table struct {
	Name    string            `json:"name"`
	Columns map[string]Column `json:"columns"`
}

// Schema a database schema element.
type Schema struct {
	Name   string           `json:"name"`
	Tables map[string]Table `json:"tables"`
}

// baseDefinition a description of a state, either desired or current.
type baseDefinition struct {
	StdSchemas map[string]Schema `json:"schemas"`
}

// Schemas get the schemas.
func (s baseDefinition) Schemas() map[string]Schema {
	return s.StdSchemas
}
