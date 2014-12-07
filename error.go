package libvirt

// #include <libvirt/virterror.h>
import "C"
import (
	"fmt"
	"log"
)

// ErrorCode is the error code.
type ErrorCode uint32

// Possible values for ErrorCode.
const (
	ErrOK                    ErrorCode = C.VIR_ERR_OK
	ErrInternal              ErrorCode = C.VIR_ERR_INTERNAL_ERROR
	ErrNoMemory              ErrorCode = C.VIR_ERR_NO_MEMORY
	ErrNoSupport             ErrorCode = C.VIR_ERR_NO_SUPPORT
	ErrUnknownHost           ErrorCode = C.VIR_ERR_UNKNOWN_HOST
	ErrNoConnect             ErrorCode = C.VIR_ERR_NO_CONNECT
	ErrInvalidConn           ErrorCode = C.VIR_ERR_INVALID_CONN
	ErrInvalidDomain         ErrorCode = C.VIR_ERR_INVALID_DOMAIN
	ErrInvalidArg            ErrorCode = C.VIR_ERR_INVALID_ARG
	ErrOperationFailed       ErrorCode = C.VIR_ERR_OPERATION_FAILED
	ErrGetFailed             ErrorCode = C.VIR_ERR_GET_FAILED
	ErrPostFailed            ErrorCode = C.VIR_ERR_POST_FAILED
	ErrHTTP                  ErrorCode = C.VIR_ERR_HTTP_ERROR
	ErrSExprSerial           ErrorCode = C.VIR_ERR_SEXPR_SERIAL
	ErrNoXen                 ErrorCode = C.VIR_ERR_NO_XEN
	ErrXenCall               ErrorCode = C.VIR_ERR_XEN_CALL
	ErrOSType                ErrorCode = C.VIR_ERR_OS_TYPE
	ErrNoKernel              ErrorCode = C.VIR_ERR_NO_KERNEL
	ErrNoRoot                ErrorCode = C.VIR_ERR_NO_ROOT
	ErrNoSource              ErrorCode = C.VIR_ERR_NO_SOURCE
	ErrNoTarget              ErrorCode = C.VIR_ERR_NO_TARGET
	ErrNoName                ErrorCode = C.VIR_ERR_NO_NAME
	ErrNoOS                  ErrorCode = C.VIR_ERR_NO_OS
	ErrNoDevice              ErrorCode = C.VIR_ERR_NO_DEVICE
	ErrNoXenStore            ErrorCode = C.VIR_ERR_NO_XENSTORE
	ErrDriverFull            ErrorCode = C.VIR_ERR_DRIVER_FULL
	ErrCallFailed            ErrorCode = C.VIR_ERR_CALL_FAILED
	ErrXML                   ErrorCode = C.VIR_ERR_XML_ERROR
	ErrDomExist              ErrorCode = C.VIR_ERR_DOM_EXIST
	ErrOperationDenied       ErrorCode = C.VIR_ERR_OPERATION_DENIED
	ErrOpenFailed            ErrorCode = C.VIR_ERR_OPEN_FAILED
	ErrReadFailed            ErrorCode = C.VIR_ERR_READ_FAILED
	ErrParseFailed           ErrorCode = C.VIR_ERR_PARSE_FAILED
	ErrConfSyntax            ErrorCode = C.VIR_ERR_CONF_SYNTAX
	ErrWriteFailed           ErrorCode = C.VIR_ERR_WRITE_FAILED
	ErrXMLDetail             ErrorCode = C.VIR_ERR_XML_DETAIL
	ErrInvalidNetwork        ErrorCode = C.VIR_ERR_INVALID_NETWORK
	ErrNetworkExist          ErrorCode = C.VIR_ERR_NETWORK_EXIST
	ErrSystem                ErrorCode = C.VIR_ERR_SYSTEM_ERROR
	ErrRPC                   ErrorCode = C.VIR_ERR_RPC
	ErrGNUTLS                ErrorCode = C.VIR_ERR_GNUTLS_ERROR
	WarNoNetwork             ErrorCode = C.VIR_WAR_NO_NETWORK
	ErrNoDomain              ErrorCode = C.VIR_ERR_NO_DOMAIN
	ErrNoNetwork             ErrorCode = C.VIR_ERR_NO_NETWORK
	ErrInvalidMAC            ErrorCode = C.VIR_ERR_INVALID_MAC
	ErrAuthFailed            ErrorCode = C.VIR_ERR_AUTH_FAILED
	ErrInvalidStoragePool    ErrorCode = C.VIR_ERR_INVALID_STORAGE_POOL
	ErrInvalidStorageVol     ErrorCode = C.VIR_ERR_INVALID_STORAGE_VOL
	WarNoStorage             ErrorCode = C.VIR_WAR_NO_STORAGE
	ErrNoStoragePool         ErrorCode = C.VIR_ERR_NO_STORAGE_POOL
	ErrNoStorageVol          ErrorCode = C.VIR_ERR_NO_STORAGE_VOL
	WarNoNode                ErrorCode = C.VIR_WAR_NO_NODE
	ErrInvalidNodeDevice     ErrorCode = C.VIR_ERR_INVALID_NODE_DEVICE
	ErrNoNodeDevice          ErrorCode = C.VIR_ERR_NO_NODE_DEVICE
	ErrNoSecurityModel       ErrorCode = C.VIR_ERR_NO_SECURITY_MODEL
	ErrOperationInvalid      ErrorCode = C.VIR_ERR_OPERATION_INVALID
	WarNoInterface           ErrorCode = C.VIR_WAR_NO_INTERFACE
	ErrNoInterface           ErrorCode = C.VIR_ERR_NO_INTERFACE
	ErrInvalidInterface      ErrorCode = C.VIR_ERR_INVALID_INTERFACE
	ErrMultipleInterfaces    ErrorCode = C.VIR_ERR_MULTIPLE_INTERFACES
	WarNoNwFilter            ErrorCode = C.VIR_WAR_NO_NWFILTER
	ErrInvalidNwFilter       ErrorCode = C.VIR_ERR_INVALID_NWFILTER
	ErrNoNwFilter            ErrorCode = C.VIR_ERR_NO_NWFILTER
	ErrBuildFirewall         ErrorCode = C.VIR_ERR_BUILD_FIREWALL
	WarNoSecret              ErrorCode = C.VIR_WAR_NO_SECRET
	ErrInvalidSecret         ErrorCode = C.VIR_ERR_INVALID_SECRET
	ErrNoSecret              ErrorCode = C.VIR_ERR_NO_SECRET
	ErrConfigUnsupported     ErrorCode = C.VIR_ERR_CONFIG_UNSUPPORTED
	ErrOperationTimeout      ErrorCode = C.VIR_ERR_OPERATION_TIMEOUT
	ErrMigratePersistFailed  ErrorCode = C.VIR_ERR_MIGRATE_PERSIST_FAILED
	ErrHookScriptFailed      ErrorCode = C.VIR_ERR_HOOK_SCRIPT_FAILED
	ErrInvalidDomainSnapshot ErrorCode = C.VIR_ERR_INVALID_DOMAIN_SNAPSHOT
	ErrNoDomainSnapshot      ErrorCode = C.VIR_ERR_NO_DOMAIN_SNAPSHOT
	ErrInvalidStream         ErrorCode = C.VIR_ERR_INVALID_STREAM
	ErrArgumentUnsupported   ErrorCode = C.VIR_ERR_ARGUMENT_UNSUPPORTED
	ErrStorageProbeFailed    ErrorCode = C.VIR_ERR_STORAGE_PROBE_FAILED
	ErrStoragePoolBuilt      ErrorCode = C.VIR_ERR_STORAGE_POOL_BUILT
	ErrSnapshotRevertRisky   ErrorCode = C.VIR_ERR_SNAPSHOT_REVERT_RISKY
	ErrOperationAborted      ErrorCode = C.VIR_ERR_OPERATION_ABORTED
	ErrAuthCancelled         ErrorCode = C.VIR_ERR_AUTH_CANCELLED
	ErrNoDomainMetadata      ErrorCode = C.VIR_ERR_NO_DOMAIN_METADATA
	ErrMigrateUnsafe         ErrorCode = C.VIR_ERR_MIGRATE_UNSAFE
	ErrOverflow              ErrorCode = C.VIR_ERR_OVERFLOW
	ErrBlockCopyActive       ErrorCode = C.VIR_ERR_BLOCK_COPY_ACTIVE
	ErrOperationUnsupported  ErrorCode = C.VIR_ERR_OPERATION_UNSUPPORTED
	ErrSSH                   ErrorCode = C.VIR_ERR_SSH
	ErrAgentUnresponsive     ErrorCode = C.VIR_ERR_AGENT_UNRESPONSIVE
	ErrResourceBusy          ErrorCode = C.VIR_ERR_RESOURCE_BUSY
	ErrAccessDenied          ErrorCode = C.VIR_ERR_ACCESS_DENIED
	ErrDBusService           ErrorCode = C.VIR_ERR_DBUS_SERVICE
	ErrStorageVolExist       ErrorCode = C.VIR_ERR_STORAGE_VOL_EXIST
	ErrCPUIncompatible       ErrorCode = C.VIR_ERR_CPU_INCOMPATIBLE
)

// ErrorDomain describes what part of the library raised the error.
type ErrorDomain uint32

// Possible values for ErrorDomain.
const (
	ErrDomNone           ErrorDomain = C.VIR_FROM_NONE
	ErrDomXen            ErrorDomain = C.VIR_FROM_XEN
	ErrDomXend           ErrorDomain = C.VIR_FROM_XEND
	ErrDomXenStore       ErrorDomain = C.VIR_FROM_XENSTORE
	ErrDomSExpr          ErrorDomain = C.VIR_FROM_SEXPR
	ErrDomXML            ErrorDomain = C.VIR_FROM_XML
	ErrDomDom            ErrorDomain = C.VIR_FROM_DOM
	ErrDomRPC            ErrorDomain = C.VIR_FROM_RPC
	ErrDomProxy          ErrorDomain = C.VIR_FROM_PROXY
	ErrDomConf           ErrorDomain = C.VIR_FROM_CONF
	ErrDomQEMU           ErrorDomain = C.VIR_FROM_QEMU
	ErrDomNet            ErrorDomain = C.VIR_FROM_NET
	ErrDomTest           ErrorDomain = C.VIR_FROM_TEST
	ErrDomRemote         ErrorDomain = C.VIR_FROM_REMOTE
	ErrDomOpenVZ         ErrorDomain = C.VIR_FROM_OPENVZ
	ErrDomXenXM          ErrorDomain = C.VIR_FROM_XENXM
	ErrDomStatsLinux     ErrorDomain = C.VIR_FROM_STATS_LINUX
	ErrDomLXC            ErrorDomain = C.VIR_FROM_LXC
	ErrDomStorage        ErrorDomain = C.VIR_FROM_STORAGE
	ErrDomNetwork        ErrorDomain = C.VIR_FROM_NETWORK
	ErrDomDomain         ErrorDomain = C.VIR_FROM_DOMAIN
	ErrDomUML            ErrorDomain = C.VIR_FROM_UML
	ErrDomNodeDev        ErrorDomain = C.VIR_FROM_NODEDEV
	ErrDomXenInotify     ErrorDomain = C.VIR_FROM_XEN_INOTIFY
	ErrDomSecurity       ErrorDomain = C.VIR_FROM_SECURITY
	ErrDomVBox           ErrorDomain = C.VIR_FROM_VBOX
	ErrDomInterface      ErrorDomain = C.VIR_FROM_INTERFACE
	ErrDomONE            ErrorDomain = C.VIR_FROM_ONE
	ErrDomESX            ErrorDomain = C.VIR_FROM_ESX
	ErrDomPHYP           ErrorDomain = C.VIR_FROM_PHYP
	ErrDomSecret         ErrorDomain = C.VIR_FROM_SECRET
	ErrDomCPU            ErrorDomain = C.VIR_FROM_CPU
	ErrDomXenAPI         ErrorDomain = C.VIR_FROM_XENAPI
	ErrDomNwFilter       ErrorDomain = C.VIR_FROM_NWFILTER
	ErrDomHook           ErrorDomain = C.VIR_FROM_HOOK
	ErrDomDomainSnapshot ErrorDomain = C.VIR_FROM_DOMAIN_SNAPSHOT
	ErrDomAudit          ErrorDomain = C.VIR_FROM_AUDIT
	ErrDomSysinfo        ErrorDomain = C.VIR_FROM_SYSINFO
	ErrDomStreams        ErrorDomain = C.VIR_FROM_STREAMS
	ErrDomVMWare         ErrorDomain = C.VIR_FROM_VMWARE
	ErrDomEvent          ErrorDomain = C.VIR_FROM_EVENT
	ErrDomLibXL          ErrorDomain = C.VIR_FROM_LIBXL
	ErrDomLocking        ErrorDomain = C.VIR_FROM_LOCKING
	ErrDomHyperv         ErrorDomain = C.VIR_FROM_HYPERV
	ErrDomCapabilities   ErrorDomain = C.VIR_FROM_CAPABILITIES
	ErrDomURI            ErrorDomain = C.VIR_FROM_URI
	ErrDomAuth           ErrorDomain = C.VIR_FROM_AUTH
	ErrDomDBus           ErrorDomain = C.VIR_FROM_DBUS
	ErrDomParallels      ErrorDomain = C.VIR_FROM_PARALLELS
	ErrDomDevice         ErrorDomain = C.VIR_FROM_DEVICE
	ErrDomSSH            ErrorDomain = C.VIR_FROM_SSH
	ErrDomLockspace      ErrorDomain = C.VIR_FROM_LOCKSPACE
	ErrDomInitctl        ErrorDomain = C.VIR_FROM_INITCTL
	ErrDomIdentity       ErrorDomain = C.VIR_FROM_IDENTITY
	ErrDomCgroup         ErrorDomain = C.VIR_FROM_CGROUP
	ErrDomAccess         ErrorDomain = C.VIR_FROM_ACCESS
	ErrDomSystemd        ErrorDomain = C.VIR_FROM_SYSTEMD
	ErrDomBhyve          ErrorDomain = C.VIR_FROM_BHYVE
	ErrDomCrypto         ErrorDomain = C.VIR_FROM_CRYPTO
	ErrDomFirewall       ErrorDomain = C.VIR_FROM_FIREWALL
	ErrDomPolkit         ErrorDomain = C.VIR_FROM_POLKIT
)

// ErrorLevel specifies how consequent is the error.
type ErrorLevel uint32

// Possible values for ErrorLevel.
const (
	ErrLvlNone    ErrorLevel = C.VIR_ERR_NONE
	ErrLvlWarning ErrorLevel = C.VIR_ERR_WARNING
	ErrLvlError   ErrorLevel = C.VIR_ERR_ERROR
)

// Error is a wrapper for a native libvirt error.
type Error struct {
	Code             ErrorCode
	Domain           ErrorDomain
	Message          string
	Level            ErrorLevel
	Str1, Str2, Str3 string
	Int1, Int2       int32
}

func (err *Error) Error() string {
	return fmt.Sprintf("%s [error code = %d]", err.Message, err.Code)
}

// NewError creates an error based on a native libvirt error. If the libvirt
// error pointer is nil, returns nil.
func NewError(virError C.virErrorPtr) *Error {
	if virError == nil {
		return nil
	}

	return &Error{
		ErrorCode(virError.code),
		ErrorDomain(virError.domain),
		C.GoString(virError.message),
		ErrorLevel(virError.level),
		C.GoString(virError.str1),
		C.GoString(virError.str2),
		C.GoString(virError.str2),
		int32(virError.int1),
		int32(virError.int2),
	}
}

// LastError provides a pointer to the last error caught at the library level.
// The error object is kept in thread local storage, so separate threads can
// safely access this concurrently.
func LastError() *Error {
	cError := C.virGetLastError()
	if cError == nil {
		log.Println("LastError() did not return an error")
	}

	return NewError(cError)
}
