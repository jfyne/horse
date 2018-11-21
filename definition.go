package horse

// Definition a description of a state, either desired or current.
type Definition struct {
	Schemas map[string]Schema `json:"schemas"`
}

// Schema a database schema element.
type Schema struct {
	Name   string           `json:"name"`
	Tables map[string]Table `json:"tables"`
}

// Table a database table element.
type Table struct {
	Name    string            `json:"name"`
	Columns map[string]Column `json:"columns"`
}

// Column a database column element.
type Column struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	DecimalSize *int   `json:"decimalSize"`
	Length      *int   `json:"length"`
	Nullable    bool   `json:"nullable"`
}

// Compare two definitions to get changes.
func (d Definition) Compare(def Definition) (Diff, error) {
	return nil, nil
}
