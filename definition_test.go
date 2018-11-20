package horse

import (
	"testing"
)

func TestCompare(t *testing.T) {
	d, err := NewDefinitionFromJSONFile("example-definition.json")
	if err != nil {
		t.Error(err)
		return
	}

	t.Error(d)
}
