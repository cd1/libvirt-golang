package libvirt

// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"

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
