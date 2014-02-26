package libvirt

// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"
import (
	"log"
	"unsafe"
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

type DomainMetadataType uint

const (
	DomMetaDescription DomainMetadataType = iota
	DomMetaTitle
	DomMetaElement
)

type DomainModificationImpact uint

const (
	DomAffectCurrent = iota
	DomAffectLive
	DomAffectConfig
)

type DomainXMLFlag uint

const (
	DomXMLSecure DomainXMLFlag = (1 << iota)
	DomXMLInactive
	DomXMLUpdateCPU
	DomXMLMigratable
	DomXMLDefault = 0
)

type DomainCreateFlag uint

const (
	DomCreateStartPaused DomainCreateFlag = (1 << iota)
	DomCreateStartAutodestroy
	DomCreateStartBypassCache
	DomCreateStartForceBoot
	DomCreateDefault = 0
)

type DomainDestroyFlag uint

const (
	DomDestroyGraceful DomainDestroyFlag = (1 << iota)
	DomDestroyDefault                    = 0
)

type DomainUndefineFlag uint

const (
	DomUndefineManagedSave DomainUndefineFlag = (1 << iota)
	DomUndefineSnapshotsMetadata
	DomUndefineDefault = 0
)

type DomainRebootFlag uint

const (
	DomRebootACPIPowerBtn DomainRebootFlag = (1 << iota)
	DomRebootGuestAgent
	DomRebootInitctl
	DomRebootSignal
	DomRebootDefault = 0
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

// OSType gets the type of domain operation system.
func (dom Domain) OSType() (string, *Error) {
	os := C.virDomainGetOSType(dom.virDomain)
	if os == nil {
		return "", LastError()
	}
	defer C.free(unsafe.Pointer(os))

	return C.GoString(os), nil
}

// Name gets the public name for that domain.
func (dom Domain) Name() string {
	cName := C.virDomainGetName(dom.virDomain)
	return C.GoString(cName)
}

// Hostname gets the hostname for that domain.
func (dom Domain) Hostname() (string, *Error) {
	cHostname := C.virDomainGetHostname(dom.virDomain, 0)
	if cHostname == nil {
		return "", LastError()
	}
	defer C.free(unsafe.Pointer(cHostname))

	return C.GoString(cHostname), nil
}

// ID gets the hypervisor ID number for the domain.
func (dom Domain) ID() (uint, *Error) {
	cID := C.virDomainGetID(dom.virDomain)
	id := uint(cID)

	if id == ^uint(0) { // Go: ^uint(0) == C: (unsigned int) -1
		return 0, LastError()
	}

	return id, nil
}

// UUID gets the UUID for a domain as string. For more information about UUID
// see RFC4122.
func (dom Domain) UUID() (string, *Error) {
	cUUID := (*C.char)(C.malloc(C.size_t(C.VIR_UUID_STRING_BUFLEN)))
	defer C.free(unsafe.Pointer(cUUID))

	cRet := C.virDomainGetUUIDString(dom.virDomain, cUUID)
	ret := int(cRet)

	if ret == -1 {
		return "", LastError()
	}

	return C.GoString(cUUID), nil
}

// XML provides an XML description of the domain. The description may be reused
// later to relaunch the domain with CreateXML().
func (dom Domain) XML(typ DomainXMLFlag) (string, *Error) {
	cXML := C.virDomainGetXMLDesc(dom.virDomain, C.uint(typ))
	if cXML == nil {
		return "", LastError()
	}
	defer C.free(unsafe.Pointer(cXML))

	return C.GoString(cXML), nil
}

// Metadata retrieves the appropriate domain element given by "type".
func (dom Domain) Metadata(typ DomainMetadataType, xmlns string, impact DomainModificationImpact) (string, *Error) {
	cXMLNS := C.CString(xmlns)
	defer C.free(unsafe.Pointer(cXMLNS))

	cMetadata := C.virDomainGetMetadata(dom.virDomain, C.int(typ), cXMLNS, C.uint(impact))
	if cMetadata == nil {
		return "", LastError()
	}
	defer C.free(unsafe.Pointer(cMetadata))

	return C.GoString(cMetadata), nil
}

// Destroy destroys the domain object. The running instance is shutdown if not
// down already and all resources used by it are given back to the hypervisor.
// This does not free the associated virDomainPtr object. This function may
// require privileged access.
func (dom Domain) Destroy(flags DomainDestroyFlag) *Error {
	cRet := C.virDomainDestroyFlags(dom.virDomain, C.uint(flags))
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}

// Create launches a defined domain. If the call succeeds the domain moves from
// the defined to the running domains pools.
func (dom Domain) Create(flags DomainCreateFlag) *Error {
	cRet := C.virDomainCreateWithFlags(dom.virDomain, C.uint(flags))
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}

// Undefine undefines a domain. If the domain is running, it's converted to
// transient domain, without stopping it. If the domain is inactive, the domain
// configuration is removed.
func (dom Domain) Undefine(flags DomainUndefineFlag) *Error {
	cRet := C.virDomainUndefineFlags(dom.virDomain, C.uint(flags))
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}

// Reboot reboots a domain, the domain object is still usable thereafter, but
// the domain OS is being stopped for a restart. Note that the guest OS may
// ignore the request. Additionally, the hypervisor may check and support the
// domain 'on_reboot' XML setting resulting in a domain that shuts down instead
// of rebooting.
func (dom Domain) Reboot(flags DomainRebootFlag) *Error {
	cRet := C.virDomainReboot(dom.virDomain, C.uint(flags))
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}

// Reset resets a domain immediately without any guest OS shutdown. Reset
// emulates the power reset button on a machine, where all hardware sees the
// RST line set and reinitializes internal state.
// Note that there is a risk of data loss caused by reset without any guest
// OS shutdown.
func (dom Domain) Reset() *Error {
	cRet := C.virDomainReset(dom.virDomain, 0)
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}

// Shutdown shuts down a domain, the domain object is still usable thereafter,
// but the domain OS is being stopped. Note that the guest OS may ignore the
// request. Additionally, the hypervisor may check and support the domain
// 'on_poweroff' XML setting resulting in a domain that reboots instead of
// shutting down. For guests that react to a shutdown request, the differences
// from Destroy() are that the guests disk storage will be in a stable state
// rather than having the (virtual) power cord pulled, and this command returns
// as soon as the shutdown request is issued rather than blocking until the
// guest is no longer running.
func (dom Domain) Shutdown() *Error {
	cRet := C.virDomainShutdown(dom.virDomain)
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}
