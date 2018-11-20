package horse

import (
	"testing"
)

func TestNew(t *testing.T) {
	d, err := newPostgresqlDescriptor()
	if err != nil {
		t.Error(err)
		return
	}

	ps, err := d.Schemas(db)
	if err != nil {
		t.Error(err)
		return
	}

	t.Error(ps)
}
