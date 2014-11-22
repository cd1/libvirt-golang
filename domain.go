package libvirt

// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"
import (
	"errors"
	"log"
	"time"
	"unicode/utf8"
	"unsafe"
)

type DomainFlag uint32

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

type DomainMetadataType uint32

const (
	DomMetaDescription DomainMetadataType = iota
	DomMetaTitle
	DomMetaElement
)

type DomainModificationImpact uint32

const (
	DomAffectCurrent = iota
	DomAffectLive
	DomAffectConfig
)

type DomainXMLFlag uint32

const (
	DomXMLSecure DomainXMLFlag = (1 << iota)
	DomXMLInactive
	DomXMLUpdateCPU
	DomXMLMigratable
	DomXMLDefault = 0
)

type DomainCreateFlag uint32

const (
	DomCreateStartPaused DomainCreateFlag = (1 << iota)
	DomCreateStartAutodestroy
	DomCreateStartBypassCache
	DomCreateStartForceBoot
	DomCreateDefault = 0
)

type DomainDestroyFlag uint32

const (
	DomDestroyGraceful DomainDestroyFlag = (1 << iota)
	DomDestroyDefault                    = 0
)

type DomainUndefineFlag uint32

const (
	DomUndefineManagedSave DomainUndefineFlag = (1 << iota)
	DomUndefineSnapshotsMetadata
	DomUndefineDefault = 0
)

type DomainRebootFlag uint32

const (
	DomRebootACPIPowerBtn DomainRebootFlag = (1 << iota)
	DomRebootGuestAgent
	DomRebootInitctl
	DomRebootSignal
	DomRebootDefault = 0
)

type DomainState uint32

const (
	DomStateNone DomainState = iota
	DomStateRunning
	DomStateBlocked
	DomStatePaused
	DomStateShutdown
	DomStateShutoff
	DomStateCrashed
	DomStatePMSuspended
)

type DomainNostateReason uint32

const (
	DomNostateReasonUnknown DomainNostateReason = iota
)

type DomainRunningReason uint32

const (
	DomRunningReasonUnknown DomainRunningReason = iota
	DomRunningReasonBooted
	DomRunningReasonMigrated
	DomRunningReasonRestored
	DomRunningReasonFromSnapshot
	DomRunningReasonUnpaused
	DomRunningReasonMigrationCancelled
	DomRunningReasonSaveCancelled
	DomRunningReasonWakeUp
	DomRunningReasonCrashed
)

type DomainBlockedReason uint32

const (
	DomBlockedReasonUnkwown DomainBlockedReason = iota
)

type DomainPausedReason uint32

const (
	DomPausedReasonUnknown DomainPausedReason = iota
	DomPausedReasonUser
	DomPausedReasonMigration
	DomPausedReasonSave
	DomPausedReasonDump
	DomPausedReasonIOError
	DomPausedReasonWatchdog
	DomPausedReasonFromSnapshot
	DomPausedReasonShuttingDown
	DomPausedReasonSnapshot
	DomPausedReasonCrashed
)

type DomainShutdownReason uint32

const (
	DomShutdownReasonUnknown DomainShutdownReason = iota
	DomShutdownReasonUser
)

type DomainShutoffReason uint32

const (
	DomShutoffReasonUnknown DomainShutoffReason = iota
	DomShutoffReasonShutdown
	DomShutoffReasonDestroyed
	DomShutoffReasonCrashed
	DomShutoffReasonMigrated
	DomShutoffReasonSaved
	DomShutoffReasonFailed
	DomShutoffReasonFromSnapshot
)

type DomainCrashedReason uint32

const (
	DomCrashedReasonUnknown DomainCrashedReason = iota
	DomCrashedReasonPanicked
)

type DomainPMSuspendedReason uint32

const (
	DomPMSuspendedReasonUnknown DomainPMSuspendedReason = iota
)

type DomainDumpFlag uint32

const (
	DomDumpCrash DomainDumpFlag = (1 << iota)
	DomDumpLive
	DomDumpBypassCache
	DomDumpReset
	DomDumpMemoryOnly
	DomDumpDefault = 0
)

type DomainVCPUsFlag uint32

const (
	DomVCPusConfig  DomainVCPUsFlag = DomAffectConfig
	DomVCPUsCurrent                 = DomAffectCurrent
	DomVCPUsLive                    = DomAffectLive
	DomVCPUsMaximum                 = 4
	DomVCPUsGuest                   = 8
)

type DomainSaveFlag uint32

const (
	DomSaveBypassCache DomainSaveFlag = (1 << iota)
	DomSaveRunning
	DomSavePaused
	DomSaveDefault = 0
)

type DomainDeviceModifyFlag uint32

const (
	DomDeviceModifyConfig  DomainDeviceModifyFlag = DomAffectConfig
	DomDeviceModifyCurrent                        = DomAffectCurrent
	DomDeviceModifyLive                           = DomAffectLive
	DomDeviceModifyForce                          = 4
)

type DomainMemoryFlag uint32

const (
	DomMemoryConfig  DomainMemoryFlag = DomAffectConfig
	DomMemoryCurrent                  = DomAffectCurrent
	DomMemoryLive                     = DomAffectLive
	DomMemoryMaximum                  = 4
)

type DomainKeycodeSet uint32

const (
	DomKeycodeSetLinux DomainKeycodeSet = iota
	DomKeycodeSetXT
	DomKeycodeSetATSet1
	DomKeycodeSetATSet2
	DomKeycodeSetATSet3
	DomKeycodeSetOSX
	DomKeycodeSetXTKbd
	DomKeycodeSetUSB
	DomKeycodeSetWin32
	DomKeycodeSetRFB
)

type DomainProcessSignal uint32

const (
	DomSIGNOP = iota
	DomSIGHUP
	DomSIGINT
	DomSIGQUIT
	DomSIGILL
	DomSIGTRAP
	DomSIGABRT
	DomSIGBUS
	DomSIGFPE
	DomSIGKILL
	DomSIGUSR1
	DomSIGSEGV
	DomSIGUSR2
	DomSIGPIPE
	DomSIGALRM
	DomSIGTERM
	DomSIGSTKFLT
	DomSIGCHLD
	DomSIGCONT
	DomSIGSTOP
	DomSIGTSTP
	DomSIGTTIN
	DomSIGTTOU
	DomSIGURG
	DomSIGXCPU
	DomSIGXFSZ
	DomSIGVTALRM
	DomSIGPROF
	DomSIGWINCH
	DomSIGPOLL
	DomSIGPWR
	DomSIGSYS
	DomSIGRT0
	DomSIGRT1
	DomSIGRT2
	DomSIGRT3
	DomSIGRT4
	DomSIGRT5
	DomSIGRT6
	DomSIGRT7
	DomSIGRT8
	DomSIGRT9
	DomSIGRT10
	DomSIGRT11
	DomSIGRT12
	DomSIGRT13
	DomSIGRT14
	DomSIGRT15
	DomSIGRT16
	DomSIGRT17
	DomSIGRT18
	DomSIGRT19
	DomSIGRT20
	DomSIGRT21
	DomSIGRT22
	DomSIGRT23
	DomSIGRT24
	DomSIGRT25
	DomSIGRT26
	DomSIGRT27
	DomSIGRT28
	DomSIGRT29
	DomSIGRT30
	DomSIGRT31
	DomSIGRT32
)

type Domain struct {
	log       *log.Logger
	virDomain C.virDomainPtr
}

// Free frees the domain object. The running instance is kept alive. The data
// structure is freed and should not be used thereafter.
func (dom Domain) Free() error {
	dom.log.Println("freeing domain object...")
	cRet := C.virDomainFree(dom.virDomain)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("domain freed")

	return nil
}

// Autostart provides a boolean value indicating whether the domain configured
// to be automatically started when the host machine boots.
func (dom Domain) Autostart() bool {
	var cAutostart C.int
	dom.log.Println("checking whether domain autostarts...")
	cRet := C.virDomainGetAutostart(dom.virDomain, &cAutostart)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return false
	}

	autostart := (int32(cAutostart) == 1)
	if autostart {
		dom.log.Println("domain autostarts")
	} else {
		dom.log.Println("domain does not autostart")
	}

	return autostart
}

// HasCurrentSnapshot determines if the domain has a current snapshot.
func (dom Domain) HasCurrentSnapshot() bool {
	dom.log.Println("checking whether domain has current snapshot...")
	cRet := C.virDomainHasCurrentSnapshot(dom.virDomain, 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return false
	}

	hasCurrentSnapshot := (ret == 1)

	if hasCurrentSnapshot {
		dom.log.Println("domain has current snapshot")
	} else {
		dom.log.Println("domain does not have current snapshot")
	}

	return hasCurrentSnapshot
}

// HasManagedSaveImage checks if a domain has a managed save image as created
// by ManagedSave(). Note that any running domain should not have such an
// image, as it should have been removed on restart.
func (dom Domain) HasManagedSaveImage() bool {
	dom.log.Println("checking whether domain has managed save...")
	cRet := C.virDomainHasManagedSaveImage(dom.virDomain, 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return false
	}

	hasManagedSave := (ret == 1)

	if hasManagedSave {
		dom.log.Println("domain has managed save")
	} else {
		dom.log.Println("domain does not have managed save")
	}

	return hasManagedSave
}

// IsActive determines if the domain is currently running.
func (dom Domain) IsActive() bool {
	dom.log.Println("checking whether domain is active...")
	cRet := C.virDomainIsActive(dom.virDomain)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return false
	}

	active := (ret == 1)
	if active {
		dom.log.Println("domain is active")
	} else {
		dom.log.Println("domain is not active")
	}

	return active
}

// IsPersistent determines if the domain has a persistent configuration which
// means it will still exist after shutting down
func (dom Domain) IsPersistent() bool {
	dom.log.Println("checking whether domain is persistent...")
	cRet := C.virDomainIsPersistent(dom.virDomain)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return false
	}

	persistent := (ret == 1)

	if persistent {
		dom.log.Println("domain is persistent")
	} else {
		dom.log.Println("domain is not persistent")
	}

	return persistent
}

// IsUpdated determines if the domain has been updated.
func (dom Domain) IsUpdated() bool {
	dom.log.Println("checking whether domain is updated...")
	cRet := C.virDomainIsUpdated(dom.virDomain)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return false
	}

	updated := (ret == 1)

	if updated {
		dom.log.Println("domain is updated")
	} else {
		dom.log.Println("domain is not updated")
	}

	return updated
}

// OSType gets the type of domain operation system.
func (dom Domain) OSType() (string, error) {
	dom.log.Println("reading domain OS type...")
	cOS := C.virDomainGetOSType(dom.virDomain)
	if cOS == nil {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return "", err
	}
	defer C.free(unsafe.Pointer(cOS))

	os := C.GoString(cOS)
	dom.log.Printf("OS type: %v\n", os)

	return os, nil
}

// Name gets the public name for that domain.
func (dom Domain) Name() string {
	dom.log.Println("reading domain name...")
	cName := C.virDomainGetName(dom.virDomain)

	name := C.GoString(cName)
	dom.log.Printf("domain name: %v\n", name)

	return name
}

// Hostname gets the hostname for that domain.
func (dom Domain) Hostname() (string, error) {
	dom.log.Println("reading domain hostname...")
	cHostname := C.virDomainGetHostname(dom.virDomain, 0)
	if cHostname == nil {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return "", err
	}
	defer C.free(unsafe.Pointer(cHostname))

	hostname := C.GoString(cHostname)
	dom.log.Printf("domain hostname: %v\n", hostname)

	return hostname, nil
}

// ID gets the hypervisor ID number for the domain.
func (dom Domain) ID() (uint32, error) {
	dom.log.Println("reading domain ID...")
	cID := C.virDomainGetID(dom.virDomain)
	id := uint32(cID)

	if id == ^uint32(0) { // Go: ^uint32(0) == C: (unsigned int) -1
		err := errors.New("domain doesn't have an ID")
		dom.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	dom.log.Printf("domain ID: %v\n", id)

	return id, nil
}

// UUID gets the UUID for a domain as string. For more information about UUID
// see RFC4122.
func (dom Domain) UUID() (string, error) {
	cUUID := (*C.char)(C.malloc(C.size_t(C.VIR_UUID_STRING_BUFLEN)))
	defer C.free(unsafe.Pointer(cUUID))

	dom.log.Println("reading domain UUID...")
	cRet := C.virDomainGetUUIDString(dom.virDomain, cUUID)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return "", err
	}

	uuid := C.GoString(cUUID)
	dom.log.Printf("UUID: %v\n", uuid)

	return uuid, nil
}

// XML provides an XML description of the domain. The description may be reused
// later to relaunch the domain with CreateXML().
func (dom Domain) XML(typ DomainXMLFlag) (string, error) {
	dom.log.Printf("reading domain XML (flags = %v)...\n", typ)
	cXML := C.virDomainGetXMLDesc(dom.virDomain, C.uint(typ))
	if cXML == nil {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return "", err
	}
	defer C.free(unsafe.Pointer(cXML))

	xml := C.GoString(cXML)
	dom.log.Printf("XML length: %v runes\n", utf8.RuneCountInString(xml))

	return xml, nil
}

// Metadata retrieves the appropriate domain element given by "type".
func (dom Domain) Metadata(typ DomainMetadataType, xmlns string, impact DomainModificationImpact) (string, error) {
	cXMLNS := C.CString(xmlns)
	defer C.free(unsafe.Pointer(cXMLNS))

	dom.log.Printf("reading domain metadata (type = %v, namespace = %v, impact = %v)...\n", typ, xmlns, impact)
	cMetadata := C.virDomainGetMetadata(dom.virDomain, C.int(typ), cXMLNS, C.uint(impact))
	if cMetadata == nil {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return "", err
	}
	defer C.free(unsafe.Pointer(cMetadata))

	metadata := C.GoString(cMetadata)
	dom.log.Printf("metadata XML length: %v runes\n", utf8.RuneCountInString(metadata))

	return metadata, nil
}

// Destroy destroys the domain object. The running instance is shutdown if not
// down already and all resources used by it are given back to the hypervisor.
// This does not free the associated virDomainPtr object. This function may
// require privileged access.
func (dom Domain) Destroy(flags DomainDestroyFlag) error {
	dom.log.Printf("destroying domain (flags = %v)...\n", flags)
	cRet := C.virDomainDestroyFlags(dom.virDomain, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("domain destroyed")

	return nil
}

// Create launches a defined domain. If the call succeeds the domain moves from
// the defined to the running domains pools.
func (dom Domain) Create(flags DomainCreateFlag) error {
	dom.log.Printf("starting domain (flags = %v)...\n", flags)
	cRet := C.virDomainCreateWithFlags(dom.virDomain, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("domain started")

	return nil
}

// Undefine undefines a domain. If the domain is running, it's converted to
// transient domain, without stopping it. If the domain is inactive, the domain
// configuration is removed.
func (dom Domain) Undefine(flags DomainUndefineFlag) error {
	dom.log.Printf("undefining domain (flags = %v)...\n", flags)
	cRet := C.virDomainUndefineFlags(dom.virDomain, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("domain undefined")

	return nil
}

// Reboot reboots a domain, the domain object is still usable thereafter, but
// the domain OS is being stopped for a restart. Note that the guest OS may
// ignore the request. Additionally, the hypervisor may check and support the
// domain 'on_reboot' XML setting resulting in a domain that shuts down instead
// of rebooting.
func (dom Domain) Reboot(flags DomainRebootFlag) error {
	dom.log.Printf("rebooting domain (flags = %v)...\n", flags)
	cRet := C.virDomainReboot(dom.virDomain, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("domain rebooted")

	return nil
}

// Reset resets a domain immediately without any guest OS shutdown. Reset
// emulates the power reset button on a machine, where all hardware sees the
// RST line set and reinitializes internal state.
// Note that there is a risk of data loss caused by reset without any guest
// OS shutdown.
func (dom Domain) Reset() error {
	dom.log.Println("resetting domain...")
	cRet := C.virDomainReset(dom.virDomain, 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("domain reset")

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
func (dom Domain) Shutdown() error {
	dom.log.Println("shutting down domain...")
	cRet := C.virDomainShutdown(dom.virDomain)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("domain shut down")

	return nil
}

// State extracts domain state. Each state can be accompanied with a reason
// (if known) which led to the state.
func (dom Domain) State() (DomainState, int32, error) {
	var cState, cReason C.int
	dom.log.Println("reading domain state...")
	cRet := C.virDomainGetState(dom.virDomain, &cState, &cReason, 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return 0, 0, err
	}

	state := DomainState(cState)
	reason := int32(cReason)
	dom.log.Printf("state: %v (reason = %v)\n", state, reason)

	return state, reason, nil
}

// Suspend suspends an active domain, the process is frozen without further
// access to CPU resources and I/O but the memory used by the domain at the
// hypervisor level will stay allocated. Use Resume() to reactivate the domain.
// This function may require privileged access. Moreover, suspend may not be
// supported if domain is in some special state like DomStatePMSuspended.
func (dom Domain) Suspend() error {
	dom.log.Println("suspending domain...")
	cRet := C.virDomainSuspend(dom.virDomain)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("domain suspended")

	return nil
}

// Resume resumes a suspended domain, the process is restarted from the state
// where it was frozen by calling Suspend(). This function may require
// privileged access. Moreover, resume may not be supported if domain is in
// some special state like DomStatePMSuspended.
func (dom Domain) Resume() error {
	dom.log.Println("resuming domain...")
	cRet := C.virDomainResume(dom.virDomain)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("domain resumed")

	return nil
}

// CoreDump dumps the core of a domain on a given file for analysis. Note that
// for remote Xen Daemon the file path will be interpreted in the remote host.
// Hypervisors may require the user to manually ensure proper permissions on
// the file named by "to".
func (dom Domain) CoreDump(file string, flags DomainDumpFlag) error {
	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))

	dom.log.Printf("dumping domain's core to file %v (flags = %v)...", file, flags)
	cRet := C.virDomainCoreDump(dom.virDomain, cFile, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("core dump saved")

	return nil
}

// Ref increments the reference count on the domain. For each additional call
// to this method, there shall be a corresponding call to virDomainFree to
// release the reference count, once the caller no longer needs the reference
// to this object.
func (dom Domain) Ref() error {
	dom.log.Println("incrementing domain's reference count...")
	cRet := C.virDomainRef(dom.virDomain)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("reference count incremented")

	return nil
}

// MaxMemory retrieves the maximum amount of physical memory allocated to
// a domain.
func (dom Domain) MaxMemory() (uint64, error) {
	dom.log.Println("reading domain maximum memory...")
	cRet := C.virDomainGetMaxMemory(dom.virDomain)
	ret := uint64(cRet)

	if ret == 0 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	dom.log.Printf("max memory: %v kiB\n", ret)

	return ret, nil
}

// VCPUs queries the number of virtual CPUs used by the domain. Note that this
// call may fail if the underlying virtualization hypervisor does not support
// it. This function may require privileged access to the hypervisor.
func (dom Domain) VCPUs(flags DomainVCPUsFlag) (int32, error) {
	dom.log.Println("reading domain VCPUs count...")
	cRet := C.virDomainGetVcpusFlags(dom.virDomain, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	dom.log.Printf("VCPUs count: %v\n", ret)

	return ret, nil
}

// InfoState extracts the state of the domain.
func (dom Domain) InfoState() (DomainState, error) {
	var cInfo C.virDomainInfo
	cRet := C.virDomainGetInfo(dom.virDomain, &cInfo)
	ret := int32(cRet)

	if ret == -1 {
		return 0, LastError()
	}

	return DomainState(cInfo.state), nil
}

// InfoMaxMemory extracts the maximum memory in KBytes allowed in the domain.
func (dom Domain) InfoMaxMemory() (uint64, error) {
	var cInfo C.virDomainInfo
	cRet := C.virDomainGetInfo(dom.virDomain, &cInfo)
	ret := int32(cRet)

	if ret == -1 {
		return 0, LastError()
	}

	return uint64(cInfo.maxMem), nil
}

// InfoMemory extracts the memory in KBytes used by the domain.
func (dom Domain) InfoMemory() (uint64, error) {
	var cInfo C.virDomainInfo
	cRet := C.virDomainGetInfo(dom.virDomain, &cInfo)
	ret := int32(cRet)

	if ret == -1 {
		return 0, LastError()
	}

	return uint64(cInfo.memory), nil
}

// InfoVCPUs extracts the number of virtual CPUs for the domain.
func (dom Domain) InfoVCPUs() (uint16, error) {
	var cInfo C.virDomainInfo
	cRet := C.virDomainGetInfo(dom.virDomain, &cInfo)
	ret := int32(cRet)

	if ret == -1 {
		return 0, LastError()
	}

	return uint16(cInfo.nrVirtCpu), nil
}

// InfoCPUTime extracts the CPU time used in nanoseconds.
func (dom Domain) InfoCPUTime() (uint64, error) {
	var cInfo C.virDomainInfo
	cRet := C.virDomainGetInfo(dom.virDomain, &cInfo)
	ret := int32(cRet)

	if ret == -1 {
		return 0, LastError()
	}

	return uint64(cInfo.cpuTime), nil
}

// Save suspends a domain and save its memory contents to a file on disk. After
// the call, if successful, the domain is not listed as running anymore (this
// ends the life of a transient domain). Use Restore() to restore a domain
// after saving.
func (dom Domain) Save(to string, xml string, flags DomainSaveFlag) error {
	cTo := C.CString(to)
	defer C.free(unsafe.Pointer(cTo))

	var cXML *C.char
	if xml != "" {
		cXML = C.CString(xml)
		defer C.free(unsafe.Pointer(cXML))
	} else {
		cXML = nil
	}

	dom.log.Printf("saving domain's memory to file %v (flags = %v)...\n", to, flags)
	cRet := C.virDomainSaveFlags(dom.virDomain, cTo, cXML, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("domain saved")

	return nil
}

// AttachDevice attaches a virtual device to a domain, using the flags
// parameter to control how the device is attached. DomDeviceModifyCurrent
// specifies that the device allocation is made based on current domain state.
// DomDeviceModifyLive specifies that the device shall be allocated to the
// active domain instance only and is not added to the persisted domain
// configuration. DomDeviceModifyConfig specifies that the device shall be
// allocated to the persisted domain configuration only. Note that the target
// hypervisor must return an error if unable to satisfy flags. E.g. the
// hypervisor driver will return failure if DomDeviceModifyLive is specified
// but it only supports modifying the persisted device allocation.
func (dom Domain) AttachDevice(deviceXML string, flags DomainDeviceModifyFlag) error {
	cXML := C.CString(deviceXML)
	defer C.free(unsafe.Pointer(cXML))

	dom.log.Printf("attaching a virtual device to domain (flags = %v)...\n", flags)
	cRet := C.virDomainAttachDeviceFlags(dom.virDomain, cXML, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("device attached")

	return nil
}

// DetachDevice detaches a virtual device from a domain, using the flags
// parameter to control how the device is detached. DomDeviceModifyCurrent
// specifies that the device allocation is removed based on current domain
// state. DomDeviceModifyLive specifies that the device shall be deallocated
// from the active domain instance only and is not from the persisted domain
// configuration. DomDeviceModifyConfig specifies that the device shall be
// deallocated from the persisted domain configuration only. Note that the
// target hypervisor must return an error if unable to satisfy flags. E.g. the
// hypervisor driver will return failure if DomDeviceModifyLive is specified
// but it only supports removing the persisted device allocation.
func (dom Domain) DetachDevice(deviceXML string, flags DomainDeviceModifyFlag) error {
	cXML := C.CString(deviceXML)
	defer C.free(unsafe.Pointer(cXML))

	dom.log.Printf("detaching a virtual device from domain (flags = %v)...\n", flags)
	cRet := C.virDomainDetachDeviceFlags(dom.virDomain, cXML, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("device detached")

	return nil
}

// UpdateDevice changes a virtual device on a domain, using the flags parameter
// to control how the device is changed. DomDeviceModifyCurrent specifies that
// the device change is made based on current domain state. DomDeviceModifyLive
// specifies that the device shall be changed on the active domain instance
// only and is not added to the persisted domain configuration.
// DomDeviceModifyConfig specifies that the device shall be changed on the
// persisted domain configuration only. Note that the target hypervisor must
// return an error if unable to satisfy flags. E.g. the hypervisor driver will
// return failure if DomDeviceModifyLive is specified but it only supports
// modifying the persisted device allocation.
func (dom Domain) UpdateDevice(deviceXML string, flags DomainDeviceModifyFlag) error {
	cXML := C.CString(deviceXML)
	defer C.free(unsafe.Pointer(cXML))

	dom.log.Printf("updating a virtual device on domain (flags = %v)...\n", flags)
	cRet := C.virDomainUpdateDeviceFlags(dom.virDomain, cXML, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("device updated")

	return nil
}

// SetAutostart configures the domain to be automatically started when the host
// machine boots.
func (dom Domain) SetAutostart(autostart bool) error {
	var cAutostart C.int
	if autostart {
		dom.log.Println("enabling domain autostart...")
		cAutostart = 1
	} else {
		dom.log.Println("disabling domain autostart...")
		cAutostart = 0
	}

	cRet := C.virDomainSetAutostart(dom.virDomain, cAutostart)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	if autostart {
		dom.log.Println("autostart enabled")
	} else {
		dom.log.Println("autostart disabled")
	}

	return nil
}

// SetMemory dynamically changes the target amount of physical memory allocated
// to a domain. This function may require privileged access to the hypervisor.
func (dom Domain) SetMemory(memory uint64, flags DomainMemoryFlag) error {
	dom.log.Printf("changing domain memory to %v kiB (flags = %v)...\n", memory, flags)
	cRet := C.virDomainSetMemoryFlags(dom.virDomain, C.ulong(memory), C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("memory changed")

	return nil
}

// SetMetadata sets the appropriate domain element given by "type" to the value
// of "description". A "type" of DomMetaDescription is free-form text;
// DomMetaTitle is free-form, but no newlines are permitted, and should be
// short (although the length is not enforced). For these two options "key" and
// "uri" are irrelevant and must be set to NULL.
func (dom Domain) SetMetadata(typ DomainMetadataType, metadata string, key string, uri string, impact DomainModificationImpact) error {
	cMetadata := C.CString(metadata)
	defer C.free(unsafe.Pointer(cMetadata))

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cURI := C.CString(uri)
	defer C.free(unsafe.Pointer(cURI))

	dom.log.Printf("changing domain metadata key '<%v:%v>' (type = %v, impact = %v)...\n", key, uri, typ, impact)
	cRet := C.virDomainSetMetadata(dom.virDomain, C.int(typ), cMetadata, cKey, cURI, C.uint(impact))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("metadata changed")

	return nil
}

// SetVCPUs dynamically changes the number of virtual CPUs used by the domain.
// Note that this call may fail if the underlying virtualization hypervisor
// does not support it or if growing the number is arbitrary limited. This
// function may require privileged access to the hypervisor.
func (dom Domain) SetVCPUs(vcpus uint32, flags DomainVCPUsFlag) error {
	dom.log.Printf("changing domain VCPUs count to %v (flags = %v)...\n", vcpus, flags)
	cRet := C.virDomainSetVcpusFlags(dom.virDomain, C.uint(vcpus), C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("VCPUs count changed")

	return nil
}

// ManagedSave suspends a domain and save its memory contents to a file on
// disk. After the call, if successful, the domain is not listed as running
// anymore. The difference from Save() is that libvirt is keeping track of the
// saved state itself, and will reuse it once the domain is being restarted
// (automatically or via an explicit libvirt call). As a result any running
// domain is sure to not have a managed saved image. This also implies that
// managed save only works on persistent domains, since the domain must still
// exist in order to use Create() to restart it.
func (dom Domain) ManagedSave(flags DomainSaveFlag) error {
	dom.log.Printf("saving domain's memory to a libvirt-managed location (flags = %v)...\n", flags)
	cRet := C.virDomainManagedSave(dom.virDomain, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("domain saved")

	return nil
}

// ManagedSaveRemove removes any managed save image for this domain.
func (dom Domain) ManagedSaveRemove() error {
	dom.log.Println("removing libvirt-managed domain save image...")
	cRet := C.virDomainManagedSaveRemove(dom.virDomain, 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("save image removed")

	return nil
}

// SendKey send key(s) to the guest.
func (dom Domain) SendKey(codeSet DomainKeycodeSet, hold time.Duration, keycodes []uint32) error {
	dom.log.Printf("sending keys %v (keycode set = %v) to domain during %v...\n", keycodes, codeSet, time.Duration(hold))
	cRet := C.virDomainSendKey(dom.virDomain, C.uint(codeSet), C.uint(hold*time.Millisecond), (*C.uint)(unsafe.Pointer(&keycodes[0])), C.int(len(keycodes)), 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("keys sent")

	return nil
}

// SendProcessSignal sends a signal to the designated process in the guest.
func (dom Domain) SendProcessSignal(pid int64, signal DomainProcessSignal) error {
	dom.log.Printf("sending signal %v to domain's process %v...\n", signal, pid)
	cRet := C.virDomainSendProcessSignal(dom.virDomain, C.longlong(pid), C.uint(signal), 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return err
	}

	dom.log.Println("signal sent")

	return nil
}
