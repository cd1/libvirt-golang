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

type DomainState uint

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

type DomainNostateReason uint

const (
	DomNostateReasonUnknown DomainNostateReason = iota
)

type DomainRunningReason uint

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

type DomainBlockedReason uint

const (
	DomBlockedReasonUnkwown DomainBlockedReason = iota
)

type DomainPausedReason uint

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

type DomainShutdownReason uint

const (
	DomShutdownReasonUnknown DomainShutdownReason = iota
	DomShutdownReasonUser
)

type DomainShutoffReason uint

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

type DomainCrashedReason uint

const (
	DomCrashedReasonUnknown DomainCrashedReason = iota
	DomCrashedReasonPanicked
)

type DomainPMSuspendedReason uint

const (
	DomPMSuspendedReasonUnknown DomainPMSuspendedReason = iota
)

type DomainDumpFlag uint

const (
	DomDumpCrash DomainDumpFlag = (1 << iota)
	DomDumpLive
	DomDumpBypassCache
	DomDumpReset
	DomDumpMemoryOnly
	DomDumpDefault = 0
)

type DomainVCPUsFlag uint

const (
	DomVCPusConfig  DomainVCPUsFlag = DomAffectConfig
	DomVCPUsCurrent                 = DomAffectCurrent
	DomVCPUsLive                    = DomAffectLive
	DomVCPUsMaximum                 = 4
	DomVCPUsGuest                   = 8
)

type DomainSaveFlag uint

const (
	DomSaveBypassCache DomainSaveFlag = (1 << iota)
	DomSaveRunning
	DomSavePaused
	DomSaveDefault = 0
)

type DomainDeviceModifyFlag uint

const (
	DomDeviceModifyConfig  DomainDeviceModifyFlag = DomAffectConfig
	DomDeviceModifyCurrent                        = DomAffectCurrent
	DomDeviceModifyLive                           = DomAffectLive
	DomDeviceModifyForce                          = 4
)

type DomainMemoryFlag uint

const (
	DomMemoryConfig  DomainMemoryFlag = DomAffectConfig
	DomMemoryCurrent                  = DomAffectCurrent
	DomMemoryLive                     = DomAffectLive
	DomMemoryMaximum                  = 4
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

// State extracts domain state. Each state can be accompanied with a reason
// (if known) which led to the state.
func (dom Domain) State() (DomainState, int, *Error) {
	var cState, cReason C.int
	cRet := C.virDomainGetState(dom.virDomain, &cState, &cReason, 0)
	ret := int(cRet)

	if ret == -1 {
		return 0, 0, LastError()
	}

	return DomainState(cState), int(cReason), nil
}

// Suspend suspends an active domain, the process is frozen without further
// access to CPU resources and I/O but the memory used by the domain at the
// hypervisor level will stay allocated. Use Resume() to reactivate the domain.
// This function may require privileged access. Moreover, suspend may not be
// supported if domain is in some special state like DomStatePMSuspended.
func (dom Domain) Suspend() *Error {
	cRet := C.virDomainSuspend(dom.virDomain)
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}

// Resume resumes a suspended domain, the process is restarted from the state
// where it was frozen by calling Suspend(). This function may require
// privileged access. Moreover, resume may not be supported if domain is in
// some special state like DomStatePMSuspended.
func (dom Domain) Resume() *Error {
	cRet := C.virDomainResume(dom.virDomain)
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}

// CoreDump dumps the core of a domain on a given file for analysis. Note that
// for remote Xen Daemon the file path will be interpreted in the remote host.
// Hypervisors may require the user to manually ensure proper permissions on
// the file named by "to".
func (dom Domain) CoreDump(file string, flags DomainDumpFlag) *Error {
	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))

	cRet := C.virDomainCoreDump(dom.virDomain, cFile, C.uint(flags))
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}

// Ref increments the reference count on the domain. For each additional call
// to this method, there shall be a corresponding call to virDomainFree to
// release the reference count, once the caller no longer needs the reference
// to this object.
func (dom Domain) Ref() *Error {
	cRet := C.virDomainRef(dom.virDomain)
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}

// MaxMemory retrieves the maximum amount of physical memory allocated to
// a domain.
func (dom Domain) MaxMemory() (uint64, *Error) {
	cRet := C.virDomainGetMaxMemory(dom.virDomain)
	ret := uint64(cRet)

	if ret == 0 {
		return 0, LastError()
	}

	return ret, nil
}

// VCPUs queries the number of virtual CPUs used by the domain. Note that this
// call may fail if the underlying virtualization hypervisor does not support
// it. This function may require privileged access to the hypervisor.
func (dom Domain) VCPUs(flags DomainVCPUsFlag) (int, *Error) {
	cRet := C.virDomainGetVcpusFlags(dom.virDomain, C.uint(flags))
	ret := int(cRet)

	if ret == -1 {
		return 0, LastError()
	}

	return ret, nil
}

// InfoState extracts the state of the domain.
func (dom Domain) InfoState() (DomainState, *Error) {
	var cInfo C.virDomainInfo
	cRet := C.virDomainGetInfo(dom.virDomain, &cInfo)
	ret := int(cRet)

	if ret == -1 {
		return 0, LastError()
	}

	return DomainState(cInfo.state), nil
}

// InfoMaxMemory extracts the maximum memory in KBytes allowed in the domain.
func (dom Domain) InfoMaxMemory() (uint64, *Error) {
	var cInfo C.virDomainInfo
	cRet := C.virDomainGetInfo(dom.virDomain, &cInfo)
	ret := int(cRet)

	if ret == -1 {
		return 0, LastError()
	}

	return uint64(cInfo.maxMem), nil
}

// InfoMemory extracts the memory in KBytes used by the domain.
func (dom Domain) InfoMemory() (uint64, *Error) {
	var cInfo C.virDomainInfo
	cRet := C.virDomainGetInfo(dom.virDomain, &cInfo)
	ret := int(cRet)

	if ret == -1 {
		return 0, LastError()
	}

	return uint64(cInfo.memory), nil
}

// InfoVCPUs extracts the number of virtual CPUs for the domain.
func (dom Domain) InfoVCPUs() (uint16, *Error) {
	var cInfo C.virDomainInfo
	cRet := C.virDomainGetInfo(dom.virDomain, &cInfo)
	ret := int(cRet)

	if ret == -1 {
		return 0, LastError()
	}

	return uint16(cInfo.nrVirtCpu), nil
}

// InfoCPUTime extracts the CPU time used in nanoseconds.
func (dom Domain) InfoCPUTime() (uint64, *Error) {
	var cInfo C.virDomainInfo
	cRet := C.virDomainGetInfo(dom.virDomain, &cInfo)
	ret := int(cRet)

	if ret == -1 {
		return 0, LastError()
	}

	return uint64(cInfo.cpuTime), nil
}

// Save suspends a domain and save its memory contents to a file on disk. After
// the call, if successful, the domain is not listed as running anymore (this
// ends the life of a transient domain). Use Restore() to restore a domain
// after saving.
func (dom Domain) Save(to string, xml string, flags DomainSaveFlag) *Error {
	cTo := C.CString(to)
	defer C.free(unsafe.Pointer(cTo))

	var cXML *C.char
	if xml != "" {
		cXML = C.CString(xml)
		defer C.free(unsafe.Pointer(cXML))
	} else {
		cXML = nil
	}

	cRet := C.virDomainSaveFlags(dom.virDomain, cTo, cXML, C.uint(flags))
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

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
func (dom Domain) AttachDevice(deviceXML string, flags DomainDeviceModifyFlag) *Error {
	cXML := C.CString(deviceXML)
	defer C.free(unsafe.Pointer(cXML))

	cRet := C.virDomainAttachDeviceFlags(dom.virDomain, cXML, C.uint(flags))
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

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
func (dom Domain) DetachDevice(deviceXML string, flags DomainDeviceModifyFlag) *Error {
	cXML := C.CString(deviceXML)
	defer C.free(unsafe.Pointer(cXML))

	cRet := C.virDomainDetachDeviceFlags(dom.virDomain, cXML, C.uint(flags))
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

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
func (dom Domain) UpdateDevice(deviceXML string, flags DomainDeviceModifyFlag) *Error {
	cXML := C.CString(deviceXML)
	defer C.free(unsafe.Pointer(cXML))

	cRet := C.virDomainUpdateDeviceFlags(dom.virDomain, cXML, C.uint(flags))
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}

// SetAutostart configures the domain to be automatically started when the host
// machine boots.
func (dom Domain) SetAutostart(autostart bool) *Error {
	var cAutostart C.int
	if autostart {
		cAutostart = 1
	} else {
		cAutostart = 0
	}

	cRet := C.virDomainSetAutostart(dom.virDomain, cAutostart)
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}

// SetMemory dynamically changes the target amount of physical memory allocated
// to a domain. This function may require privileged access to the hypervisor.
func (dom Domain) SetMemory(memory uint64, flags DomainMemoryFlag) *Error {
	cRet := C.virDomainSetMemoryFlags(dom.virDomain, C.ulong(memory), C.uint(flags))
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}

// SetMetadata sets the appropriate domain element given by "type" to the value
// of "description". A "type" of DomMetaDescription is free-form text;
// DomMetaTitle is free-form, but no newlines are permitted, and should be
// short (although the length is not enforced). For these two options "key" and
// "uri" are irrelevant and must be set to NULL.
func (dom Domain) SetMetadata(typ DomainMetadataType, metadata string, key string, uri string, impact DomainModificationImpact) *Error {
	cMetadata := C.CString(metadata)
	defer C.free(unsafe.Pointer(cMetadata))

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cURI := C.CString(uri)
	defer C.free(unsafe.Pointer(cURI))

	cRet := C.virDomainSetMetadata(dom.virDomain, C.int(typ), cMetadata, cKey, cURI, C.uint(impact))
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}

// SetVCPUs dynamically changes the number of virtual CPUs used by the domain.
// Note that this call may fail if the underlying virtualization hypervisor
// does not support it or if growing the number is arbitrary limited. This
// function may require privileged access to the hypervisor.
func (dom Domain) SetVCPUs(vcpus uint, flags DomainVCPUsFlag) *Error {
	cRet := C.virDomainSetVcpusFlags(dom.virDomain, C.uint(vcpus), C.uint(flags))
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}
