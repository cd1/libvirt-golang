package libvirt

// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"
import (
	"log"
)

type DomainFlag uint

const (
	DomActive DomainFlag = (1 << iota)
	DomInactive
	DomPersistent
	DomTransient
	DomRunning
	DomPaused
	DomShutOff
	DomOther
	DomManagedSave
	DomNoManagedSave
	DomAutostart
	DomNoAutostart
	DomHasSnapshot
	DomNoSnapshot
	DomAll = 0
)

type Domain struct {
	virDomain C.virDomainPtr
}

// Free frees the domain object. The running instance is kept alive. The data
// structure is freed and should not be used thereafter.
func (dom Domain) Free() *Error {
	cRet := C.virDomainFree(dom.virDomain)
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}

// Autostart provides a boolean value indicating whether the domain configured
// to be automatically started when the host machine boots.
func (dom Domain) Autostart() bool {
	var cAutostart C.int
	cRet := C.virDomainGetAutostart(dom.virDomain, &cAutostart)
	ret := int(cRet)

	if ret == -1 {
		if err := LastError(); err != nil {
			log.Println(err)
		}
		return false
	}

	autostart := int(cAutostart)
	return (autostart == 1)
}

// HasCurrentSnapshot determines if the domain has a current snapshot.
func (dom Domain) HasCurrentSnapshot() bool {
	cRet := C.virDomainHasCurrentSnapshot(dom.virDomain, 0)
	ret := int(cRet)

	if ret == -1 {
		if err := LastError(); err != nil {
			log.Println(err)
		}
		return false
	}

	return (ret == 1)
}

// HasManagedSaveImage checks if a domain has a managed save image as created
// by ManagedSave(). Note that any running domain should not have such an
// image, as it should have been removed on restart.
func (dom Domain) HasManagedSaveImage() bool {
	cRet := C.virDomainHasManagedSaveImage(dom.virDomain, 0)
	ret := int(cRet)

	if ret == -1 {
		if err := LastError(); err != nil {
			log.Println(err)
		}
		return false
	}

	return (ret == 1)
}

// IsActive determines if the domain is currently running.
func (dom Domain) IsActive() bool {
	cRet := C.virDomainIsActive(dom.virDomain)
	ret := int(cRet)

	if ret == -1 {
		if err := LastError(); err != nil {
			log.Println(err)
		}
		return false
	}

	return (ret == 1)
}

// IsPersistent determines if the domain has a persistent configuration which
// means it will still exist after shutting down
func (dom Domain) IsPersistent() bool {
	cRet := C.virDomainIsPersistent(dom.virDomain)
	ret := int(cRet)

	if ret == -1 {
		if err := LastError(); err != nil {
			log.Println(err)
		}
		return false
	}

	return (ret == 1)
}

// IsUpdated determines if the domain has been updated.
func (dom Domain) IsUpdated() bool {
	cRet := C.virDomainIsUpdated(dom.virDomain)
	ret := int(cRet)

	if ret == -1 {
		if err := LastError(); err != nil {
			log.Println(err)
		}
		return false
	}

	return (ret == 1)
}
