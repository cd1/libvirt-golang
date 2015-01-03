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

// StoragePoolListFlag defines a filter when listing storage pools.
type StoragePoolListFlag uint32

// Possible values for StoragePoolListFlag.
const (
	PoolListAll         StoragePoolListFlag = 0
	PoolListInactive    StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_INACTIVE
	PoolListActive      StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_ACTIVE
	PoolListPersistent  StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_PERSISTENT
	PoolListTransient   StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_TRANSIENT
	PoolListAutostart   StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_AUTOSTART
	PoolListNoAutostart StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_NO_AUTOSTART
	PoolListDir         StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_DIR
	PoolListFS          StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_FS
	PoolListNetFS       StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_NETFS
	PoolListLogical     StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_LOGICAL
	PoolListDisk        StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_DISK
	PoolListISCSI       StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_ISCSI
	PoolListSCSI        StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_SCSI
	PoolListMPath       StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_MPATH
	PoolListRBD         StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_RBD
	PoolListSheepdog    StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_SHEEPDOG
	PoolListGluster     StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_GLUSTER
	PoolListZFS         StoragePoolListFlag = C.VIR_CONNECT_LIST_STORAGE_POOLS_ZFS
)

// StoragePoolDeleteFlag defines how a storage pool should be deleted.
type StoragePoolDeleteFlag uint32

// Possible values for StoragePoolDeleteFlag.
const (
	PoolDeleteNormal StoragePoolDeleteFlag = C.VIR_STORAGE_POOL_DELETE_NORMAL
	PoolDeleteZeroed StoragePoolDeleteFlag = C.VIR_STORAGE_POOL_DELETE_ZEROED
)

// StorageXMLFlag defines how the XML content should be read from a storage resource.
type StorageXMLFlag uint32

// Possible values for StorageXMLFlag.
const (
	StorageXMLDefault  StorageXMLFlag = 0
	StorageXMLInactive StorageXMLFlag = C.VIR_STORAGE_XML_INACTIVE
)

// StoragePoolState represents the state of a storage pool.
type StoragePoolState uint32

// Possible values for StoragePoolState.
const (
	PoolStateInactive     StoragePoolState = C.VIR_STORAGE_POOL_INACTIVE
	PoolStateBuilding     StoragePoolState = C.VIR_STORAGE_POOL_BUILDING
	PoolStateRunning      StoragePoolState = C.VIR_STORAGE_POOL_RUNNING
	PoolStateDegraded     StoragePoolState = C.VIR_STORAGE_POOL_DEGRADED
	PoolStateInaccessible StoragePoolState = C.VIR_STORAGE_POOL_INACCESSIBLE
)

// StoragePoolBuildFlag defines how a storage pool should be built.
type StoragePoolBuildFlag uint32

// Possible values for StoragePoolBuildFlag.
const (
	PoolBuildNew         StoragePoolBuildFlag = C.VIR_STORAGE_POOL_BUILD_NEW
	PoolBuildRepair      StoragePoolBuildFlag = C.VIR_STORAGE_POOL_BUILD_REPAIR
	PoolBuildResize      StoragePoolBuildFlag = C.VIR_STORAGE_POOL_BUILD_RESIZE
	PoolBuildNoOverwrite StoragePoolBuildFlag = C.VIR_STORAGE_POOL_BUILD_NO_OVERWRITE
	PoolBuildOverwrite   StoragePoolBuildFlag = C.VIR_STORAGE_POOL_BUILD_OVERWRITE
)

// StoragePool holds a libvirt storage pool. There are no exported fields.
type StoragePool struct {
	log            *log.Logger
	virStoragePool C.virStoragePoolPtr
}

// Free frees a storage pool object, releasing all memory associated with it.
// Does not change the state of the pool on the host.
func (pool StoragePool) Free() error {
	pool.log.Println("freeing storage pool object...")
	cRet := C.virStoragePoolFree(pool.virStoragePool)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return err
	}

	pool.log.Println("pool freed")

	return nil
}

// Undefine undefines an inactive storage pool.
func (pool StoragePool) Undefine() error {
	pool.log.Println("undefining storage pool...")
	cRet := C.virStoragePoolUndefine(pool.virStoragePool)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return err
	}

	pool.log.Println("pool undefined")

	return nil
}

// Create starts an inactive storage pool.
func (pool StoragePool) Create() error {
	pool.log.Println("creating storage pool...")
	cRet := C.virStoragePoolCreate(pool.virStoragePool, 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return err
	}

	pool.log.Println("pool created")

	return nil
}

// Destroy destroys an active storage pool. This will deactivate the pool on the
// host, but keep any persistent config associated with it. If it has a
// persistent config it can later be restarted with "Create". This does not free
// the associated StoragePool object.
func (pool StoragePool) Destroy() error {
	pool.log.Println("destroying storage pool...")
	cRet := C.virStoragePoolDestroy(pool.virStoragePool)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return err
	}

	pool.log.Println("pool destroyed")

	return nil
}

// Delete deletes the underlying pool resources. This is a non-recoverable
// operation. The StoragePool object itself is not free'd.
func (pool StoragePool) Delete(flags StoragePoolDeleteFlag) error {
	pool.log.Printf("deleting storage pool (flags = %v)...\n", flags)
	cRet := C.virStoragePoolDelete(pool.virStoragePool, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return err
	}

	pool.log.Println("pool deleted")

	return nil
}

// IsActive determines if the storage pool is currently running.
func (pool StoragePool) IsActive() (bool, error) {
	pool.log.Println("checking whether storage pool is active...")
	cRet := C.virStoragePoolIsActive(pool.virStoragePool)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return false, err
	}

	active := (ret == 1)

	if active {
		pool.log.Println("pool is active")
	} else {
		pool.log.Println("pool is not active")
	}

	return active, nil
}

// IsPersistent determines if the storage pool has a persistent configuration
// which means it will still exist after shutting down.
func (pool StoragePool) IsPersistent() (bool, error) {
	pool.log.Println("checking whether storage pool is persistent...")
	cRet := C.virStoragePoolIsPersistent(pool.virStoragePool)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return false, err
	}

	persistent := (ret == 1)

	if persistent {
		pool.log.Println("pool is persistent")
	} else {
		pool.log.Println("pool is not persistent")
	}

	return persistent, nil
}

// Name fetches the locally unique name of the storage pool.
func (pool StoragePool) Name() (string, error) {
	pool.log.Println("reading storage pool name...")
	cName := C.virStoragePoolGetName(pool.virStoragePool)

	if cName == nil {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return "", err
	}

	name := C.GoString(cName)
	pool.log.Printf("name: %v\n", name)

	return name, nil
}

// UUID fetches the globally unique ID of the storage pool as a string.
func (pool StoragePool) UUID() (string, error) {
	cUUID := (*C.char)(C.malloc(C.size_t(C.VIR_UUID_STRING_BUFLEN)))
	defer C.free(unsafe.Pointer(cUUID))

	pool.log.Println("reading storage pool UUID...")
	cRet := C.virStoragePoolGetUUIDString(pool.virStoragePool, cUUID)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return "", err
	}

	uuid := C.GoString(cUUID)
	pool.log.Printf("UUID: %v\n", uuid)

	return uuid, nil
}

// XML fetches an XML document describing all aspects of the storage pool. This
// is suitable for later feeding back into the
// "<Connection>.CreateStoragePool" method.
func (pool StoragePool) XML(flags StorageXMLFlag) (string, error) {
	pool.log.Printf("reading storage pool XML (flags = %v)...\n", flags)
	cXML := C.virStoragePoolGetXMLDesc(pool.virStoragePool, C.uint(flags))

	if cXML == nil {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return "", err
	}
	defer C.free(unsafe.Pointer(cXML))

	xml := C.GoString(cXML)
	pool.log.Printf("XML length: %v runes\n", utf8.RuneCountInString(xml))

	return xml, nil
}

// InfoState extracts the storage pool state.
func (pool StoragePool) InfoState() (StoragePoolState, error) {
	var cInfo C.virStoragePoolInfo

	pool.log.Println("reading storage pool state...")
	cRet := C.virStoragePoolGetInfo(pool.virStoragePool, &cInfo)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	state := StoragePoolState(cInfo.state)
	pool.log.Printf("state: %v\n", state)

	return state, nil
}

// InfoCapacity extracts the storage pool logical size (bytes).
func (pool StoragePool) InfoCapacity() (uint64, error) {
	var cInfo C.virStoragePoolInfo

	pool.log.Println("reading storage pool capacity...")
	cRet := C.virStoragePoolGetInfo(pool.virStoragePool, &cInfo)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	capacity := uint64(cInfo.capacity)
	pool.log.Printf("capacity: %v bytes\n", capacity)

	return capacity, nil
}

// InfoAllocation extracts the storage pool current allocation (bytes).
func (pool StoragePool) InfoAllocation() (uint64, error) {
	var cInfo C.virStoragePoolInfo

	pool.log.Println("reading storage pool allocation...")
	cRet := C.virStoragePoolGetInfo(pool.virStoragePool, &cInfo)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	allocation := uint64(cInfo.allocation)
	pool.log.Printf("allocation: %v bytes\n", allocation)

	return allocation, nil
}

// InfoAvailable extracts the storage pool remaining free space (bytes)
func (pool StoragePool) InfoAvailable() (uint64, error) {
	var cInfo C.virStoragePoolInfo

	pool.log.Println("reading storage pool available space...")
	cRet := C.virStoragePoolGetInfo(pool.virStoragePool, &cInfo)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	available := uint64(cInfo.available)
	pool.log.Printf("available: %v bytes\n", available)

	return available, nil
}

// Autostart fetches the value of the autostart flag, which determines whether
// the pool is automatically started at boot time.
func (pool StoragePool) Autostart() (bool, error) {
	var cAutostart C.int

	pool.log.Println("checking whether storage pool autostarts...")
	cRet := C.virStoragePoolGetAutostart(pool.virStoragePool, &cAutostart)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return false, err
	}

	autostart := (int32(cAutostart) == 1)

	if autostart {
		pool.log.Println("pool autostarts")
	} else {
		pool.log.Println("pool does not autostart")
	}

	return autostart, nil
}

// SetAutostart sets the autostart flag.
func (pool StoragePool) SetAutostart(autostart bool) error {
	var autostartInt int32
	if autostart {
		pool.log.Println("enabling storage pool autostart...")
		autostartInt = 1
	} else {
		pool.log.Println("disabling storage pool autostart...")
		autostartInt = 0
	}

	cRet := C.virStoragePoolSetAutostart(pool.virStoragePool, C.int(autostartInt))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return err
	}

	if autostart {
		pool.log.Printf("autostart enabled")
	} else {
		pool.log.Printf("autostart disabled")
	}

	return nil
}

// Build builds the underlying storage pool.
func (pool StoragePool) Build(flags StoragePoolBuildFlag) error {
	pool.log.Printf("building storage pool (flags = %v)...\n", flags)
	cRet := C.virStoragePoolBuild(pool.virStoragePool, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return err
	}

	pool.log.Println("pool built")

	return nil
}

// Refresh requests that the pool refresh its list of volumes. This may involve
// communicating with a remote server, and/or initializing new devices at the
// OS layer.
func (pool StoragePool) Refresh() error {
	pool.log.Println("refreshing storage pool...")
	cRet := C.virStoragePoolRefresh(pool.virStoragePool, 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return err
	}

	pool.log.Println("pool refreshed")

	return nil
}

// Ref increments the reference count on the pool. For each additional call to
// this method, there shall be a corresponding call to "Free" to release the
// reference count, once the caller no longer needs the reference to
// this object.
// This method is typically useful for applications where multiple threads are
// using a connection, and it is required that the connection remain open until
// all threads have finished using it. ie, each new thread using a pool would
// increment the reference count.
func (pool StoragePool) Ref() error {
	pool.log.Println("incrementing storage pool's reference count...")
	cRet := C.virStoragePoolRef(pool.virStoragePool)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return err
	}

	return nil
}

// ListStorageVolumes collects the list of storage volumes, and allocate an
// array to store those objects.
func (pool StoragePool) ListStorageVolumes() ([]StorageVolume, error) {
	var cStorageVolumes []C.virStorageVolPtr
	cStorageVolumesSH := (*reflect.SliceHeader)(unsafe.Pointer(&cStorageVolumes))

	pool.log.Println("reading storage volumes...")
	cRet := C.virStoragePoolListAllVolumes(pool.virStoragePool, (**C.virStorageVolPtr)(unsafe.Pointer(&cStorageVolumesSH.Data)), 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		pool.log.Printf("an error occurred: %v\n", err)
		return nil, err
	}
	defer C.free(unsafe.Pointer(cStorageVolumesSH.Data))

	cStorageVolumesSH.Cap = int(ret)
	cStorageVolumesSH.Len = int(ret)

	storageVolumes := make([]StorageVolume, ret)
	for i, cVol := range cStorageVolumes {
		storageVolumes[i] = StorageVolume{
			log:           pool.log,
			virStorageVol: cVol,
		}
	}

	pool.log.Printf("volumes count: %v\n", ret)

	return storageVolumes, nil
}
