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

type operation struct {
	action Action
	schema *Schema
	table  *Table
	column *Column
}

func compare(source, target *Definition) ([]operation, error) {
	operations := []operation{}

	for targetSchemaName, targetSchema := range target.Schemas {
		sourceSchema, ok := source.Schemas[targetSchemaName]
		if !ok {
			op := operation{
				action: CreateSchema,
				schema: &targetSchema,
			}
			operations = append(operations, op)
			continue
		}
		ops, err := compareTables(&sourceSchema, &targetSchema)
		if err != nil {
			return nil, err
		}
		operations = append(operations, ops...)
	}

	return operations, nil
}

func compareTables(source, target *Schema) ([]operation, error) {
	operations := []operation{}

	for targetTableName, targetTable := range target.Tables {
		sourceTable, ok := source.Tables[targetTableName]
		if !ok {
			op := operation{
				action: CreateTable,
				table:  &targetTable,
			}
			operations = append(operations, op)
			continue
		}
		ops, err := compareColumns(&sourceTable, &targetTable)
		if err != nil {
			return nil, err
		}
		operations = append(operations, ops...)
	}

	return operations, nil
}

func compareColumns(source, target *Table) ([]operation, error) {
	operations := []operation{}

	for targetColumnName, targetColumn := range target.Columns {
		sourceColumn, ok := source.Columns[targetColumnName]
		if !ok {
			op := operation{
				action: CreateColumn,
				column: &targetColumn,
			}
			operations = append(operations, op)
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
			op := operation{
				action: AlterColumn,
				column: &targetColumn,
			}
			operations = append(operations, op)
		}
	}

	return operations, nil
}
