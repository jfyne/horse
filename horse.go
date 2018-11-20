package horse

// NewDescriptor returns a Descriptor for a type of database.
func NewDescriptor(d DatabaseType) (Descriptor, error) {
	switch d {
	case Postgresql:
		return newPostgresqlDescriptor()
	}

	return nil, ErrUnknownDatabase
}
