package libvirt

// #include <libvirt/libvirt.h>
import "C"
import (
	"log"
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
