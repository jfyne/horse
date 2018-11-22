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

	def2.Schemas["a"] = Schema{Name: "a"}
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
