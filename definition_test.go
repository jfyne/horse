package horse

import (
	"testing"
)

func TestCompare(t *testing.T) {
	def1, err := NewDefinitionFromFile("example-definition.json")
	if err != nil {
		t.Error(err)
		return
	}
	def2, err := NewDefinitionFromFile("example-definition.json")
	if err != nil {
		t.Error(err)
		return
	}

	ops, err := compare(def1, def2)
	if err != nil {
		t.Error(err)
		return
	}

	if len(ops) != 0 {
		t.Error("Comparision failed", ops)
		return
	}

	def2.stdSchemas["a"] = Schema{Name: "a"}
	ops, err = compare(def1, def2)
	if err != nil {
		t.Error(err)
		return
	}
	if len(ops) != 1 {
		t.Error("Comparison failed", ops)
		return
	}
}

func TestMigrations(t *testing.T) {
	des, err := NewDescriptor(Postgresql)
	if err != nil {
		t.Error(err)
		return
	}

	defDB, err := des.Definition(db, "public")
	if err != nil {
		t.Error(err)
		return
	}

	defF, err := NewDefinitionFromFile("example-definition.json")
	if err != nil {
		t.Error(err)
		return
	}

	ops, err := compare(defDB, defF)
	if err != nil {
		t.Error(err)
		return
	}

	migrations, err := des.Migrations(db, ops)
	if err != nil {
		t.Error(err)
		return
	}

	for _, m := range migrations {
		t.Error(m)
	}
}
