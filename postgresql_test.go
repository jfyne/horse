package horse

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func setupDBState() error {
	def, err := NewDefinitionFromJSON(testDef)
	if err != nil {
		return errors.Wrap(err, "NewDefinitionFromJSON")
	}

	d, err := NewDatabase(Postgresql, db)
	if err != nil {
		return errors.Wrap(err, "newPostgresqlDatabase")
	}

	if err := d.Migrate(def); err != nil {
		return errors.Wrap(err, "d.Migrate")
	}

	return nil
}

func TestSchema(t *testing.T) {
	if err := setupDBState(); err != nil {
		t.Error(err)
		return
	}

	d, err := NewDatabase(Postgresql, db)
	if err != nil {
		t.Error(err)
		return
	}

	if _, err := d.Schema("publix"); err == nil {
		t.Error(err)
		return
	}

	s, err := d.Schema("public")
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
	d, err := NewDatabase(Postgresql, db)
	if err != nil {
		t.Error(err)
		return
	}

	ta, err := d.Table("public", "test1")
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
	d, err := NewDatabase(Postgresql, db)
	if err != nil {
		t.Error(err)
		return
	}

	co, err := d.Column("public", "test1", "first")
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
	d, err := NewDatabase(Postgresql, db)
	if err != nil {
		t.Error(err)
		return
	}

	def, err := d.Definition("public")
	if err != nil {
		t.Error(err)
		return
	}

	s, err := def.Schema("public")
	if err == ErrNotFound {
		t.Error("Schema public not present in definition")
		return
	}

	ta, err := s.Table("test1")
	if err == ErrNotFound {
		t.Error("Table test1 not present in definition")
		return
	}

	col, err := ta.Column("second")
	if err == ErrNotFound {
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
