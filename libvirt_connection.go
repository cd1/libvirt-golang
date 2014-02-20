package libvirt

// #cgo pkg-config: libvirt
// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"
import (
	"log"
	"unsafe"
)

// Connection holds a libvirt connection. There are no exported fields.
type Connection struct {
	virConnect C.virConnectPtr
}

// Open creates a new libvirt connection to the Hypervisor. The URIs are
// documented at http://libvirt.org/uri.html.
func Open(uri string) (Connection, *Error) {
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))

	cConn := C.virConnectOpen(cUri)
	if cConn == nil {
		return Connection{}, lastError()
	}

	return Connection{cConn}, nil
}

// OpenReadOnly creates a restricted libvirt connection. The URIs are
// documented at http://libvirt.org/uri.html.
func OpenReadOnly(uri string) (Connection, *Error) {
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))

	cConn := C.virConnectOpenReadOnly(cUri)
	if cConn == nil {
		return Connection{}, lastError()
	}

	return Connection{cConn}, nil
}

// Close closes the connection to the Hypervisor. Connections are reference
// counted; the count is explicitly increased by the initial open (Open,
// OpenAuth, and the like) as well as Ref (not implemented yet); it is also
// temporarily increased by other API that depend on the connection remaining
// alive. The open and every Ref call should have a matching Close, and all
// other references will be released after the corresponding operation
// completes.
// It returns a positive number if at least 1 reference remains on success. The
// returned value should not be assumed to be the total reference count. A
// return of 0 implies no references remain and the connection is closed and
// memory has been freed. It is possible for the last Close to return a
// positive value if some other object still has a temporary reference to the
// connection, but the application should not try to further use a connection
// after the Close that matches the initial open.
func (conn Connection) Close() (int, *Error) {
	cRet := C.virConnectClose(conn.virConnect)
	ret := int(cRet)

	if ret == -1 {
		return 0, lastError()
	}

	return ret, nil
}

// Version gets the version level of the Hypervisor running.
func (conn Connection) Version() (uint64, *Error) {
	var cVersion C.ulong
	cRet := C.virConnectGetVersion(conn.virConnect, &cVersion)
	ret := int(cRet)

	if ret == -1 {
		return 0, lastError()
	}

	return uint64(cVersion), nil
}

// LibVersion provides the version of libvirt used by the daemon running on
// the host.
func (conn Connection) LibVersion() (uint64, *Error) {
	var cVersion C.ulong
	cRet := C.virConnectGetLibVersion(conn.virConnect, &cVersion)
	ret := int(cRet)

	if ret == -1 {
		return 0, lastError()
	}

	return uint64(cVersion), nil
}

// IsAlive determines if the connection to the hypervisor is still alive.
// If an error occurs, the function will also return "false" and the error
// message will be written to the log.
func (conn Connection) IsAlive() bool {
	cRet := C.virConnectIsAlive(conn.virConnect)
	ret := int(cRet)

	if ret == 1 {
		return true
	}

	if ret == -1 {
		if err := lastError(); err != nil {
			log.Println(err)
		}
	}

	return false
}

// IsEncrypted determines if the connection to the hypervisor is encrypted.
// If an error occurs, the function will also return "false" and the error
// message will be written to the log.
func (conn Connection) IsEncrypted() bool {
	cRet := C.virConnectIsEncrypted(conn.virConnect)
	ret := int(cRet)

	if ret == 1 {
		return true
	}

	if ret == -1 {
		if err := lastError(); err != nil {
			log.Println(err)
		}
	}

	return false
}

// IsSecure determines if the connection to the hypervisor is secure.
// If an error occurs, the function will also return "false" and the error
// message will be written to the log.
func (conn Connection) IsSecure() bool {
	cRet := C.virConnectIsSecure(conn.virConnect)
	ret := int(cRet)

	if ret == 1 {
		return true
	}

	if ret == -1 {
		if err := lastError(); err != nil {
			log.Println(err)
		}
	}

	return false
}

// Capabilities provides capabilities of the hypervisor/driver.
func (conn Connection) Capabilities() (string, *Error) {
	cCap := C.virConnectGetCapabilities(conn.virConnect)
	if cCap == nil {
		return "", lastError()
	}
	defer C.free(unsafe.Pointer(cCap))

	return C.GoString(cCap), nil
}

// Hostname returns a system hostname on which the hypervisor is running
// (based on the result of the gethostname system call, but possibly expanded
// to a fully-qualified domain name via getaddrinfo). If we are connected to a
// remote system, then this returns the hostname of the remote system.
func (conn Connection) Hostname() (string, *Error) {
	cHostname := C.virConnectGetHostname(conn.virConnect)
	if cHostname == nil {
		return "", lastError()
	}
	defer C.free(unsafe.Pointer(cHostname))

	return C.GoString(cHostname), nil
}

// Sysinfo returns the XML description of the sysinfo details for the host on
// which the hypervisor is running, in the same format as the <sysinfo> element
// of a domain XML. This information is generally available only for
// hypervisors running with root privileges.
func (conn Connection) Sysinfo() (string, *Error) {
	cSysinfo := C.virConnectGetSysinfo(conn.virConnect, 0)
	if cSysinfo == nil {
		return "", lastError()
	}
	defer C.free(unsafe.Pointer(cSysinfo))

	return C.GoString(cSysinfo), nil
}

// Type gets the name of the Hypervisor driver used. This is merely the driver
// name; for example, both KVM and QEMU guests are serviced by the driver for
// the qemu:// URI, so a return of "QEMU" does not indicate whether KVM
// acceleration is present. For more details about the hypervisor, use
// Capabilities.
func (conn Connection) Type() (string, *Error) {
	cType := C.virConnectGetType(conn.virConnect)
	if cType == nil {
		return "", lastError()
	}

	return C.GoString(cType), nil
}

// Uri returns the URI (name) of the hypervisor connection. Normally this is
// the same as or similar to the string passed to the Open/OpenReadOnly call,
// but the driver may make the URI canonical. If uri == "" was passed to Open,
// then the driver will return a non-NULL URI which can be used to connect tos
// the same hypervisor later.
func (conn Connection) Uri() (string, *Error) {
	cUri := C.virConnectGetURI(conn.virConnect)
	if cUri == nil {
		return "", lastError()
	}
	defer C.free(unsafe.Pointer(cUri))

	return C.GoString(cUri), nil
}

// Ref increments the reference count on the connection. For each additional
// call to this method, there shall be a corresponding call to Close to release
// the reference count, once the caller no longer needs the reference to
// this object.
func (conn Connection) Ref() *Error {
	cRet := C.virConnectRef(conn.virConnect)
	ret := int(cRet)
	if ret == 0 {
		return nil
	} else {
		return lastError()
	}
}

// CpuModelNames gets the list of supported CPU models for a
// specific architecture.
func (conn Connection) CpuModelNames(arch string) ([]string, *Error) {
	var cModels **C.char
	cArch := C.CString(arch)
	defer C.free(unsafe.Pointer(cArch))

	cRet := C.virConnectGetCPUModelNames(conn.virConnect, cArch, &cModels, 0)
	ret := int(cRet)
	defer C.free(unsafe.Pointer(cModels))

	if ret == -1 {
		return nil, lastError()
	}

	cBackedModels := (*[1 << 30]*C.char)(unsafe.Pointer(cModels))[:ret:ret]
	models := make([]string, 0, ret)
	for i := 0; i < ret; i++ {
		models = append(models, C.GoString(cBackedModels[i]))
		defer C.free(unsafe.Pointer(cBackedModels[i]))
	}

	return models, nil
}

// MaxVcpus provides the maximum number of virtual CPUs supported for a guest
// VM of a specific type. The 'type' parameter here corresponds to the 'type'
// attribute in the <domain> element of the XML
func (conn Connection) MaxVcpus(typ string) (int, *Error) {
	cTyp := C.CString(typ)
	defer C.free(unsafe.Pointer(cTyp))

	cRet := C.virConnectGetMaxVcpus(conn.virConnect, cTyp)
	ret := int(cRet)

	if ret == -1 {
		return 0, lastError()
	}

	return ret, nil
}
