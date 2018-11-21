package horse

import (
	"testing"
)

func TestCompare(t *testing.T) {
	_, err := NewDefinitionFromJSONFile("example-definition.json")
	if err != nil {
		t.Error(err)
		return
	}
}
