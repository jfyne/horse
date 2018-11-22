package horse

import "errors"

// ErrUnknownDatabase returned when the type of database is unrecognised.
var ErrUnknownDatabase = errors.New("Unknown database type")

// ErrNotFound returned when an element is not found.
var ErrNotFound = errors.New("Element not found")

// ErrTooManyResults returned when filtering is ambiguous.
var ErrTooManyResults = errors.New("Too many results")

// ErrEmptyDefinition returned when a definition is empty.
var ErrEmptyDefinition = errors.New("Empty definition")
