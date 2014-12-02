package libvirt

import (
	"testing"
)

func TestSnapshotInit(t *testing.T) {
	env := newTestEnvironment(t).withSnapshot()
	defer env.cleanUp()

	if name := env.snap.Name(); name != env.snapData.Name {
		t.Errorf("unexpected snapshot name; got=%v, want=%v", name, env.snapData.Name)
	}

	_, err := env.snap.Parent()
	if err == nil {
		t.Error("an error was not returned when querying the parent snapshot of a root snapshot")
	} else {
		virtErr := err.(*Error)
		if virtErr.Code != ErrNoDomainSnapshot {
			t.Error(err)
		}
	}

	if !env.snap.HasMetadata() {
		t.Error("snapshot should have metadata (but it does not)")
	}
}

func TestSnapshotXML(t *testing.T) {
	env := newTestEnvironment(t).withSnapshot()
	defer env.cleanUp()

	if _, err := env.snap.XML(DomainXMLFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid flag")
	}

	xml, err := env.snap.XML(DomXMLDefault)
	if err != nil {
		t.Error(err)
	}

	if len(xml) == 0 {
		t.Error("empty snapshot XML")
	}
}
