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
