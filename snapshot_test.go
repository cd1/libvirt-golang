package libvirt

import (
	"bytes"
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

func TestSnapshotRef(t *testing.T) {
	env := newTestEnvironment(t).withSnapshot()
	defer env.cleanUp()

	if err := env.snap.Ref(); err != nil {
		t.Fatal(err)
	}

	if err := env.snap.Free(); err != nil {
		t.Error(err)
	}
}

func TestSnapshotListChildren(t *testing.T) {
	env := newTestEnvironment(t).withSnapshot()
	defer env.cleanUp()

	var xml bytes.Buffer
	data := newTestSnapshotData()

	if err := testSnapshotTmpl.Execute(&xml, data); err != nil {
		t.Fatal(err)
	}

	childSnap, err := env.dom.CreateSnapshot(xml.String(), SnapCreateDefault)
	if err != nil {
		t.Fatal(err)
	}
	defer childSnap.Free()

	snapshots, err := env.snap.ListChildren(SnapListDescendants)
	if err != nil {
		t.Fatal(err)
	}

	for _, snap := range snapshots {
		defer snap.Free()
	}

	if l := len(snapshots); l != 1 {
		t.Errorf("unexpected snapshot children count; got=%v, want=1", l)
	}

	if childName := snapshots[0].Name(); childName != data.Name {
		t.Errorf("unexpected snapshot child name; got=%v, want=%v\n", childName, data.Name)
	}
}

func TestSnapshotRevert(t *testing.T) {
	env := newTestEnvironment(t).withSnapshot()
	defer env.cleanUp()

	if err := env.snap.Revert(SnapshotRevertFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid flag")
	}

	if err := env.snap.Revert(SnapRevertDefault); err != nil {
		t.Fatal(err)
	}
}
