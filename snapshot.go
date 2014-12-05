package libvirt

// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"
import (
	"log"
	"reflect"
	"unicode/utf8"
	"unsafe"
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

//SnapshotCreateFlag defines how a snapshot should be created.
type SnapshotCreateFlag uint32

// Possible values for SnapshotCreateFlag.
const (
	SnapCreateDefault    SnapshotCreateFlag = 0
	SnapCreateRedefine   SnapshotCreateFlag = C.VIR_DOMAIN_SNAPSHOT_CREATE_REDEFINE
	SnapCreateCurrent    SnapshotCreateFlag = C.VIR_DOMAIN_SNAPSHOT_CREATE_CURRENT
	SnapCreateNoMetadata SnapshotCreateFlag = C.VIR_DOMAIN_SNAPSHOT_CREATE_NO_METADATA
	SnapCreateHalt       SnapshotCreateFlag = C.VIR_DOMAIN_SNAPSHOT_CREATE_HALT
	SnapCreateDiskOnly   SnapshotCreateFlag = C.VIR_DOMAIN_SNAPSHOT_CREATE_DISK_ONLY
	SnapCreateReuseExt   SnapshotCreateFlag = C.VIR_DOMAIN_SNAPSHOT_CREATE_REUSE_EXT
	SnapCreateQuiesce    SnapshotCreateFlag = C.VIR_DOMAIN_SNAPSHOT_CREATE_QUIESCE
	SnapCreateAtomic     SnapshotCreateFlag = C.VIR_DOMAIN_SNAPSHOT_CREATE_ATOMIC
	SnapCreateLive       SnapshotCreateFlag = C.VIR_DOMAIN_SNAPSHOT_CREATE_LIVE
)

// SnapshotDeleteFlag defines how a snapshot should be deleted.
type SnapshotDeleteFlag uint32

// Possible values for SnapshotDeleteFlag.
const (
	SnapDeleteDefault      SnapshotDeleteFlag = 0
	SnapDeleteChildren     SnapshotDeleteFlag = C.VIR_DOMAIN_SNAPSHOT_DELETE_CHILDREN
	SnapDeleteMetadataOnly SnapshotDeleteFlag = C.VIR_DOMAIN_SNAPSHOT_DELETE_METADATA_ONLY
	SnapDeleteChildrenOnly SnapshotDeleteFlag = C.VIR_DOMAIN_SNAPSHOT_DELETE_CHILDREN_ONLY
)

// SnapshotRevertFlag defines how a snapshot revert operation should be performed.
type SnapshotRevertFlag uint32

// Possible values for SnapshotRevertFlag.
const (
	SnapRevertDefault SnapshotRevertFlag = 0
	SnapRevertRunning SnapshotRevertFlag = C.VIR_DOMAIN_SNAPSHOT_REVERT_RUNNING
	SnapRevertPaused  SnapshotRevertFlag = C.VIR_DOMAIN_SNAPSHOT_REVERT_PAUSED
	SnapRevertForce   SnapshotRevertFlag = C.VIR_DOMAIN_SNAPSHOT_REVERT_FORCE
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

// Delete deletes the snapshot.
func (snap Snapshot) Delete(flags SnapshotDeleteFlag) error {
	snap.log.Printf("deleting snapshot (flags = %v)...\n", flags)
	cRet := C.virDomainSnapshotDelete(snap.virSnapshot, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		snap.log.Printf("an error occurred: %v\n", err)
		return err
	}

	snap.log.Println("snapshot deleted")

	return nil
}

// Name gets the public name for that snapshot.
func (snap Snapshot) Name() string {
	snap.log.Println("reading snapshot name...")
	cName := C.virDomainSnapshotGetName(snap.virSnapshot)

	name := C.GoString(cName)
	snap.log.Printf("snapshot name: %v\n", name)

	return name
}

// Parent gets the parent snapshot for "snap", if any.
func (snap Snapshot) Parent() (Snapshot, error) {
	snap.log.Println("reading snapshot parent...")
	cParent := C.virDomainSnapshotGetParent(snap.virSnapshot, 0)
	if cParent == nil {
		err := LastError()
		snap.log.Printf("an error occurred: %v\n", err)
		return Snapshot{}, err
	}

	parent := Snapshot{
		log:         snap.log,
		virSnapshot: cParent,
	}

	snap.log.Println("parent obtained")

	return parent, nil
}

// XML provides an XML description of the domain snapshot.
func (snap Snapshot) XML(flags DomainXMLFlag) (string, error) {
	snap.log.Printf("reading snapshot XML (flags = %v)...\n", flags)
	cXML := C.virDomainSnapshotGetXMLDesc(snap.virSnapshot, C.uint(flags))
	if cXML == nil {
		err := LastError()
		snap.log.Printf("an error occurred: %v\n", err)
		return "", err
	}
	defer C.free(unsafe.Pointer(cXML))

	xml := C.GoString(cXML)
	snap.log.Printf("XML length: %v runes\n", utf8.RuneCountInString(xml))

	return xml, nil
}

// HasMetadata determines if the given snapshot is associated with libvirt
// metadata that would prevent the deletion of the domain.
func (snap Snapshot) HasMetadata() bool {
	snap.log.Println("checking whether snapshot has metadata...")
	cRet := C.virDomainSnapshotHasMetadata(snap.virSnapshot, 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		snap.log.Printf("an error occurred: %v\n", err)
		return false
	}

	metadata := (ret == 1)

	if metadata {
		snap.log.Println("snapshot has metadata")
	} else {
		snap.log.Println("snapshot doesn't have metadata")
	}

	return metadata
}

// IsCurrent determines if the given snapshot is the domain's current snapshot.
// See also "<Domain>.HasCurrentSnapshot".
func (snap Snapshot) IsCurrent() bool {
	snap.log.Println("checking whether snapshot is current...")
	cRet := C.virDomainSnapshotIsCurrent(snap.virSnapshot, 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		snap.log.Printf("an error occurred: %v\n", err)
		return false
	}

	current := (ret == 1)

	if current {
		snap.log.Println("snapshot is current")
	} else {
		snap.log.Println("snapshot isn't current")
	}

	return current
}

// Ref increments the reference count on the snapshot. For each additional call
// to this method, there shall be a corresponding call to "<Snapshot>.Free" to
// release the reference count, once the caller no longer needs the reference to this object.
func (snap Snapshot) Ref() error {
	snap.log.Println("incrementing snapshot's reference count...")
	cRet := C.virDomainSnapshotRef(snap.virSnapshot)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		snap.log.Printf("an error occurred: %v\n", err)
		return err
	}

	snap.log.Println("reference count incremented")

	return nil
}

// ListChildren collects the list of domain snapshots that are children of the
// given snapshot, and allocate an array to store those objects.
// By default, this command covers only direct children; it is also possible to
// expand things to cover all descendants, when "flags" includes
// SnapshotListDescendants. Also, some filters are provided in groups, where
// each group contains bits that describe mutually exclusive attributes of a
// snapshot, and where all bits within a group describe all possible snapshots.
// Some hypervisors might reject explicit bits from a group where the hypervisor
// cannot make a distinction. For a group supported by a given hypervisor, the
// behavior when no bits of a group are set is identical to the behavior when
// all bits in that group are set. When setting bits from more than one group,
// it is possible to select an impossible combination, in that case a hypervisor
// may return either 0 or an error.
func (snap Snapshot) ListChildren(flags SnapshotListFlag) ([]Snapshot, error) {
	var cSnaps []C.virDomainSnapshotPtr
	snapsSH := (*reflect.SliceHeader)(unsafe.Pointer(&cSnaps))

	snap.log.Printf("reading snapshot children (flags = %v)...\n", flags)
	cRet := C.virDomainSnapshotListAllChildren(snap.virSnapshot, (**C.virDomainSnapshotPtr)(unsafe.Pointer(&snapsSH.Data)), C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		snap.log.Printf("an error occurred: %v\n", err)
		return nil, err
	}
	defer C.free(unsafe.Pointer(snapsSH.Data))

	snapsSH.Cap = int(ret)
	snapsSH.Len = int(ret)

	snaps := make([]Snapshot, ret)

	for i := range snaps {
		snaps[i] = Snapshot{
			log:         snap.log,
			virSnapshot: cSnaps[i],
		}
	}

	snap.log.Printf("snapshots count: %v\n", ret)

	return snaps, nil
}

// Revert reverts the domain to a given snapshot.
// Normally, the domain will revert to the same state the domain was in while
// the snapshot was taken (whether inactive, running, or paused), except that
// disk snapshots default to reverting to inactive state. Including
// SnapRevertRunning in "flags" overrides the snapshot state to guarantee a
// running domain after the revert; or including SnapRevertPaused in "flags"
// guarantees a paused domain after the revert. These two flags are mutually
// exclusive. While a persistent domain does not need either flag, it is not
// possible to revert a transient domain into an inactive state, so transient
// domains require the use of one of these two flags.
// Reverting to any snapshot discards all configuration changes made since the
// last snapshot. Additionally, reverting to a snapshot from a running domain is
// a form of data loss, since it discards whatever is in the guest's RAM at the
// time. Since the very nature of keeping snapshots implies the intent to roll
// back state, no additional confirmation is normally required for these
// lossy effects.
// However, there are two particular situations where reverting will be refused
// by default, and where "flags" must include SnapRevertForce to acknowledge the
// risks. 1) Any attempt to revert to a snapshot that lacks the metadata to
// perform ABI compatibility checks (generally the case for snapshots that lack
// a full <domain> when listed by "<Snapshot>.XML()", such as those created
// prior to libvirt 0.9.5). 2) Any attempt to revert a running domain to an
// active state that requires starting a new hypervisor instance rather than
// reusing the existing hypervisor (since this would terminate all connections
// to the domain, such as such as VNC or Spice graphics) - this condition arises
// from active snapshots that are provably ABI incomaptible, as well as from
// inactive snapshots with a "flags" request to start the domain after
// the revert.
func (snap Snapshot) Revert(flags SnapshotRevertFlag) error {
	snap.log.Printf("reverting to snapshot (flags = %v)...\n", flags)
	cRet := C.virDomainRevertToSnapshot(snap.virSnapshot, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		snap.log.Printf("an error occurred: %v\n", err)
		return err
	}

	snap.log.Println("domain reverted")

	return nil
}
