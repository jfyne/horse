package horse

import "log"

// StdDefinition a description of a state, either desired or current.
type StdDefinition struct {
	stdSchemas map[string]Schema `json:"schemas"`
}

// Schemas get the schemas.
func (s StdDefinition) Schemas() map[string]Schema {
	return s.stdSchemas
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
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
	Default  string `json:"default"`
}

// Operation an action to perform on an element.
type Operation struct {
	action Action
	schema Schema
	table  Table
	column Column
}

func compare(source, target Definition) ([]Operation, error) {
	Operations := []Operation{}

	sourceSchemas := source.Schemas()
	for targetSchemaName, targetSchema := range target.Schemas() {
		sourceSchema, ok := sourceSchemas[targetSchemaName]
		if !ok {
			op := Operation{
				action: CreateSchema,
				schema: targetSchema,
			}
			Operations = append(Operations, op)
		}
		ops, err := compareTables(&sourceSchema, &targetSchema)
		if err != nil {
			return nil, err
		}
		Operations = append(Operations, ops...)
	}

	return Operations, nil
}

func compareTables(source, target *Schema) ([]Operation, error) {
	Operations := []Operation{}

	for targetTableName, targetTable := range target.Tables {
		sourceTable, ok := source.Tables[targetTableName]
		if !ok {
			op := Operation{
				action: CreateTable,
				schema: *target,
				table:  targetTable,
			}
			Operations = append(Operations, op)
		}
		ops, err := compareColumns(&sourceTable, target, &targetTable)
		if err != nil {
			return nil, err
		}
		Operations = append(Operations, ops...)
	}

	return Operations, nil
}

func compareColumns(source *Table, schema *Schema, target *Table) ([]Operation, error) {
	Operations := []Operation{}

	for targetColumnName, targetColumn := range target.Columns {
		sourceColumn, ok := source.Columns[targetColumnName]
		if !ok {
			op := Operation{
				action: CreateColumn,
				schema: *schema,
				table:  *target,
				column: targetColumn,
			}
			Operations = append(Operations, op)
			continue
		}
		alteration := false
		if targetColumn.Type != sourceColumn.Type {
			log.Println("change", "column", "type", targetColumn.Type, sourceColumn.Type)
			alteration = true
		}
		if targetColumn.Nullable != sourceColumn.Nullable {
			log.Println("change", "column", "nullable", targetColumn.Nullable, sourceColumn.Nullable)
			alteration = true
		}
		if targetColumn.Default != sourceColumn.Default {
			log.Println("change", "column", "default", targetColumn.Default, sourceColumn.Default)
			alteration = true
		}

		if alteration {
			op := Operation{
				action: AlterColumn,
				schema: *schema,
				table:  *target,
				column: targetColumn,
			}
			Operations = append(Operations, op)
		}
	}

	return Operations, nil
}
