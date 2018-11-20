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

func newPostgresqlDescriptor() (Descriptor, error) {
	d := postgresDescriptor{}
	return d, nil
}

func (p postgresDescriptor) Schema(db *sql.DB, name string) (Element, error) {
	dbx := sqlx.NewDb(db, "postgres")
	q := `select
		catalog_name,
		schema_name,
		schema_owner
	from information_schema.schemata
	where schema_name = $1`
	var ps postgresSchema

	row := dbx.QueryRowx(q, name)
	if err := row.StructScan(&ps); err != nil {
		return nil, err
	}

	return ps, nil
}

func (p postgresDescriptor) Table(db *sql.DB, schema, name string) (Element, error) {
	dbx := sqlx.NewDb(db, "postgres")
	q := `select
		table_catalog,
		table_schema,
		table_name,
		table_type
	from information_schema.tables
	where table_schema = $1
	and table_name = $2`
	var pt postgresTable

	row := dbx.QueryRowx(q, schema, name)
	if err := row.StructScan(&pt); err != nil {
		return nil, err
	}

	return pt, nil
}

func (p postgresDescriptor) Column(db *sql.DB, schema, table, column string) (Element, error) {
	dbx := sqlx.NewDb(db, "postgres")
	q := `select
		*
	from information_schema.columns
	where table_schema = $1
	and table_name = $2
	and column_name = $3`
	var pc postgresColumn

	row := dbx.QueryRowx(q, schema, table, column)
	if err := row.StructScan(&pc); err != nil {
		return nil, err
	}

	return pc, nil
}
