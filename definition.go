package horse

import "log"

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

// Operation an action to perform on an element.
type Operation struct {
	action Action
	schema *Schema
	table  *Table
	column *Column
}

func compare(source, target *Definition) ([]Operation, error) {
	Operations := []Operation{}

	for targetSchemaName, targetSchema := range target.Schemas {
		sourceSchema, ok := source.Schemas[targetSchemaName]
		if !ok {
			op := Operation{
				action: CreateSchema,
				schema: &targetSchema,
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
				schema: target,
				table:  &targetTable,
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
				schema: schema,
				table:  target,
				column: &targetColumn,
			}
			Operations = append(Operations, op)
			continue
		}
		alteration := false
		if targetColumn.Type != sourceColumn.Type {
			log.Println("change", "column", "type", targetColumn.Type, sourceColumn.Type)
			alteration = true
		}
		if targetColumn.DecimalSize == nil && sourceColumn.DecimalSize == nil {
			break
		} else if (targetColumn.DecimalSize == nil && sourceColumn.DecimalSize != nil) ||
			(targetColumn.DecimalSize != nil && sourceColumn.DecimalSize == nil) ||
			(*targetColumn.DecimalSize != *sourceColumn.DecimalSize) {
			log.Println("change", "column", "decimalsize", *targetColumn.DecimalSize, *sourceColumn.DecimalSize)
			alteration = true
		}
		if targetColumn.Length == nil && sourceColumn.DecimalSize == nil {
			break
		} else if (targetColumn.Length == nil && sourceColumn.Length != nil) ||
			(targetColumn.Length != nil && sourceColumn.Length == nil) ||
			(*targetColumn.Length != *sourceColumn.Length) {
			log.Println("change", "column", "length", *targetColumn.Length, *sourceColumn.Length)
			alteration = true
		}
		if targetColumn.Nullable != sourceColumn.Nullable {
			log.Println("change", "column", "nullable", targetColumn.Nullable, sourceColumn.Nullable)
			alteration = true
		}

		if alteration {
			op := Operation{
				action: AlterColumn,
				schema: schema,
				table:  target,
				column: &targetColumn,
			}
			Operations = append(Operations, op)
		}
	}

	return Operations, nil
}
