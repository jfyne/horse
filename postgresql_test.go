package horse

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func setupDBState() error {
	def, err := NewDefinitionFromJSON(testDef)
	if err != nil {
		return err
	}

	d, err := newPostgresqlDatabase()
	if err != nil {
		return err
	}

	dbDef, err := d.Definition(db, "public")
	if err != nil {
		return err
	}

	ops, err := OperationsToMatch(dbDef, def)
	if err != nil {
		return err
	}

	migrations, err := d.Migrations(db, ops)
	if err != nil {
		return err
	}

	if err := d.Migrate(db, migrations); err != nil {
		return err
	}

	return nil
}

func TestSchema(t *testing.T) {
	if err := setupDBState(); err != nil {
		t.Error(err)
		return
	}

	d, err := newPostgresqlDatabase()
	if err != nil {
		t.Error(err)
		return
	}

	if _, err := d.Schema(db, "publix"); err == nil {
		t.Error(err)
		return
	}

	s, err := d.Schema(db, "public")
	if err != nil {
		t.Error(err)
		return
	}
	ps, ok := s.(*postgresSchema)
	if !ok {
		t.Error("Cannot convert to postgres schema")
		return
	}
	if ps.SchemaName != "public" {
		t.Error("Unmatching schema name")
		return
	}
}

func TestTable(t *testing.T) {
	d, err := newPostgresqlDatabase()
	if err != nil {
		t.Error(err)
		return
	}

	ta, err := d.Table(db, "public", "test1")
	if err != nil {
		t.Error(err)
		return
	}
	pt, ok := ta.(*postgresTable)
	if !ok {
		t.Error("Cannot convert to postgres table")
		return
	}
	if pt.TableName != "test1" {
		t.Error("Unmatching table name")
		return
	}
}

func TestColumn(t *testing.T) {
	d, err := newPostgresqlDatabase()
	if err != nil {
		t.Error(err)
		return
	}

	co, err := d.Column(db, "public", "test1", "first")
	if err != nil {
		t.Error(err)
		return
	}
	pc, ok := co.(*postgresColumn)
	if !ok {
		t.Error("Cannot conver to postgres column")
		return
	}
	if pc.ColumnName != "first" {
		t.Error(pc)
		t.Error("Unmatching column name")
		return
	}
}

func TestDefinition(t *testing.T) {
	d, err := newPostgresqlDatabase()
	if err != nil {
		t.Error(err)
		return
	}

	def, err := d.Definition(db, "public")
	if err != nil {
		t.Error(err)
		return
	}

	schemas := def.Schemas()

	s, ok := schemas["public"]
	if !ok {
		t.Error("Schema public not present in definition")
		return
	}

	ta, ok := s.Tables["test1"]
	if !ok {
		t.Error("Table test1 not present in definition")
		return
	}

	col, ok := ta.Columns["second"]
	if !ok {
		t.Error("Column second not present in definition")
		return
	}

	if col.Name != "second" {
		t.Error("Column name not filled")
		return
	}

	if col.Type != "numeric" {
		t.Error("Column data type not filled")
		return
	}
}

func TestTypeMap(t *testing.T) {
	good := map[string]string{
		"decimal":      "numeric",
		"varchar":      "character varying",
		"varchar(255)": "character varying",
		"text":         "text",
	}

	pgdef := postgresDefinition{db: sqlx.NewDb(db, "postgres")}

	for source, target := range good {
		pgType, err := pgdef.ExpectedType(source)
		if err != nil {
			t.Error(err)
			return
		}
		if pgType != target {
			t.Error(target, source)
		}
	}

}
