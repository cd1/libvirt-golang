package libvirt

// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"
import (
	"errors"
	"log"
	"reflect"
	"time"
	"unicode/utf8"
	"unsafe"
)

// DomainListFlag defines a filter when listing domains.
type DomainListFlag uint32

// Possible values for DomainListFlag.
const (
	DomListAll           DomainListFlag = 0
	DomListActive        DomainListFlag = C.VIR_CONNECT_LIST_DOMAINS_ACTIVE
	DomListInactive      DomainListFlag = C.VIR_CONNECT_LIST_DOMAINS_INACTIVE
	DomListPersistent    DomainListFlag = C.VIR_CONNECT_LIST_DOMAINS_PERSISTENT
	DomListTransient     DomainListFlag = C.VIR_CONNECT_LIST_DOMAINS_TRANSIENT
	DomListRunning       DomainListFlag = C.VIR_CONNECT_LIST_DOMAINS_RUNNING
	DomListPaused        DomainListFlag = C.VIR_CONNECT_LIST_DOMAINS_PAUSED
	DomListShutOff       DomainListFlag = C.VIR_CONNECT_LIST_DOMAINS_SHUTOFF
	DomListOther         DomainListFlag = C.VIR_CONNECT_LIST_DOMAINS_OTHER
	DomListManagedSave   DomainListFlag = C.VIR_CONNECT_LIST_DOMAINS_MANAGEDSAVE
	DomListNoManagedSave DomainListFlag = C.VIR_CONNECT_LIST_DOMAINS_NO_MANAGEDSAVE
	DomListAutostart     DomainListFlag = C.VIR_CONNECT_LIST_DOMAINS_AUTOSTART
	DomListNoAutostart   DomainListFlag = C.VIR_CONNECT_LIST_DOMAINS_NO_AUTOSTART
	DomListHasSnapshot   DomainListFlag = C.VIR_CONNECT_LIST_DOMAINS_HAS_SNAPSHOT
	DomListNoSnapshot    DomainListFlag = C.VIR_CONNECT_LIST_DOMAINS_NO_SNAPSHOT
)

// DomainMetadataType defines a type of metadata element.
type DomainMetadataType uint32

// Possible values for DomainMetadataType.
const (
	DomMetaDescription DomainMetadataType = C.VIR_DOMAIN_METADATA_DESCRIPTION
	DomMetaTitle       DomainMetadataType = C.VIR_DOMAIN_METADATA_TITLE
	DomMetaElement     DomainMetadataType = C.VIR_DOMAIN_METADATA_ELEMENT
)

// DomainModificationImpact controls whether the live domain or persistent
// configuration (or both) will be queried.
type DomainModificationImpact uint32

// Possible values for DomainModificationImpact
const (
	DomAffectCurrent DomainModificationImpact = C.VIR_DOMAIN_AFFECT_CURRENT
	DomAffectLive    DomainModificationImpact = C.VIR_DOMAIN_AFFECT_LIVE
	DomAffectConfig  DomainModificationImpact = C.VIR_DOMAIN_AFFECT_CONFIG
)

// DomainXMLFlag defines how the XML content should be read from a domain.
type DomainXMLFlag uint32

// Possible values for DomainXMLFlag.
const (
	DomXMLDefault    DomainXMLFlag = 0
	DomXMLSecure     DomainXMLFlag = C.VIR_DOMAIN_XML_SECURE
	DomXMLInactive   DomainXMLFlag = C.VIR_DOMAIN_XML_INACTIVE
	DomXMLUpdateCPU  DomainXMLFlag = C.VIR_DOMAIN_XML_UPDATE_CPU
	DomXMLMigratable DomainXMLFlag = C.VIR_DOMAIN_XML_MIGRATABLE
)

// DomainCreateFlag defines how a domain should be created.
type DomainCreateFlag uint32

// Possible values for DomainCreateFlag.
const (
	DomCreateDefault     DomainCreateFlag = C.VIR_DOMAIN_NONE
	DomCreatePaused      DomainCreateFlag = C.VIR_DOMAIN_START_PAUSED
	DomCreateAutodestroy DomainCreateFlag = C.VIR_DOMAIN_START_AUTODESTROY
	DomCreateBypassCache DomainCreateFlag = C.VIR_DOMAIN_START_BYPASS_CACHE
	DomCreateForceBoot   DomainCreateFlag = C.VIR_DOMAIN_START_FORCE_BOOT
)

// DomainDestroyFlag defines how a domain should be destroyed.
type DomainDestroyFlag uint32

// Possible values for DomainDestroyFlag.
const (
	DomDestroyDefault  DomainDestroyFlag = C.VIR_DOMAIN_DESTROY_DEFAULT
	DomDestroyGraceful DomainDestroyFlag = C.VIR_DOMAIN_DESTROY_GRACEFUL
)

// DomainUndefineFlag defines how a domain should be undefined.
type DomainUndefineFlag uint32

// Possible values for DomainUndefineFlag.
const (
	DomUndefineDefault           DomainUndefineFlag = 0
	DomUndefineManagedSave       DomainUndefineFlag = C.VIR_DOMAIN_UNDEFINE_MANAGED_SAVE
	DomUndefineSnapshotsMetadata DomainUndefineFlag = C.VIR_DOMAIN_UNDEFINE_SNAPSHOTS_METADATA
	DomUndefineNVRAM             DomainUndefineFlag = C.VIR_DOMAIN_UNDEFINE_NVRAM
)

// DomainRebootFlag defines how a domain should be rebooted.
type DomainRebootFlag uint32

// Possible values for DomainRebootFlag.
const (
	DomRebootDefault      DomainRebootFlag = C.VIR_DOMAIN_REBOOT_DEFAULT
	DomRebootACPIPowerBtn DomainRebootFlag = C.VIR_DOMAIN_REBOOT_ACPI_POWER_BTN
	DomRebootGuestAgent   DomainRebootFlag = C.VIR_DOMAIN_REBOOT_GUEST_AGENT
	DomRebootInitctl      DomainRebootFlag = C.VIR_DOMAIN_REBOOT_INITCTL
	DomRebootSignal       DomainRebootFlag = C.VIR_DOMAIN_REBOOT_SIGNAL
	DomRebootParavirt     DomainRebootFlag = C.VIR_DOMAIN_REBOOT_PARAVIRT
)

// DomainState represents the state of a domain.
type DomainState uint32

// Possible values for DomainState.
const (
	DomStateNone        DomainState = C.VIR_DOMAIN_NOSTATE
	DomStateRunning     DomainState = C.VIR_DOMAIN_RUNNING
	DomStateBlocked     DomainState = C.VIR_DOMAIN_BLOCKED
	DomStatePaused      DomainState = C.VIR_DOMAIN_PAUSED
	DomStateShutdown    DomainState = C.VIR_DOMAIN_SHUTDOWN
	DomStateShutoff     DomainState = C.VIR_DOMAIN_SHUTOFF
	DomStateCrashed     DomainState = C.VIR_DOMAIN_CRASHED
	DomStatePMSuspended DomainState = C.VIR_DOMAIN_PMSUSPENDED
)

// DomainNostateReason describes the reason which led a domain to be on "DomStateNone".
type DomainNostateReason uint32

// Possible values for DomainNostateReason.
const (
	DomNostateReasonUnknown DomainNostateReason = C.VIR_DOMAIN_NOSTATE_UNKNOWN
)

// DomainRunningReason describes the reason which led a domain to be on "DomStateRunning".
type DomainRunningReason uint32

// Possible values for DomainRunningReason.
const (
	DomRunningReasonUnknown            DomainRunningReason = C.VIR_DOMAIN_RUNNING_UNKNOWN
	DomRunningReasonBooted             DomainRunningReason = C.VIR_DOMAIN_RUNNING_BOOTED
	DomRunningReasonMigrated           DomainRunningReason = C.VIR_DOMAIN_RUNNING_MIGRATED
	DomRunningReasonRestored           DomainRunningReason = C.VIR_DOMAIN_RUNNING_RESTORED
	DomRunningReasonFromSnapshot       DomainRunningReason = C.VIR_DOMAIN_RUNNING_FROM_SNAPSHOT
	DomRunningReasonUnpaused           DomainRunningReason = C.VIR_DOMAIN_RUNNING_UNPAUSED
	DomRunningReasonMigrationCancelled DomainRunningReason = C.VIR_DOMAIN_RUNNING_MIGRATION_CANCELED
	DomRunningReasonSaveCancelled      DomainRunningReason = C.VIR_DOMAIN_RUNNING_SAVE_CANCELED
	DomRunningReasonWakeUp             DomainRunningReason = C.VIR_DOMAIN_RUNNING_WAKEUP
	DomRunningReasonCrashed            DomainRunningReason = C.VIR_DOMAIN_RUNNING_CRASHED
)

// DomainBlockedReason describes the reason which led a domain to be on "DomStateBlocked".
type DomainBlockedReason uint32

// Possible values for DomainBlockedReason.
const (
	DomBlockedReasonUnkwown DomainBlockedReason = C.VIR_DOMAIN_BLOCKED_UNKNOWN
)

// DomainPausedReason describes the reason which led a domain to be on "DomStatePaused".
type DomainPausedReason uint32

// Possible values for DomainPausedReason.
const (
	DomPausedReasonUnknown      DomainPausedReason = C.VIR_DOMAIN_PAUSED_UNKNOWN
	DomPausedReasonUser         DomainPausedReason = C.VIR_DOMAIN_PAUSED_USER
	DomPausedReasonMigration    DomainPausedReason = C.VIR_DOMAIN_PAUSED_MIGRATION
	DomPausedReasonSave         DomainPausedReason = C.VIR_DOMAIN_PAUSED_SAVE
	DomPausedReasonDump         DomainPausedReason = C.VIR_DOMAIN_PAUSED_DUMP
	DomPausedReasonIOError      DomainPausedReason = C.VIR_DOMAIN_PAUSED_IOERROR
	DomPausedReasonWatchdog     DomainPausedReason = C.VIR_DOMAIN_PAUSED_WATCHDOG
	DomPausedReasonFromSnapshot DomainPausedReason = C.VIR_DOMAIN_PAUSED_FROM_SNAPSHOT
	DomPausedReasonShuttingDown DomainPausedReason = C.VIR_DOMAIN_PAUSED_SHUTTING_DOWN
	DomPausedReasonSnapshot     DomainPausedReason = C.VIR_DOMAIN_PAUSED_SNAPSHOT
	DomPausedReasonCrashed      DomainPausedReason = C.VIR_DOMAIN_PAUSED_CRASHED
)

// DomainShutdownReason describes the reason which led a domain to be on "DomStateShutdown".
type DomainShutdownReason uint32

// Possible values for DomainShutdownReason.
const (
	DomShutdownReasonUnknown DomainShutdownReason = C.VIR_DOMAIN_SHUTDOWN_UNKNOWN
	DomShutdownReasonUser    DomainShutdownReason = C.VIR_DOMAIN_SHUTDOWN_USER
)

// DomainShutoffReason describes the reason which led a domain to be on "DomStateShutoff".
type DomainShutoffReason uint32

// Possible values for DomainShutoffReason.
const (
	DomShutoffReasonUnknown      DomainShutoffReason = C.VIR_DOMAIN_SHUTOFF_UNKNOWN
	DomShutoffReasonShutdown     DomainShutoffReason = C.VIR_DOMAIN_SHUTOFF_SHUTDOWN
	DomShutoffReasonDestroyed    DomainShutoffReason = C.VIR_DOMAIN_SHUTOFF_DESTROYED
	DomShutoffReasonCrashed      DomainShutoffReason = C.VIR_DOMAIN_SHUTOFF_CRASHED
	DomShutoffReasonMigrated     DomainShutoffReason = C.VIR_DOMAIN_SHUTOFF_MIGRATED
	DomShutoffReasonSaved        DomainShutoffReason = C.VIR_DOMAIN_SHUTOFF_SAVED
	DomShutoffReasonFailed       DomainShutoffReason = C.VIR_DOMAIN_SHUTOFF_FAILED
	DomShutoffReasonFromSnapshot DomainShutoffReason = C.VIR_DOMAIN_SHUTOFF_FROM_SNAPSHOT
)

// DomainCrashedReason describes the reason which led a domain to be on "DomStateCrashed".
type DomainCrashedReason uint32

// Possible values for DomainCrashedReason.
const (
	DomCrashedReasonUnknown  DomainCrashedReason = C.VIR_DOMAIN_CRASHED_UNKNOWN
	DomCrashedReasonPanicked DomainCrashedReason = C.VIR_DOMAIN_CRASHED_PANICKED
)

// DomainPMSuspendedReason describes the reason which led a domain to be on "DomStatePMSuspended".
type DomainPMSuspendedReason uint32

// Possible values for DomainPMSuspendedReason.
const (
	DomPMSuspendedReasonUnknown DomainPMSuspendedReason = C.VIR_DOMAIN_PMSUSPENDED_UNKNOWN
)

// DomainDumpFlag defines how a domain coredump should be taken.
type DomainDumpFlag uint32

// Possible values for DomainDumpFlag.
const (
	DomDumpDefault     DomainDumpFlag = 0
	DomDumpCrash       DomainDumpFlag = C.VIR_DUMP_CRASH
	DomDumpLive        DomainDumpFlag = C.VIR_DUMP_LIVE
	DomDumpBypassCache DomainDumpFlag = C.VIR_DUMP_BYPASS_CACHE
	DomDumpReset       DomainDumpFlag = C.VIR_DUMP_RESET
	DomDumpMemoryOnly  DomainDumpFlag = C.VIR_DUMP_MEMORY_ONLY
)

// DomainDumpFormat defines the format of a domain core dump.
type DomainDumpFormat uint32

// Possible values for DomainDumpFormat.
const (
	DomDumpFormatRaw         DomainDumpFormat = C.VIR_DOMAIN_CORE_DUMP_FORMAT_RAW
	DomDumpFormatKdumpZlib   DomainDumpFormat = C.VIR_DOMAIN_CORE_DUMP_FORMAT_KDUMP_ZLIB
	DomDumpFormatKdumpLzo    DomainDumpFormat = C.VIR_DOMAIN_CORE_DUMP_FORMAT_KDUMP_LZO
	DomDumpFormatKdumpSnappy DomainDumpFormat = C.VIR_DOMAIN_CORE_DUMP_FORMAT_KDUMP_SNAPPY
)

// DomainVCPUsFlag defines how a domain VCPUs count should be handled.
type DomainVCPUsFlag uint32

// Possible values for DomainVCPUsFlag.
const (
	DomVCPusConfig  DomainVCPUsFlag = C.VIR_DOMAIN_VCPU_CONFIG
	DomVCPUsCurrent DomainVCPUsFlag = C.VIR_DOMAIN_VCPU_CURRENT
	DomVCPUsLive    DomainVCPUsFlag = C.VIR_DOMAIN_VCPU_LIVE
	DomVCPUsMaximum DomainVCPUsFlag = C.VIR_DOMAIN_VCPU_MAXIMUM
	DomVCPUsGuest   DomainVCPUsFlag = C.VIR_DOMAIN_VCPU_GUEST
)

// DomainSaveFlag defines how a domain should be saved/restored.
type DomainSaveFlag uint32

// Possible values for DomainSaveFlag.
const (
	DomSaveDefault     DomainSaveFlag = 0
	DomSaveBypassCache DomainSaveFlag = C.VIR_DOMAIN_SAVE_BYPASS_CACHE
	DomSaveRunning     DomainSaveFlag = C.VIR_DOMAIN_SAVE_RUNNING
	DomSavePaused      DomainSaveFlag = C.VIR_DOMAIN_SAVE_PAUSED
)

// DomainDeviceModifyFlag defines how a domain device should be attached/detached/modified.
type DomainDeviceModifyFlag uint32

// Possible values for DomainDeviceModifyFlag.
const (
	DomDeviceModifyConfig  DomainDeviceModifyFlag = C.VIR_DOMAIN_DEVICE_MODIFY_CONFIG
	DomDeviceModifyCurrent DomainDeviceModifyFlag = C.VIR_DOMAIN_DEVICE_MODIFY_CURRENT
	DomDeviceModifyLive    DomainDeviceModifyFlag = C.VIR_DOMAIN_DEVICE_MODIFY_LIVE
	DomDeviceModifyForce   DomainDeviceModifyFlag = C.VIR_DOMAIN_DEVICE_MODIFY_FORCE
)

// DomainMemoryModifyFlag controls how the domain memory should be modified.
type DomainMemoryModifyFlag uint32

// Possible values for DomainMemoryModifyFlag.
const (
	DomMemoryConfig  DomainMemoryModifyFlag = C.VIR_DOMAIN_MEM_CONFIG
	DomMemoryCurrent DomainMemoryModifyFlag = C.VIR_DOMAIN_MEM_CURRENT
	DomMemoryLive    DomainMemoryModifyFlag = C.VIR_DOMAIN_MEM_LIVE
	DomMemoryMaximum DomainMemoryModifyFlag = C.VIR_DOMAIN_MEM_MAXIMUM
)

// DomainKeycodeSet defines a code set of keycodes.
type DomainKeycodeSet uint32

// Possible values for DomainKeycodeSet.
const (
	DomKeycodeSetLinux  DomainKeycodeSet = C.VIR_KEYCODE_SET_LINUX
	DomKeycodeSetXT     DomainKeycodeSet = C.VIR_KEYCODE_SET_XT
	DomKeycodeSetATSet1 DomainKeycodeSet = C.VIR_KEYCODE_SET_ATSET1
	DomKeycodeSetATSet2 DomainKeycodeSet = C.VIR_KEYCODE_SET_ATSET2
	DomKeycodeSetATSet3 DomainKeycodeSet = C.VIR_KEYCODE_SET_ATSET3
	DomKeycodeSetOSX    DomainKeycodeSet = C.VIR_KEYCODE_SET_OSX
	DomKeycodeSetXTKbd  DomainKeycodeSet = C.VIR_KEYCODE_SET_XT_KBD
	DomKeycodeSetUSB    DomainKeycodeSet = C.VIR_KEYCODE_SET_USB
	DomKeycodeSetWin32  DomainKeycodeSet = C.VIR_KEYCODE_SET_WIN32
	DomKeycodeSetRFB    DomainKeycodeSet = C.VIR_KEYCODE_SET_RFB
)

// DomainProcessSignal defines the valid signals which can be sent to a domain.
type DomainProcessSignal uint32

// Possible values for DomainProcessSignal.
const (
	DomSIGNOP    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_NOP
	DomSIGHUP    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_HUP
	DomSIGINT    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_INT
	DomSIGQUIT   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_QUIT
	DomSIGILL    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_ILL
	DomSIGTRAP   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_TRAP
	DomSIGABRT   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_ABRT
	DomSIGBUS    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_BUS
	DomSIGFPE    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_FPE
	DomSIGKILL   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_KILL
	DomSIGUSR1   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_USR1
	DomSIGSEGV   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_SEGV
	DomSIGUSR2   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_USR2
	DomSIGPIPE   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_PIPE
	DomSIGALRM   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_ALRM
	DomSIGTERM   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_TERM
	DomSIGSTKFLT DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_STKFLT
	DomSIGCHLD   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_CHLD
	DomSIGCONT   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_CONT
	DomSIGSTOP   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_STOP
	DomSIGTSTP   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_TSTP
	DomSIGTTIN   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_TTIN
	DomSIGTTOU   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_TTOU
	DomSIGURG    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_URG
	DomSIGXCPU   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_XCPU
	DomSIGXFSZ   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_XFSZ
	DomSIGVTALRM DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_VTALRM
	DomSIGPROF   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_PROF
	DomSIGWINCH  DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_WINCH
	DomSIGPOLL   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_POLL
	DomSIGPWR    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_PWR
	DomSIGSYS    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_SYS
	DomSIGRT0    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT0
	DomSIGRT1    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT1
	DomSIGRT2    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT2
	DomSIGRT3    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT3
	DomSIGRT4    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT4
	DomSIGRT5    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT5
	DomSIGRT6    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT6
	DomSIGRT7    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT7
	DomSIGRT8    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT8
	DomSIGRT9    DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT9
	DomSIGRT10   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT10
	DomSIGRT11   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT11
	DomSIGRT12   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT12
	DomSIGRT13   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT13
	DomSIGRT14   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT14
	DomSIGRT15   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT15
	DomSIGRT16   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT16
	DomSIGRT17   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT17
	DomSIGRT18   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT18
	DomSIGRT19   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT19
	DomSIGRT20   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT20
	DomSIGRT21   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT21
	DomSIGRT22   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT22
	DomSIGRT23   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT23
	DomSIGRT24   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT24
	DomSIGRT25   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT25
	DomSIGRT26   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT26
	DomSIGRT27   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT27
	DomSIGRT28   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT28
	DomSIGRT29   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT29
	DomSIGRT30   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT30
	DomSIGRT31   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT31
	DomSIGRT32   DomainProcessSignal = C.VIR_DOMAIN_PROCESS_SIGNAL_RT32
)

// Domain holds a libvirt domain. There are no exported fields.
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
func (dom Domain) Autostart() (bool, error) {
	var cAutostart C.int
	dom.log.Println("checking whether domain autostarts...")
	cRet := C.virDomainGetAutostart(dom.virDomain, &cAutostart)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return false, err
	}

	autostart := (int32(cAutostart) == 1)

	if autostart {
		dom.log.Println("domain autostarts")
	} else {
		dom.log.Println("domain does not autostart")
	}

	return autostart, nil
}

// HasCurrentSnapshot determines if the domain has a current snapshot.
func (dom Domain) HasCurrentSnapshot() (bool, error) {
	dom.log.Println("checking whether domain has current snapshot...")
	cRet := C.virDomainHasCurrentSnapshot(dom.virDomain, 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return false, err
	}

	hasCurrentSnapshot := (ret == 1)

	if hasCurrentSnapshot {
		dom.log.Println("domain has current snapshot")
	} else {
		dom.log.Println("domain does not have current snapshot")
	}

	return hasCurrentSnapshot, nil
}

// HasManagedSaveImage checks if a domain has a managed save image as created
// by ManagedSave(). Note that any running domain should not have such an
// image, as it should have been removed on restart.
func (dom Domain) HasManagedSaveImage() (bool, error) {
	dom.log.Println("checking whether domain has managed save...")
	cRet := C.virDomainHasManagedSaveImage(dom.virDomain, 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return false, err
	}

	hasManagedSave := (ret == 1)

	if hasManagedSave {
		dom.log.Println("domain has managed save")
	} else {
		dom.log.Println("domain does not have managed save")
	}

	return hasManagedSave, nil
}

// IsActive determines if the domain is currently running.
func (dom Domain) IsActive() (bool, error) {
	dom.log.Println("checking whether domain is active...")
	cRet := C.virDomainIsActive(dom.virDomain)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return false, err
	}

	active := (ret == 1)

	if active {
		dom.log.Println("domain is active")
	} else {
		dom.log.Println("domain is not active")
	}

	return active, nil
}

// IsPersistent determines if the domain has a persistent configuration which
// means it will still exist after shutting down
func (dom Domain) IsPersistent() (bool, error) {
	dom.log.Println("checking whether domain is persistent...")
	cRet := C.virDomainIsPersistent(dom.virDomain)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return false, err
	}

	persistent := (ret == 1)

	if persistent {
		dom.log.Println("domain is persistent")
	} else {
		dom.log.Println("domain is not persistent")
	}

	return persistent, nil
}

// IsUpdated determines if the domain has been updated.
func (dom Domain) IsUpdated() (bool, error) {
	dom.log.Println("checking whether domain is updated...")
	cRet := C.virDomainIsUpdated(dom.virDomain)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return false, err
	}

	updated := (ret == 1)

	if updated {
		dom.log.Println("domain is updated")
	} else {
		dom.log.Println("domain is not updated")
	}

	return updated, nil
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
func (dom Domain) Name() (string, error) {
	dom.log.Println("reading domain name...")
	cName := C.virDomainGetName(dom.virDomain)

	if cName == nil {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return "", err
	}

	name := C.GoString(cName)
	dom.log.Printf("domain name: %v\n", name)

	return name, nil
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
// Hypervisors may require the user to manually ensure proper permissions on the
// file named by "to".
// "dumpformat" controls which format the dump will have. Not all hypervisors
// are able to support all formats.
// If "flags" includes DomDumpCrash, then leave the guest shut off with a
// crashed state after the dump completes. If "flags" includes DomDumpLive, then
// make the core dump while continuing to allow the guest to run; otherwise, the
// guest is suspended during the dump. DomDumpReset flag forces reset of the
// guest after dump. The above three flags are mutually exclusive.
// Additionally, if "flags" includes DomDumpBypassCache, then libvirt will
// attempt to bypass the file system cache while creating the file, or fail if
// it cannot do so for the given system; this can allow less pressure on file
// system cache, but also risks slowing saves to NFS.
func (dom Domain) CoreDump(file string, format DomainDumpFormat, flags DomainDumpFlag) error {
	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))

	dom.log.Printf("dumping domain's core to file %v (format = %v, flags = %v)...", file, format, flags)
	cRet := C.virDomainCoreDumpWithFormat(dom.virDomain, cFile, C.uint(format), C.uint(flags))
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
func (dom Domain) SetMemory(memory uint64, flags DomainMemoryModifyFlag) error {
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

// ListSnapshots collects the list of domain snapshots for the given domain, and
// allocate an array to store those objects.
func (dom Domain) ListSnapshots(flags SnapshotListFlag) ([]Snapshot, error) {
	var cSnaps []C.virDomainSnapshotPtr
	snapsSH := (*reflect.SliceHeader)(unsafe.Pointer(&cSnaps))

	dom.log.Printf("reading domain snapshots (flags = %v)...\n", flags)
	cRet := C.virDomainListAllSnapshots(dom.virDomain, (**C.virDomainSnapshotPtr)(unsafe.Pointer(&snapsSH.Data)), C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return nil, err
	}
	defer C.free(unsafe.Pointer(snapsSH.Data))

	snapsSH.Cap = int(ret)
	snapsSH.Len = int(ret)

	snaps := make([]Snapshot, ret)

	for i := range snaps {
		snaps[i] = Snapshot{
			log:         dom.log,
			virSnapshot: cSnaps[i],
		}
	}

	dom.log.Printf("snapshots count: %v\n", len(snaps))

	return snaps, nil
}

// CreateSnapshot creates a new snapshot of a domain based on a snapshot XML.
func (dom Domain) CreateSnapshot(xml string, flags SnapshotCreateFlag) (Snapshot, error) {
	cXML := C.CString(xml)
	defer C.free(unsafe.Pointer(cXML))

	dom.log.Printf("creating domain snapshot (flags = %v)...\n", flags)
	cSnapshot := C.virDomainSnapshotCreateXML(dom.virDomain, cXML, C.uint(flags))
	if cSnapshot == nil {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return Snapshot{}, err
	}

	snap := Snapshot{
		log:         dom.log,
		virSnapshot: cSnapshot,
	}

	dom.log.Println("snapshot created")

	return snap, nil
}

// LookupSnapshotByName tries to lookup a domain snapshot based on its name.
func (dom Domain) LookupSnapshotByName(name string) (Snapshot, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	dom.log.Printf("looking up snapshot with name = %v...\n", name)
	cSnap := C.virDomainSnapshotLookupByName(dom.virDomain, cName, 0)
	if cSnap == nil {
		err := LastError()
		dom.log.Printf("an error occurred: %v\n", err)
		return Snapshot{}, err
	}

	snap := Snapshot{
		log:         dom.log,
		virSnapshot: cSnap,
	}

	dom.log.Println("snapshot found")

	return snap, nil
}
