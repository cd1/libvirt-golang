package libvirt

// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"

type DomainFlag uint

const (
	ALL_DOMAINS DomainFlag = (0 << iota)
	ACTIVE
	INACTIVE
	PERSISTENT
	TRANSIENT
	RUNNING
	PAUSED
	SHUTOFF
	OTHER
	MANAGEDSAVE
	NO_MANAGEDSAVE
	AUTOSTART
	NO_AUTOSTART
	HAS_SNAPSHOT
	NO_SNAPSHOT
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
	} else {
		return nil
	}
}
