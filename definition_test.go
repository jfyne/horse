package horse

import (
	"testing"
)

func TestCompare(t *testing.T) {
	_, err := NewDefinitionFromFile("example-definition.json")
	if err != nil {
		t.Error(err)
		return
	}
}
