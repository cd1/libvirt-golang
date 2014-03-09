package libvirt

// #cgo pkg-config: libvirt
// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"
import (
	"log"
	"reflect"
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
		return Connection{}, LastError()
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
		return Connection{}, LastError()
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
		return 0, LastError()
	}

	return ret, nil
}

// Version gets the version level of the Hypervisor running.
func (conn Connection) Version() (uint64, *Error) {
	var cVersion C.ulong
	cRet := C.virConnectGetVersion(conn.virConnect, &cVersion)
	ret := int(cRet)

	if ret == -1 {
		return 0, LastError()
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
		return 0, LastError()
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
		if err := LastError(); err != nil {
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
		if err := LastError(); err != nil {
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
		if err := LastError(); err != nil {
			log.Println(err)
		}
	}

	return false
}

// Capabilities provides capabilities of the hypervisor/driver.
func (conn Connection) Capabilities() (string, *Error) {
	cCap := C.virConnectGetCapabilities(conn.virConnect)
	if cCap == nil {
		return "", LastError()
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
		return "", LastError()
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
		return "", LastError()
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
		return "", LastError()
	}

	return C.GoString(cType), nil
}

// URI returns the URI (name) of the hypervisor connection. Normally this is
// the same as or similar to the string passed to the Open/OpenReadOnly call,
// but the driver may make the URI canonical. If uri == "" was passed to Open,
// then the driver will return a non-NULL URI which can be used to connect tos
// the same hypervisor later.
func (conn Connection) URI() (string, *Error) {
	cURI := C.virConnectGetURI(conn.virConnect)
	if cURI == nil {
		return "", LastError()
	}
	defer C.free(unsafe.Pointer(cURI))

	return C.GoString(cURI), nil
}

// Ref increments the reference count on the connection. For each additional
// call to this method, there shall be a corresponding call to Close to release
// the reference count, once the caller no longer needs the reference to
// this object.
func (conn Connection) Ref() *Error {
	cRet := C.virConnectRef(conn.virConnect)
	ret := int(cRet)
	if ret == -1 {
		return LastError()
	}

	return nil
}

// CPUModelNames gets the list of supported CPU models for a
// specific architecture.
func (conn Connection) CPUModelNames(arch string) ([]string, *Error) {
	cArch := C.CString(arch)
	defer C.free(unsafe.Pointer(cArch))

	var cModels []*C.char
	modelsSH := (*reflect.SliceHeader)(unsafe.Pointer(&cModels))

	cRet := C.virConnectGetCPUModelNames(conn.virConnect, cArch, (***C.char)(unsafe.Pointer(&modelsSH.Data)), 0)
	ret := int(cRet)

	if ret == -1 {
		return nil, LastError()
	}
	defer C.free(unsafe.Pointer(modelsSH.Data))

	modelsSH.Cap = ret
	modelsSH.Len = ret

	models := make([]string, 0, ret)
	for i := 0; i < ret; i++ {
		models = append(models, C.GoString(cModels[i]))
		defer C.free(unsafe.Pointer(cModels[i]))
	}

	return models, nil
}

// MaxVCPUs provides the maximum number of virtual CPUs supported for a guest
// VM of a specific type. The 'type' parameter here corresponds to the 'type'
// attribute in the <domain> element of the XML
func (conn Connection) MaxVCPUs(typ string) (int, *Error) {
	cTyp := C.CString(typ)
	defer C.free(unsafe.Pointer(cTyp))

	cRet := C.virConnectGetMaxVcpus(conn.virConnect, cTyp)
	ret := int(cRet)

	if ret == -1 {
		return 0, LastError()
	}

	return ret, nil
}

// ListDomains collects a possibly-filtered list of all domains, and return an
// array of information for each.
func (conn Connection) ListDomains(flags DomainFlag) ([]Domain, *Error) {
	var cDomains []C.virDomainPtr
	domainsSH := (*reflect.SliceHeader)(unsafe.Pointer(&cDomains))

	cRet := C.virConnectListAllDomains(conn.virConnect, (**C.virDomainPtr)(unsafe.Pointer(&domainsSH.Data)), C.uint(flags))
	ret := int(cRet)

	if ret == -1 {
		return nil, LastError()
	}
	defer C.free(unsafe.Pointer(domainsSH.Data))

	domainsSH.Cap = ret
	domainsSH.Len = ret

	domains := make([]Domain, 0, ret)
	for i := 0; i < ret; i++ {
		domains = append(domains, Domain{cDomains[i]})
	}

	return domains, nil
}

// CreateDomain launches a new guest domain, based on an XML description
// similar to the one returned by Domain.XML() This function may require
// privileged access to the hypervisor. The domain is not persistent, so its
// definition will disappear when it is destroyed, or if the host is restarted
// (see Domain.Define() to define persistent domains).
func (conn Connection) CreateDomain(xml string, flags DomainCreateFlag) (Domain, *Error) {
	cXML := C.CString(xml)
	defer C.free(unsafe.Pointer(cXML))

	cDomain := C.virDomainCreateXML(conn.virConnect, cXML, C.uint(flags))
	if cDomain == nil {
		return Domain{}, LastError()
	}

	return Domain{cDomain}, nil
}

// DefineDomain defines a domain, but does not start it. This definition is
// persistent, until explicitly undefined with Domain.Undefine(). A previous
// definition for this domain would be overridden if it already exists.
func (conn Connection) DefineDomain(xml string) (Domain, *Error) {
	cXML := C.CString(xml)
	defer C.free(unsafe.Pointer(cXML))

	cDomain := C.virDomainDefineXML(conn.virConnect, cXML)
	if cDomain == nil {
		return Domain{}, LastError()
	}

	return Domain{cDomain}, nil
}

// LookupDomainByID tries to find a domain based on the hypervisor ID number.
// Note that this won't work for inactive domains which have an ID of -1, in
// that case a lookup based on the Name or UUID need to be done instead.
func (conn Connection) LookupDomainByID(id uint) (Domain, *Error) {
	cDomain := C.virDomainLookupByID(conn.virConnect, C.int(id))
	if cDomain == nil {
		return Domain{}, LastError()
	}

	return Domain{cDomain}, nil
}

// LookupDomainByName tries to lookup a domain on the given hypervisor based on
// its name.
func (conn Connection) LookupDomainByName(name string) (Domain, *Error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cDomain := C.virDomainLookupByName(conn.virConnect, cName)
	if cDomain == nil {
		return Domain{}, LastError()
	}

	return Domain{cDomain}, nil
}

// LookupDomainByUUID tries to lookup a domain on the given hypervisor based on
// its UUID.
func (conn Connection) LookupDomainByUUID(uuid string) (Domain, *Error) {
	cUUID := C.CString(uuid)
	defer C.free(unsafe.Pointer(cUUID))

	cDomain := C.virDomainLookupByUUIDString(conn.virConnect, cUUID)
	if cDomain == nil {
		return Domain{}, LastError()
	}

	return Domain{cDomain}, nil
}

// RestoreDomain restores a domain saved to disk by Save().
func (conn Connection) RestoreDomain(from string, xml string, flags DomainSaveFlag) *Error {
	cFrom := C.CString(from)
	defer C.free(unsafe.Pointer(cFrom))

	var cXML *C.char
	if xml != "" {
		cXML = C.CString(xml)
		defer C.free(unsafe.Pointer(cXML))
	} else {
		cXML = nil
	}

	cRet := C.virDomainRestoreFlags(conn.virConnect, cFrom, cXML, C.uint(flags))
	ret := int(cRet)

	if ret == -1 {
		return LastError()
	}

	return nil
}
