package horse

import "log"

type jsonDefinition struct {
	baseDefinition
}

// ExpectedType takes a target tpye and returns the expected type, some
// databases alias type names.
func (j jsonDefinition) ExpectedType(target string) (string, error) {
	log.Println("Std")
	return target, nil
}
