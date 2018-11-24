package horse

import "testing"

func TestOperationsToMatch(t *testing.T) {
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

	ops, err := OperationsToMatch(def1, def2)
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

	ops, err = OperationsToMatch(def1, def3)
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
	database, err := NewDatabase(Postgresql)
	if err != nil {
		t.Error(err)
		return
	}

	defDB, err := database.Definition(db, "public")
	if err != nil {
		t.Error(err)
		return
	}

	defF, err := NewDefinitionFromFile("example-definition.json")
	if err != nil {
		t.Error(err)
		return
	}

	ops, err := OperationsToMatch(defDB, defF)
	if err != nil {
		t.Error(err)
		return
	}

	migrations, err := database.Migrations(db, ops)
	if err != nil {
		t.Error(err)
		return
	}

	if len(migrations) != 6 {
		t.Error("Wrong number of migrations")
	}
}
