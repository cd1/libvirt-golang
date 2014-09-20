package libvirt

// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"
import (
	"log"
)

// SnapshotListFlag defines a filter when listing snapshots.
type SnapshotListFlag uint

// Possible values for SnapshotListFlag.
const (
	SnapListAll         SnapshotListFlag = 0
	SnapListDescendants SnapshotListFlag = C.VIR_DOMAIN_SNAPSHOT_LIST_DESCENDANTS
	SnapListRoots       SnapshotListFlag = C.VIR_DOMAIN_SNAPSHOT_LIST_ROOTS
	SnapListMetadata    SnapshotListFlag = C.VIR_DOMAIN_SNAPSHOT_LIST_METADATA
	SnapListLeaves      SnapshotListFlag = C.VIR_DOMAIN_SNAPSHOT_LIST_LEAVES
	SnapListNoLeaves    SnapshotListFlag = C.VIR_DOMAIN_SNAPSHOT_LIST_NO_LEAVES
	SnapListNoMetadata  SnapshotListFlag = C.VIR_DOMAIN_SNAPSHOT_LIST_NO_METADATA
	SnapListInactive    SnapshotListFlag = C.VIR_DOMAIN_SNAPSHOT_LIST_INACTIVE
	SnapListActive      SnapshotListFlag = C.VIR_DOMAIN_SNAPSHOT_LIST_ACTIVE
	SnapListDiskOnly    SnapshotListFlag = C.VIR_DOMAIN_SNAPSHOT_LIST_DISK_ONLY
	SnapListInternal    SnapshotListFlag = C.VIR_DOMAIN_SNAPSHOT_LIST_INTERNAL
	SnapListExternal    SnapshotListFlag = C.VIR_DOMAIN_SNAPSHOT_LIST_EXTERNAL
)

// Snapshot holds a libvirt domain snapshot. There are no exported fields.
type Snapshot struct {
	log         *log.Logger
	virSnapshot C.virDomainSnapshotPtr
}

// Free frees the domain snapshot object. The snapshot itself is not modified.
// The data structure is freed and should not be used thereafter.
func (snap Snapshot) Free() error {
	snap.log.Println("freeing snapshot object...")
	cRet := C.virDomainSnapshotFree(snap.virSnapshot)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		snap.log.Printf("an error occurred: %v\n", err)
		return err
	}

	snap.log.Println("snapshot freed")

	return nil
}
