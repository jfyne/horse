package horse

import "errors"

// ErrUnknownDatabase returned when the type of database is unrecognised.
var ErrUnknownDatabase = errors.New("Unknown database type")
