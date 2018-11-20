package horse

import (
	"testing"
)

func TestDescriptor(t *testing.T) {
	d, err := newPostgresqlDescriptor()
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
	ps, ok := s.(postgresSchema)
	if !ok {
		t.Error("Cannot convert to postgres schema")
		return
	}
	if ps.SchemaName != "public" {
		t.Error("Unmatching schema name")
		return
	}

	ta, err := d.Table(db, "public", "test1")
	if err != nil {
		t.Error(err)
		return
	}
	pt, ok := ta.(postgresTable)
	if !ok {
		t.Error("Cannot convert to postgres table")
		return
	}
	if pt.TableName != "test1" {
		t.Error("Unmatching table name")
		return
	}

	co, err := d.Column(db, "public", "test1", "first")
	if err != nil {
		t.Error(err)
		return
	}
	pc, ok := co.(postgresColumn)
	if !ok {
		t.Error("Cannot conver to postgres column")
		return
	}
	if pc.ColumnName != "first" {
		t.Error(pc)
		t.Error("Unmatching column name")
		return
	}

	t.Error(ps)
	t.Error(pt)
	t.Error(pc)
}
