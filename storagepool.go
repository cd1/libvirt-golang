package libvirt

// #include <libvirt/libvirt.h>
import "C"
import (
	"log"
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
