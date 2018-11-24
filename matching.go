package horse

import (
	"log"
)

// Operation an action to perform on an element.
type Operation struct {
	action Action
	schema Schema
	table  Table
	column Column
}

// Compare this source definition to a target, returns the
// Operations required in order to match the target.
func matchDef(source, target Definition) ([]Operation, error) {
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
		ops, err := compareTables(source, &sourceSchema, &targetSchema)
		if err != nil {
			return nil, err
		}
		Operations = append(Operations, ops...)
	}

	return Operations, nil
}

func compareTables(sourceDef Definition, source, target *Schema) ([]Operation, error) {
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
		ops, err := compareColumns(sourceDef, &sourceTable, target, &targetTable)
		if err != nil {
			return nil, err
		}
		Operations = append(Operations, ops...)
	}

	return Operations, nil
}

func compareColumns(sourceDef Definition, source *Table, schema *Schema, target *Table) ([]Operation, error) {
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

		// Type
		expectedType, err := sourceDef.ExpectedType(targetColumn.Type)
		if err != nil {
			return nil, err
		}
		if expectedType != sourceColumn.Type {
			log.Println("change", "column", "type", targetColumn.Name, targetColumn.Type, sourceColumn.Type)
			alteration = true
		}

		// Length
		if targetColumn.Length != sourceColumn.Length {
			log.Println("change", "column", "length", targetColumn.Name, targetColumn.Length, sourceColumn.Length)
			alteration = true
		}

		// Precision
		// If precision is blank, we are taking teh default. If it is something other than blank we are
		// potentially changing.
		if targetColumn.Precision != "" && (targetColumn.Precision != sourceColumn.Precision) {
			log.Println("change", "column", "precision", targetColumn.Name, targetColumn.Precision, sourceColumn.Precision)
			alteration = true
		}

		// Nullable
		if targetColumn.Nullable != sourceColumn.Nullable {
			log.Println("change", "column", "nullable", targetColumn.Name, targetColumn.Nullable, sourceColumn.Nullable)
			alteration = true
		}

		// Default
		if targetColumn.Default != sourceColumn.Default {
			log.Println("change", "column", "default", targetColumn.Name, targetColumn.Default, sourceColumn.Default)
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
