package libvirt

// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"
import (
	"log"
	"unicode/utf8"
	"unsafe"
)

// StorageVolumeType represents the type of a storage volume.
type StorageVolumeType uint32

// Possible values for StorageVolumeType.
const (
	VolTypeFile    StorageVolumeType = C.VIR_STORAGE_VOL_FILE
	VolTypeBlock   StorageVolumeType = C.VIR_STORAGE_VOL_BLOCK
	VolTypeDir     StorageVolumeType = C.VIR_STORAGE_VOL_DIR
	VolTypeNetwork StorageVolumeType = C.VIR_STORAGE_VOL_NETWORK
	VolTypeNetdir  StorageVolumeType = C.VIR_STORAGE_VOL_NETDIR
)

// StorageVolumeResizeFlag defines how a storage volume should be resized.
type StorageVolumeResizeFlag uint32

// Possible values for StorageVolumeResizeFlag.
const (
	VolResizeDefault  StorageVolumeResizeFlag = 0
	VolResizeAllocate StorageVolumeResizeFlag = C.VIR_STORAGE_VOL_RESIZE_ALLOCATE
	VolResizeDelta    StorageVolumeResizeFlag = C.VIR_STORAGE_VOL_RESIZE_DELTA
	VolResizeShrink   StorageVolumeResizeFlag = C.VIR_STORAGE_VOL_RESIZE_SHRINK
)

// StorageVolumeWipeAlgorithm defines the algorithm used to wipe a
// storage volume.
type StorageVolumeWipeAlgorithm uint32

// Possible values for StorageVolumeWipeAlgorithm.
const (
	VolWipeAlgZero       StorageVolumeWipeAlgorithm = C.VIR_STORAGE_VOL_WIPE_ALG_ZERO
	VolWipeAlgNNSA       StorageVolumeWipeAlgorithm = C.VIR_STORAGE_VOL_WIPE_ALG_NNSA
	VolWipeAlgDoD        StorageVolumeWipeAlgorithm = C.VIR_STORAGE_VOL_WIPE_ALG_DOD
	VolWipeAlgBSI        StorageVolumeWipeAlgorithm = C.VIR_STORAGE_VOL_WIPE_ALG_BSI
	VolWipeAlgGutmann    StorageVolumeWipeAlgorithm = C.VIR_STORAGE_VOL_WIPE_ALG_GUTMANN
	VolWipeAlgSchneier   StorageVolumeWipeAlgorithm = C.VIR_STORAGE_VOL_WIPE_ALG_SCHNEIER
	VolWipeAlgPfitzner7  StorageVolumeWipeAlgorithm = C.VIR_STORAGE_VOL_WIPE_ALG_PFITZNER7
	VolWipeAlgPfitzner33 StorageVolumeWipeAlgorithm = C.VIR_STORAGE_VOL_WIPE_ALG_PFITZNER33
	VolWipeAlgRandom     StorageVolumeWipeAlgorithm = C.VIR_STORAGE_VOL_WIPE_ALG_RANDOM
)

// StorageVolume holds a libvirt storage volume. There are no exported fields.
type StorageVolume struct {
	log           *log.Logger
	virStorageVol C.virStorageVolPtr
}

// Free releases the storage volume handle. The underlying storage volume
// continues to exist.
func (vol StorageVolume) Free() error {
	vol.log.Println("freeing storage volume object...")
	cRet := C.virStorageVolFree(vol.virStorageVol)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		vol.log.Printf("an error occurred: %v\n", err)
		return err
	}

	vol.log.Println("volume freed")

	return nil
}

// Delete deletes the storage volume from the pool.
func (vol StorageVolume) Delete() error {
	vol.log.Println("deleting storage volume...")
	cRet := C.virStorageVolDelete(vol.virStorageVol, 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		vol.log.Printf("an error occurred: %v\n", err)
		return err
	}

	vol.log.Println("volume deleted")

	return nil
}

// Key fetches the storage volume key. This is globally unique, so the same
// volume will have the same key no matter what host it is accessed from.
func (vol StorageVolume) Key() (string, error) {
	vol.log.Println("reading storage volume key...")
	cKey := C.virStorageVolGetKey(vol.virStorageVol)

	if cKey == nil {
		err := LastError()
		vol.log.Printf("an error occurred: %v\n", err)
		return "", err
	}

	key := C.GoString(cKey)
	vol.log.Printf("key: %v\n", key)

	return key, nil
}

// Name fetches the storage volume name. This is unique within the scope of
// a pool.
func (vol StorageVolume) Name() (string, error) {
	vol.log.Println("reading storage volume name...")
	cName := C.virStorageVolGetName(vol.virStorageVol)

	if cName == nil {
		err := LastError()
		vol.log.Printf("an error occurred: %v\n", err)
		return "", err
	}

	name := C.GoString(cName)
	vol.log.Printf("name: %v\n", name)

	return name, nil
}

// Path fetches the storage volume path. Depending on the pool configuration
// this is either persistent across hosts, or dynamically assigned at pool
// startup. Consult pool documentation for information on getting the
// persistent naming.
func (vol StorageVolume) Path() (string, error) {
	vol.log.Println("reading storage volume path...")
	cPath := C.virStorageVolGetPath(vol.virStorageVol)

	if cPath == nil {
		err := LastError()
		vol.log.Printf("an error occurred: %v\n", err)
		return "", err
	}
	defer C.free(unsafe.Pointer(cPath))

	path := C.GoString(cPath)
	vol.log.Printf("path: %v\n", path)

	return path, nil
}

// XML fetches an XML document describing all aspects of the storage volume.
func (vol StorageVolume) XML() (string, error) {
	vol.log.Println("reading storage volume XML...")
	cXML := C.virStorageVolGetXMLDesc(vol.virStorageVol, 0)

	if cXML == nil {
		err := LastError()
		vol.log.Printf("an error occurred: %v\n", err)
		return "", err
	}

	xml := C.GoString(cXML)
	vol.log.Printf("XML length: %v runes\n", utf8.RuneCountInString(xml))

	return xml, nil
}

// InfoType fetches volatile information about the storage volume:
// current type.
func (vol StorageVolume) InfoType() (StorageVolumeType, error) {
	var cInfo C.virStorageVolInfo

	vol.log.Println("reading storage volume type...")
	cRet := C.virStorageVolGetInfo(vol.virStorageVol, &cInfo)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		vol.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	typ := StorageVolumeType(cInfo._type)
	vol.log.Printf("type: %v\n", typ)

	return typ, nil
}

// InfoCapacity fetches volatile information about the storage volume:
// current capacity.
func (vol StorageVolume) InfoCapacity() (uint64, error) {
	var cInfo C.virStorageVolInfo

	vol.log.Println("reading storage volume capacity...")
	cRet := C.virStorageVolGetInfo(vol.virStorageVol, &cInfo)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		vol.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	capacity := uint64(cInfo.capacity)
	vol.log.Printf("capacity: %v\n", capacity)

	return capacity, nil
}

// InfoAllocation fetches volatile information about the storage volume:
// current allocation.
func (vol StorageVolume) InfoAllocation() (uint64, error) {
	var cInfo C.virStorageVolInfo

	vol.log.Println("reading storage volume allocation...")
	cRet := C.virStorageVolGetInfo(vol.virStorageVol, &cInfo)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		vol.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	allocation := uint64(cInfo.allocation)
	vol.log.Printf("allocation: %v\n", allocation)

	return allocation, nil
}

// Resize changes the capacity of the storage volume to "capacity". The
// operation will fail if the new capacity requires allocation that would exceed
// the remaining free space in the parent pool. The contents of the new capacity
// will appear as all zero bytes. The capacity value will be rounded to the
// granularity supported by the hypervisor.
// Normally, the operation will attempt to affect capacity with a minimum impact
// on allocation (that is, the default operation favors a sparse resize). If
// "flags" contains VolResizeAllocate, then the operation will ensure that
// allocation is sufficient for the new capacity; this may make the operation
// take noticeably longer.
// Normally, the operation treats "capacity" as the new size in bytes; but if
// "flags" contains VolResizeDelta, then "capacity" represents the size
// difference to add to the current size. It is up to the storage pool
// implementation whether unaligned requests are rounded up to the next valid
// boundary, or rejected.
// Normally, this operation should only be used to enlarge capacity; but if
// "flags" contains VolResizeShrink, it is possible to attempt a reduction in
// capacity even though it might cause data loss. If VolResizeDelta is also
// present, then "capacity" is subtracted from the current size; without it,
// "capacity" represents the absolute new size regardless of whether it is
// larger or smaller than the current size.
func (vol StorageVolume) Resize(capacity uint64, flags StorageVolumeResizeFlag) error {
	vol.log.Printf("resizing storage volume to %v bytes (flags = %v)...\n", capacity, flags)
	cRet := C.virStorageVolResize(vol.virStorageVol, C.ulonglong(capacity), C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		vol.log.Printf("an error occurred: %v\n", err)
		return err
	}

	vol.log.Println("volume resized")

	return nil
}

// Wipe ensure data previously on a volume is not accessible to future reads.
func (vol StorageVolume) Wipe(alg StorageVolumeWipeAlgorithm) error {
	vol.log.Printf("wiping storage volume with algorithm %v...\n", alg)
	cRet := C.virStorageVolWipePattern(vol.virStorageVol, C.uint(alg), 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		vol.log.Printf("an error occurred: %v\n", err)
		return err
	}

	vol.log.Println("volume wiped")

	return nil
}

// Ref increments the reference count on the vol. For each additional call to
// this method, there shall be a corresponding call to "Free" to release the
// reference count, once the caller no longer needs the reference to
// this object.
// This method is typically useful for applications where multiple threads are
// using a connection, and it is required that the connection remain open until
// all threads have finished using it. ie, each new thread using a vol would
// increment the reference count.
func (vol StorageVolume) Ref() error {
	vol.log.Println("incrementing storage volume's reference count...")
	cRet := C.virStorageVolRef(vol.virStorageVol)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		vol.log.Printf("an error occurred: %v\n", err)
		return err
	}

	vol.log.Println("reference count incremented")

	return nil
}

// StoragePool fetches a storage pool which contains a particular volume.
// "Free" should be used to free the resources after the storage pool object is
// no longer needed.
func (vol StorageVolume) StoragePool() (StoragePool, error) {
	vol.log.Println("looking up storage pool by storage volume...")
	cPool := C.virStoragePoolLookupByVolume(vol.virStorageVol)

	if cPool == nil {
		err := LastError()
		vol.log.Printf("an error occurred: %v\n", err)
		return StoragePool{}, err
	}

	vol.log.Println("pool found")

	pool := StoragePool{
		log:            vol.log,
		virStoragePool: cPool,
	}

	return pool, nil
}
