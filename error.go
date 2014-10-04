package libvirt

// #include <libvirt/virterror.h>
import "C"
import (
	"fmt"
	"log"
)

type ErrorCode uint32

const (
	ErrOK ErrorCode = iota
	ErrInternal
	ErrNoMemory
	ErrNoSupport
	ErrUnknownHost
	ErrNoConnect
	ErrInvalidConn
	ErrInvalidDomain
	ErrInvalidArg
	ErrOperationFailed
	ErrGetFailed
	ErrPostFailed
	ErrHTTP
	ErrSExprSerial
	ErrNoXen
	ErrXenCall
	ErrOSType
	ErrNoKernel
	ErrNoRoot
	ErrNoSource
	ErrNoTarget
	ErrNoName
	ErrNoOS
	ErrNoDevice
	ErrNoXenStore
	ErrDriverFull
	ErrCallFailed
	ErrXML
	ErrDomExist
	ErrOperationDenied
	ErrOpenFailed
	ErrReadFailed
	ErrParseFailed
	ErrConfSyntax
	ErrWriteFailed
	ErrXMLDetail
	ErrInvalidNetwork
	ErrNetworkExist
	ErrSystem
	ErrRPC
	ErrGNUTLS
	ErrVirWarNoNetwork
	ErrNoDomain
	ErrNoNetwork
	ErrInvalidMAC
	ErrAuthFailed
	ErrInvalidStoragePool
	ErrInvalidStorageVol
	ErrVirWarNoStorage
	ErrNoStoragePool
	ErrNoStorageVol
	ErrVirWarNoNode
	ErrInvalidNodeDevice
	ErrNoNodeDevice
	ErrNoSecurityLabel
	ErrOperationInvalid
	ErrVirWarNoInterface
	ErrNoInterface
	ErrInvalidInterface
	ErrMultipleInterfaces
	ErrVirWarNoNwFilter
	ErrInvalidNwFilter
	ErrNoNwFilter
	ErrBuildFirewall
	ErrVirWarNoSecret
	ErrInvalidSecret
	ErrNoSecret
	ErrConfigUnsupported
	ErrOperationTimeout
	ErrMigratePersistFailed
	ErrHookScriptFailed
	ErrInvalidDomainSnapshot
	ErrNoDomainSnapshot
	ErrInvalidStream
	ErrArgumentUnsupported
	ErrStorageProbeFailed
	ErrStoragePoolBuilt
	ErrSnapshotRevertRisky
	ErrOperationAborted
	ErrAuthCancelled
	ErrNoDomainMetadata
	ErrMigrateUnsafe
	ErrOverflow
	ErrBlockCopyActive
	ErrOperationUnsupported
	ErrSSH
	ErrAgentUnresponsive
	ErrResourceBusy
	ErrAccessDenied
	ErrDBusService
	ErrStorageVolExist
)

type ErrorDomain uint32

const (
	ErrDomNone ErrorDomain = iota
	ErrDomXen
	ErrDomXend
	ErrDomXenStore
	ErrDomSExpr
	ErrDomXML
	ErrDomDom
	ErrDomRPC
	ErrDomProxy
	ErrDomConf
	ErrDomQEMU
	ErrDomNet
	ErrDomTest
	ErrDomRemote
	ErrDomOpenVZ
	ErrDomXenXM
	ErrDomStatsLinux
	ErrDomLXC
	ErrDomStorage
	ErrDomNetwork
	ErrDomDomain
	ErrDomUML
	ErrDomNodeDev
	ErrDomXenInotify
	ErrDomSecurity
	ErrDomVBox
	ErrDomInterface
	ErrDomONE
	ErrDomESX
	ErrDomPHYP
	ErrDomSecret
	ErrDomCPU
	ErrDomXenAPI
	ErrDomNwFilter
	ErrDomHook
	ErrDomDomainSnapshot
	ErrDomAudit
	ErrDomSysinfo
	ErrDomStreams
	ErrDomVMWare
	ErrDomEvent
	ErrDomLibXL
	ErrDomLocking
	ErrDomHyperv
	ErrDomCapabilities
	ErrDomURI
	ErrDomAuth
	ErrDomDBus
	ErrDomParallels
	ErrDomDevice
	ErrDomSSH
	ErrDomLockspace
	ErrDomInitctl
	ErrDomIdentity
	ErrDomCgroup
	ErrDomAccess
	ErrDomSystemd
	ErrDomBhyve
)

type ErrorLevel uint32

const (
	ErrLvlNone ErrorLevel = (0 << iota)
	ErrLvlWarning
	ErrLvlError
)

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
