package libvirt

import (
	"testing"
)

func TestStreamAbort(t *testing.T) {
	env := newTestEnvironment(t).withStream()
	defer env.cleanUp()

	if err := env.str.Abort(); err != nil {
		t.Fatal(err)
	}
}

func TestStreamRef(t *testing.T) {
	env := newTestEnvironment(t).withStream()
	defer env.cleanUp()

	if err := env.str.Ref(); err != nil {
		t.Error(err)
	}

	if err := env.str.Free(); err != nil {
		t.Error(err)
	}
}
