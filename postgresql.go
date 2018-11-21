package horse

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type postgresDescriptor struct {
}

type postgresSchema struct {
	CatalogName                string `db:"catalog_name"`
	SchemaName                 string `db:"schema_name"`
	SchemaOwner                string `db:"schema_owner"`
	DefaultCharacterSetCatalog string `db:"default_charcter_set_catalog"`
	DefaultCharacterSetSchema  string `db:"default_charcater_set_schema"`
	DefaultCharacterSetName    string `db:"default_character_set_name"`
	SQLPath                    string `db:"sql_path"`
}

func (p postgresSchema) String() string {
	return fmt.Sprintf(`
schema (
	catalog_name: %s,
	schema_name:  %s,
	schema_owner: %s
)`, p.CatalogName, p.SchemaName, p.SchemaOwner)
}

func (p postgresSchema) Definition() (interface{}, error) {
	schema := Schema{
		Name:   p.SchemaName,
		Tables: map[string]Table{},
	}
	return schema, nil
}

type postgresTable struct {
	TableCatalog string `db:"table_catalog"`
	TableSchema  string `db:"table_schema"`
	TableName    string `db:"table_name"`
	TableType    string `db:"table_type"`
}

func (p postgresTable) String() string {
	return fmt.Sprintf(`
table (
	table_catalog: %s,
	table_schema:  %s,
	table_name:    %s,
	talbe_type:    %s
}`, p.TableCatalog, p.TableSchema, p.TableName, p.TableType)
}

func (p postgresTable) Definition() (interface{}, error) {
	table := Table{
		Name:    p.TableName,
		Columns: map[string]Column{},
	}
	return table, nil
}

type postgresColumn struct {
	TableCatalog          string  `db:"table_catalog"`
	TableSchema           string  `db:"table_schema"`
	TableName             string  `db:"table_name"`
	ColumnName            string  `db:"column_name"`
	OrdinalPosition       int64   `db:"ordinal_position"`
	ColumnDefault         *string `db:"column_default"`
	IsNullable            string  `db:"is_nullable"`
	DataType              string  `db:"data_type"`
	CharcterMaximumLength *int64  `db:"character_maximum_length"`
	CharacterOctetLength  *int64  `db:"character_octet_length"`
	NumericPrecision      *int64  `db:"numeric_precision"`
	NumericPrecisionRadix *int64  `db:"numeric_precision_radix"`
	NumericScale          *int64  `db:"numeric_scale"`
	DatetimePrecisions    *int64  `db:"datetime_precision"`
	IntervalType          *string `db:"interval_type"`
	IntervalPrecision     *int64  `db:"interval_precision"`
	CharacterSetCatalog   *string `db:"character_set_catalog"`
	CharacterSetSchema    *string `db:"character_set_schema"`
	CharacterSetName      *string `db:"character_set_name"`
	CollationCatalog      *string `db:"collation_catalog"`
	CollationSchema       *string `db:"collation_schema"`
	CollationName         *string `db:"collation_name"`
	DomainCatalog         *string `db:"domain_catalog"`
	DomainSchema          *string `db:"domain_schema"`
	DomainName            *string `db:"domain_name"`
	UdtCatalog            string  `db:"udt_catalog"`
	UdtSchema             string  `db:"udt_schema"`
	UdtName               string  `db:"udt_name"`
	ScopeCatalog          *string `db:"scope_catalog"`
	ScopeSchema           *string `db:"scope_schema"`
	ScopeName             *string `db:"scope_name"`
	MaximumCardinality    *int64  `db:"maximum_cardinality"`
	DtdIdentifier         string  `db:"dtd_identifier"`
	IsSelfReferencing     string  `db:"is_self_referencing"`
	IsIdentity            string  `db:"is_identity"`
	IdentityGeneration    *string `db:"identity_generation"`
	IdentityStart         *string `db:"identity_start"`
	IdentityIncrement     *string `db:"identity_increment"`
	IdentityMaximum       *string `db:"identity_maximum"`
	IdentityMinimum       *string `db:"identity_minimum"`
	IdentityCycle         string  `db:"identity_cycle"`
	IsGenerated           string  `db:"is_generated"`
	GenerationExpression  *string `db:"generation_expression"`
	IsUpdatable           string  `db:"is_updatable"`
}

func (p postgresColumn) String() string {
	return fmt.Sprintf(`
column (
	table_schema: %s,
	table_name:   %s,
	column_name:  %s,
	column_type:  %s
)`, p.TableSchema, p.TableName, p.ColumnName, p.DataType)
}

func (p postgresColumn) Definition() (interface{}, error) {
	nullable := false
	if p.IsNullable == "YES" {
		nullable = true
	}
	column := Column{
		Name:     p.ColumnName,
		Type:     p.DataType,
		Nullable: nullable,
	}
	return column, nil
}

func newPostgresqlDescriptor() (Descriptor, error) {
	d := postgresDescriptor{}
	return d, nil
}

func (p postgresDescriptor) schemas(db *sql.DB, name string) ([]*postgresSchema, error) {
	dbx := sqlx.NewDb(db, "postgres")
	var rows *sqlx.Rows
	var err error

	q := `select
		catalog_name,
		schema_name,
		schema_owner
	from information_schema.schemata
	`
	if name != "" {
		where := `where schema_name = $1`
		rows, err = dbx.Queryx(q+where, name)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = dbx.Queryx(q)
		if err != nil {
			return nil, err
		}
	}

	defer rows.Close()

	schemas := []*postgresSchema{}
	for rows.Next() {
		var ps postgresSchema
		if err := rows.StructScan(&ps); err != nil {
			return nil, err
		}
		schemas = append(schemas, &ps)
	}

	return schemas, nil
}

func (p postgresDescriptor) schema(db *sql.DB, name string) (*postgresSchema, error) {
	elements, err := p.schemas(db, name)
	if err != nil {
		return nil, err
	}
	if len(elements) == 0 {
		return nil, ErrNotFound
	}
	if len(elements) != 1 {
		return nil, ErrTooManyResults
	}
	return elements[0], nil
}

func (p postgresDescriptor) Schema(db *sql.DB, name string) (Element, error) {
	return p.schema(db, name)
}

func (p postgresDescriptor) tables(db *sql.DB, schema, name string) ([]*postgresTable, error) {
	dbx := sqlx.NewDb(db, "postgres")
	var rows *sqlx.Rows
	var err error

	q := `select
		table_catalog,
		table_schema,
		table_name,
		table_type
	from information_schema.tables
	where table_schema = $1
	`
	if name != "" {
		where := `and table_name = $2`
		rows, err = dbx.Queryx(q+where, schema, name)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = dbx.Queryx(q, schema)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	tables := []*postgresTable{}
	for rows.Next() {
		var pt postgresTable
		if err := rows.StructScan(&pt); err != nil {
			return nil, err
		}
		tables = append(tables, &pt)

	}
	return tables, nil
}

func (p postgresDescriptor) table(db *sql.DB, schema, name string) (*postgresTable, error) {
	elements, err := p.tables(db, schema, name)
	if err != nil {
		return nil, err
	}
	if len(elements) == 0 {
		return nil, ErrNotFound
	}
	if len(elements) != 1 {
		return nil, ErrTooManyResults
	}
	return elements[0], nil
}

func (p postgresDescriptor) Table(db *sql.DB, schema, name string) (Element, error) {
	return p.table(db, schema, name)
}

func (p postgresDescriptor) columns(db *sql.DB, schema, table, column string) ([]*postgresColumn, error) {
	dbx := sqlx.NewDb(db, "postgres")
	var rows *sqlx.Rows
	var err error

	q := `select
		*
	from information_schema.columns
	where table_schema = $1
	and table_name = $2
	`
	if column != "" {
		where := `and column_name = $3`
		rows, err = dbx.Queryx(q+where, schema, table, column)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = dbx.Queryx(q, schema, table)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	columns := []*postgresColumn{}
	for rows.Next() {
		var pc postgresColumn
		if err := rows.StructScan(&pc); err != nil {
			return nil, err
		}
		columns = append(columns, &pc)
	}
	return columns, nil
}

func (p postgresDescriptor) column(db *sql.DB, schema, table, column string) (*postgresColumn, error) {
	elements, err := p.columns(db, schema, table, column)
	if err != nil {
		return nil, err
	}
	if len(elements) == 0 {
		return nil, ErrNotFound
	}
	if len(elements) != 1 {
		return nil, ErrTooManyResults
	}
	return elements[0], nil
}

func (p postgresDescriptor) Column(db *sql.DB, schema, table, column string) (Element, error) {
	return p.column(db, schema, table, column)
}

func (p postgresDescriptor) Definition(db *sql.DB, schemas ...string) (*Definition, error) {
	d := Definition{
		Schemas: map[string]Schema{},
	}

	for _, schemaName := range schemas {
		schema, err := p.schema(db, schemaName)
		if err != nil {
			return nil, err
		}
		s, err := schema.Definition()
		if err != nil {
			return nil, err
		}
		sdef, _ := s.(Schema)
		d.Schemas[sdef.Name] = sdef

		tables, err := p.tables(db, schemaName, "")
		if err != nil {
			return nil, err
		}
		for _, table := range tables {
			t, err := table.Definition()
			if err != nil {
				return nil, err
			}
			tdef, _ := t.(Table)
			d.Schemas[sdef.Name].Tables[tdef.Name] = tdef

			columns, err := p.columns(db, schemaName, tdef.Name, "")
			if err != nil {
				return nil, err
			}

			for _, column := range columns {
				c, err := column.Definition()
				if err != nil {
					return nil, err
				}
				cdef, _ := c.(Column)
				d.Schemas[sdef.Name].Tables[tdef.Name].Columns[cdef.Name] = cdef
			}
		}
	}

	return &d, nil
}
