package horse

import (
	"testing"
)

func TestCompare(t *testing.T) {
	def1, err := NewDefinitionFromFile("testfiles/testcompare-1.json")
	if err != nil {
		t.Error(err)
		return
	}
	def2, err := NewDefinitionFromFile("testfiles/testcompare-1.json")
	if err != nil {
		t.Error(err)
		return
	}

	ops, err := Compare(def1, def2)
	if err != nil {
		t.Error(err)
		return
	}

	if len(ops) != 0 {
		t.Error("Comparision failed", ops)
		return
	}

	def3, err := NewDefinitionFromFile("testfiles/testcompare-2.json")
	if err != nil {
		t.Error(err)
		return
	}

	ops, err = Compare(def1, def3)
	if err != nil {
		t.Error(err)
		return
	}

	if len(ops) == 0 {
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

	ops, err := Compare(defDB, defF)
	if err != nil {
		t.Error(err)
		return
	}

	migrations, err := des.Migrations(db, ops)
	if err != nil {
		t.Error(err)
		return
	}

	if len(migrations) != 6 {
		t.Error("Wrong number of migrations")
	}
}
