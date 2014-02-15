package libvirt

// #cgo pkg-config: libvirt
// #include <stdlib.h>
// #include <libvirt/libvirt.h>
// #include <libvirt/virterror.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// Connection holds a libvirt connection. There are no exported fields.
type Connection struct {
	virConnect C.virConnectPtr
}

// Open creates a new libvirt connection to the Hypervisor. The URIs are
// documented at http://libvirt.org/uri.html.
func Open(uri string) (Connection, error) {
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))

	cConn := C.virConnectOpen(cUri)
	if cConn == nil {
		return Connection{}, fmt.Errorf("libvirt connection to %s failed", uri)
	}

	return Connection{cConn}, nil
}

// OpenReadOnly creates a restricted libvirt connection. The URIs are
// documented at http://libvirt.org/uri.html.
func OpenReadOnly(uri string) (Connection, error) {
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))

	cConn := C.virConnectOpenReadOnly(cUri)
	if cConn == nil {
		return Connection{}, fmt.Errorf("libvirt connection to %s failed", uri)
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
func (conn Connection) Close() (int, error) {
	cRet := C.virConnectClose(conn.virConnect)
	ret := int(cRet)

	if ret == -1 {
		return 0, errors.New("failed to close libvirt connection")
	}

	return ret, nil
}

// Version gets the version level of the Hypervisor running.
func (conn Connection) Version() (uint64, error) {
	var cVersion C.ulong
	cRet := C.virConnectGetVersion(conn.virConnect, &cVersion)
	ret := int(cRet)

	if ret == -1 {
		return 0, errors.New("failed to get hypervisor version")
	}

	return uint64(cVersion), nil
}

// LibVersion provides the version of libvirt used by the daemon running on
// the host.
func (conn Connection) LibVersion() (uint64, error) {
	var cVersion C.ulong
	cRet := C.virConnectGetLibVersion(conn.virConnect, &cVersion)
	ret := int(cRet)

	if ret == -1 {
		return 0, errors.New("failed to get libvirt version")
	}

	return uint64(cVersion), nil
}
