package libvirt

/*
#cgo pkg-config: libvirt
#include <stdlib.h>
#include <libvirt/libvirt.h>
#include <libvirt/virterror.h>

void emptyErrorFunc(void *userData, virErrorPtr error) {
    // do nothing
}
*/
import "C"
import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"unicode/utf8"
	"unsafe"
)

// Connection holds a libvirt connection. There are no exported fields.
type Connection struct {
	log        *log.Logger
	virConnect C.virConnectPtr
}

// ConnectionMode is the type of connection to the libvirt hypervisor.
type ConnectionMode uint

// Possible values for ConnectionMode.
const (
	ReadWrite ConnectionMode = iota
	ReadOnly
)

// DefaultURI is the URI chosen by libvirt to establish a default
// connection, based on the current environment.
// Check http://libvirt.org/uri.html for more details.
const DefaultURI = ""

// ErrInvalidConnectionMode is returned by "Open" when a value other than
// "ReadOnly" or "ReadWrite" is used.
var ErrInvalidConnectionMode = errors.New("invalid libvirt connection mode")

func init() {
	// Supress the native error output. There's no way to do this per
	// connection, so we have to do this globally.
	C.virSetErrorFunc(nil, C.virErrorFunc(unsafe.Pointer(C.emptyErrorFunc)))
}

// newLogger creates a logger object to be used across a libvirt
// connection. It prints the messages to the default error output.
func newLogger(output io.Writer) *log.Logger {
	return log.New(output, "libvirt-golang: ", log.LstdFlags)
}

// Open creates a new libvirt connection to the Hypervisor. The
// connection mode specifies whether the connection will be read-write
// or read-only. The URIs are documented at http://libvirt.org/uri.html.
func Open(uri string, mode ConnectionMode, logOutput io.Writer) (Connection, error) {
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))

	logger := newLogger(logOutput)

	if uri == DefaultURI {
		logger.Printf("opening connection (mode = %v) to the default URI...\n", mode)
	} else {
		logger.Printf("opening connection (mode = %v) to %v...\n", mode, uri)
	}

	var cConn C.virConnectPtr
	switch mode {
	case ReadWrite:
		cConn = C.virConnectOpen(cUri)
	case ReadOnly:
		cConn = C.virConnectOpenReadOnly(cUri)
	default:
		return Connection{}, ErrInvalidConnectionMode
	}

	if cConn == nil {
		err := LastError()
		logger.Printf("an error occurred: %v\n", err)
		return Connection{}, err
	}

	logger.Println("connection established")

	conn := Connection{
		log:        logger,
		virConnect: cConn,
	}

	return conn, nil
}

// OpenDefault creates a new read-write libvirt connection to the
// hypervisor using the default URI.
func OpenDefault() (Connection, error) {
	return Open(DefaultURI, ReadWrite, ioutil.Discard)
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
func (conn Connection) Close() (int32, error) {
	conn.log.Println("closing connection...")
	cRet := C.virConnectClose(conn.virConnect)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	conn.log.Printf("connection closed; remaining references: %v\n", ret)

	return ret, nil
}

// Version gets the version level of the Hypervisor running.
func (conn Connection) Version() (uint64, error) {
	var cVersion C.ulong
	conn.log.Println("reading hypervisor version...")
	cRet := C.virConnectGetVersion(conn.virConnect, &cVersion)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	version := uint64(cVersion)
	conn.log.Printf("hypervisor version: %v\n", version)

	return version, nil
}

// LibVersion provides the version of libvirt used by the daemon running on
// the host.
func (conn Connection) LibVersion() (uint64, error) {
	var cVersion C.ulong
	conn.log.Println("reading libvirt version...")
	cRet := C.virConnectGetLibVersion(conn.virConnect, &cVersion)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	version := uint64(cVersion)
	conn.log.Printf("libvirt version: %v\n", version)

	return version, nil
}

// IsAlive determines if the connection to the hypervisor is still alive.
// If an error occurs, the function will also return "false" and the error
// message will be written to the log.
func (conn Connection) IsAlive() (bool, error) {
	conn.log.Println("checking whether connection is alive...")
	cRet := C.virConnectIsAlive(conn.virConnect)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return false, err
	}

	alive := (ret == 1)

	if alive {
		conn.log.Println("connection is alive")
	} else {
		conn.log.Println("connection is not alive")
	}

	return alive, nil
}

// IsEncrypted determines if the connection to the hypervisor is encrypted.
// If an error occurs, the function will also return "false" and the error
// message will be written to the log.
func (conn Connection) IsEncrypted() (bool, error) {
	conn.log.Println("checking whether connection is encrypted...")
	cRet := C.virConnectIsEncrypted(conn.virConnect)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return false, err
	}

	encrypted := (ret == 1)

	if encrypted {
		conn.log.Println("connection is encrypted")
	} else {
		conn.log.Println("connection is not encrypted")
	}

	return encrypted, nil
}

// IsSecure determines if the connection to the hypervisor is secure.
// If an error occurs, the function will also return "false" and the error
// message will be written to the log.
func (conn Connection) IsSecure() (bool, error) {
	conn.log.Println("checking whether connection is secure...")
	cRet := C.virConnectIsSecure(conn.virConnect)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return false, err
	}

	secure := (ret == 1)

	if secure {
		conn.log.Println("connection is secure")
	} else {
		conn.log.Println("connection is not secure")
	}

	return secure, nil
}

// Capabilities provides capabilities of the hypervisor/driver.
func (conn Connection) Capabilities() (string, error) {
	conn.log.Println("reading connection capabilities...")
	cCap := C.virConnectGetCapabilities(conn.virConnect)
	if cCap == nil {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return "", err
	}
	defer C.free(unsafe.Pointer(cCap))

	cap := C.GoString(cCap)
	conn.log.Printf("capabilities XML length: %v runes\n", utf8.RuneCountInString(cap))

	return cap, nil
}

// Hostname returns a system hostname on which the hypervisor is running
// (based on the result of the gethostname system call, but possibly expanded
// to a fully-qualified domain name via getaddrinfo). If we are connected to a
// remote system, then this returns the hostname of the remote system.
func (conn Connection) Hostname() (string, error) {
	conn.log.Println("reading system hostname...")
	cHostname := C.virConnectGetHostname(conn.virConnect)
	if cHostname == nil {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return "", err
	}
	defer C.free(unsafe.Pointer(cHostname))

	hostname := C.GoString(cHostname)
	conn.log.Printf("system hostname: %v\n", hostname)

	return hostname, nil
}

// Sysinfo returns the XML description of the sysinfo details for the host on
// which the hypervisor is running, in the same format as the <sysinfo> element
// of a domain XML. This information is generally available only for
// hypervisors running with root privileges.
func (conn Connection) Sysinfo() (string, error) {
	conn.log.Println("reading system info...")
	cSysinfo := C.virConnectGetSysinfo(conn.virConnect, 0)
	if cSysinfo == nil {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return "", err
	}
	defer C.free(unsafe.Pointer(cSysinfo))

	sysinfo := C.GoString(cSysinfo)
	conn.log.Printf("system info XML length: %v runes\n", utf8.RuneCountInString(sysinfo))
	return sysinfo, nil
}

// Type gets the name of the Hypervisor driver used. This is merely the driver
// name; for example, both KVM and QEMU guests are serviced by the driver for
// the qemu:// URI, so a return of "QEMU" does not indicate whether KVM
// acceleration is present. For more details about the hypervisor, use
// Capabilities.
func (conn Connection) Type() (string, error) {
	conn.log.Println("reading hypervisor driver name...")
	cType := C.virConnectGetType(conn.virConnect)
	if cType == nil {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return "", err
	}

	typ := C.GoString(cType)
	conn.log.Printf("hypervisor driver name: %v\n", typ)

	return typ, nil
}

// URI returns the URI (name) of the hypervisor connection. Normally this is
// the same as or similar to the string passed to the Open/OpenReadOnly call,
// but the driver may make the URI canonical. If uri == "" was passed to Open,
// then the driver will return a non-NULL URI which can be used to connect tos
// the same hypervisor later.
func (conn Connection) URI() (string, error) {
	conn.log.Println("reading connection URI...")
	cURI := C.virConnectGetURI(conn.virConnect)
	if cURI == nil {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return "", err
	}
	defer C.free(unsafe.Pointer(cURI))

	uri := C.GoString(cURI)
	conn.log.Printf("connection URI: %v\n", uri)
	return uri, nil
}

// Ref increments the reference count on the connection. For each additional
// call to this method, there shall be a corresponding call to Close to release
// the reference count, once the caller no longer needs the reference to
// this object.
func (conn Connection) Ref() error {
	conn.log.Println("incrementing connection's reference count...")
	cRet := C.virConnectRef(conn.virConnect)
	ret := int32(cRet)
	if ret == -1 {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return err
	}

	conn.log.Println("reference count incremented")

	return nil
}

// CPUModelNames gets the list of supported CPU models for a
// specific architecture.
func (conn Connection) CPUModelNames(arch string) ([]string, error) {
	cArch := C.CString(arch)
	defer C.free(unsafe.Pointer(cArch))

	var cModels []*C.char
	modelsSH := (*reflect.SliceHeader)(unsafe.Pointer(&cModels))

	conn.log.Printf("querying supported CPU models for %v...\n", arch)
	cRet := C.virConnectGetCPUModelNames(conn.virConnect, cArch, (***C.char)(unsafe.Pointer(&modelsSH.Data)), 0)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return nil, err
	}
	defer C.free(unsafe.Pointer(modelsSH.Data))

	modelsSH.Cap = int(ret)
	modelsSH.Len = int(ret)

	models := make([]string, ret)
	for i := range models {
		models[i] = C.GoString(cModels[i])
		defer C.free(unsafe.Pointer(cModels[i]))
	}

	conn.log.Printf("CPU models count: %v\n", ret)

	return models, nil
}

// MaxVCPUs provides the maximum number of virtual CPUs supported for a guest
// VM of a specific type. The 'type' parameter here corresponds to the 'type'
// attribute in the <domain> element of the XML
func (conn Connection) MaxVCPUs(typ string) (int32, error) {
	cTyp := C.CString(typ)
	defer C.free(unsafe.Pointer(cTyp))

	conn.log.Printf("querying maximum VCPUs count for %v...\n", typ)
	cRet := C.virConnectGetMaxVcpus(conn.virConnect, cTyp)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	conn.log.Printf("max VCPUs count: %v\n", ret)

	return ret, nil
}

// ListDomains collects a possibly-filtered list of all domains, and return an
// array of information for each.
func (conn Connection) ListDomains(flags DomainListFlag) ([]Domain, error) {
	var cDomains []C.virDomainPtr
	domainsSH := (*reflect.SliceHeader)(unsafe.Pointer(&cDomains))

	conn.log.Printf("reading domains (flags = %v)...\n", flags)
	cRet := C.virConnectListAllDomains(conn.virConnect, (**C.virDomainPtr)(unsafe.Pointer(&domainsSH.Data)), C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return nil, err
	}
	defer C.free(unsafe.Pointer(domainsSH.Data))

	domainsSH.Cap = int(ret)
	domainsSH.Len = int(ret)

	domains := make([]Domain, ret)
	for i := range domains {
		domains[i] = Domain{
			log:       conn.log,
			virDomain: cDomains[i],
		}
	}

	conn.log.Printf("domains count: %v\n", ret)

	return domains, nil
}

// CreateDomain launches a new guest domain, based on an XML description
// similar to the one returned by Domain.XML() This function may require
// privileged access to the hypervisor. The domain is not persistent, so its
// definition will disappear when it is destroyed, or if the host is restarted
// (see Domain.Define() to define persistent domains).
func (conn Connection) CreateDomain(xml string, flags DomainCreateFlag) (Domain, error) {
	cXML := C.CString(xml)
	defer C.free(unsafe.Pointer(cXML))

	conn.log.Printf("creating domain (flags = %v)...\n", flags)
	cDomain := C.virDomainCreateXML(conn.virConnect, cXML, C.uint(flags))
	if cDomain == nil {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return Domain{}, err
	}

	conn.log.Println("domain created")

	dom := Domain{
		log:       conn.log,
		virDomain: cDomain,
	}

	return dom, nil
}

// DefineDomain defines a domain, but does not start it. This definition is
// persistent, until explicitly undefined with Domain.Undefine(). A previous
// definition for this domain would be overridden if it already exists.
func (conn Connection) DefineDomain(xml string) (Domain, error) {
	cXML := C.CString(xml)
	defer C.free(unsafe.Pointer(cXML))

	conn.log.Println("defining domain...")
	cDomain := C.virDomainDefineXML(conn.virConnect, cXML)
	if cDomain == nil {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return Domain{}, err
	}

	conn.log.Println("domain defined")

	dom := Domain{
		log:       conn.log,
		virDomain: cDomain,
	}

	return dom, nil
}

// LookupDomainByID tries to find a domain based on the hypervisor ID number.
// Note that this won't work for inactive domains which have an ID of -1, in
// that case a lookup based on the Name or UUID need to be done instead.
func (conn Connection) LookupDomainByID(id uint32) (Domain, error) {
	conn.log.Printf("looking up domain with ID = %v...\n", id)
	cDomain := C.virDomainLookupByID(conn.virConnect, C.int(id))
	if cDomain == nil {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return Domain{}, err
	}

	conn.log.Println("domain found")

	dom := Domain{
		log:       conn.log,
		virDomain: cDomain,
	}

	return dom, nil
}

// LookupDomainByName tries to lookup a domain on the given hypervisor based on
// its name.
func (conn Connection) LookupDomainByName(name string) (Domain, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	conn.log.Printf("looking up domain with name = %v...\n", name)
	cDomain := C.virDomainLookupByName(conn.virConnect, cName)
	if cDomain == nil {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return Domain{}, err
	}

	conn.log.Println("domain found")

	dom := Domain{
		log:       conn.log,
		virDomain: cDomain,
	}

	return dom, nil
}

// LookupDomainByUUID tries to lookup a domain on the given hypervisor based on
// its UUID.
func (conn Connection) LookupDomainByUUID(uuid string) (Domain, error) {
	cUUID := C.CString(uuid)
	defer C.free(unsafe.Pointer(cUUID))

	conn.log.Printf("looking up domain with UUID = %v...\n", uuid)
	cDomain := C.virDomainLookupByUUIDString(conn.virConnect, cUUID)
	if cDomain == nil {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return Domain{}, err
	}

	conn.log.Println("domain found")

	dom := Domain{
		log:       conn.log,
		virDomain: cDomain,
	}

	return dom, nil
}

// RestoreDomain restores a domain saved to disk by Save().
func (conn Connection) RestoreDomain(from string, xml string, flags DomainSaveFlag) error {
	cFrom := C.CString(from)
	defer C.free(unsafe.Pointer(cFrom))

	var cXML *C.char
	if xml != "" {
		cXML = C.CString(xml)
		defer C.free(unsafe.Pointer(cXML))
	} else {
		cXML = nil
	}

	conn.log.Printf("restoring domain from file %v (flags = %v)...\n", from, flags)
	cRet := C.virDomainRestoreFlags(conn.virConnect, cFrom, cXML, C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return err
	}

	conn.log.Println("domain restored")

	return nil
}

// ListSecrets collects the list of secrets, and allocate an array to store those objects.
// Normally, all secrets are returned; however, "flags" can be used to filter
// the results for a smaller list of targeted secrets. The valid flags are
// divided into groups, where each group contains bits that describe mutually
// exclusive attributes of a secret, and where all bits within a group describe
// all possible secrets.
// The first group of "flags" is used to filter secrets by its storage location.
// Flag "SecListEphemeral" selects secrets that are kept only in memory. Flag
// "SecListNoEphemeral" selects secrets that are kept in persistent storage.
// The second group of "flags" is used to filter secrets by privacy. Flag
// "SecListPrivate" selects secrets that are never revealed to any caller of
// libvirt nor to any other node. Flag SecListNoPrivate selects
// non-private secrets.
func (conn Connection) ListSecrets(flags SecretListFlag) ([]Secret, error) {
	var cSecrets []C.virSecretPtr
	secretsSH := (*reflect.SliceHeader)(unsafe.Pointer(&cSecrets))

	conn.log.Printf("reading secrets (flags = %v)...\n", flags)
	cRet := C.virConnectListAllSecrets(conn.virConnect, (**C.virSecretPtr)(unsafe.Pointer(&secretsSH.Data)), C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return nil, err
	}
	defer C.free(unsafe.Pointer(secretsSH.Data))

	secretsSH.Cap = int(ret)
	secretsSH.Len = int(ret)

	secrets := make([]Secret, ret)
	for i := range secrets {
		secrets[i] = Secret{
			log:       conn.log,
			virSecret: cSecrets[i],
		}
	}

	conn.log.Printf("secrets count: %v\n", ret)

	return secrets, nil
}

// DefineSecret creates a new secret with an automatically chosen UUID, and
// initializes its attributes from "xml".
// If "xml" specifies a UUID, locates the specified secret and replaces all
// attributes of the secret specified by UUID by attributes specified in "xml"
// (any attributes not specified in "xml" are discarded).
// "Free" should be used to free the resources after the secret object is no
// longer needed.
func (conn Connection) DefineSecret(xml string) (Secret, error) {
	cXML := C.CString(xml)
	defer C.free(unsafe.Pointer(cXML))

	conn.log.Println("defining secret...")
	cSec := C.virSecretDefineXML(conn.virConnect, cXML, 0)

	if cSec == nil {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return Secret{}, err
	}

	conn.log.Println("secret defined")

	sec := Secret{
		log:       conn.log,
		virSecret: cSec,
	}

	return sec, nil
}

// LookupSecretByUUID tries to lookup a secret on the given hypervisor based on
// its UUID. Uses the printable string value to describe the UUID.
// "Free" should be used to free the resources after the secret object is no
// longer needed.
func (conn Connection) LookupSecretByUUID(uuid string) (Secret, error) {
	cUUID := C.CString(uuid)
	defer C.free(unsafe.Pointer(cUUID))

	cSecret := C.virSecretLookupByUUIDString(conn.virConnect, cUUID)

	if cSecret == nil {
		err := LastError()
		return Secret{}, err
	}

	secret := Secret{
		log:       conn.log,
		virSecret: cSecret,
	}

	return secret, nil
}

// LookupSecretByUsage tries to lookup a secret on the given hypervisor based on
// its usage. The usageID is unique within the set of secrets sharing the same
// usageType value.
// "Free" should be used to free the resources after the secret object is no
// longer needed.
func (conn Connection) LookupSecretByUsage(usageType SecretUsageType, usageID string) (Secret, error) {
	cUsageType := C.int(usageType)
	cUsageID := C.CString(usageID)
	defer C.free(unsafe.Pointer(cUsageID))

	cSecret := C.virSecretLookupByUsage(conn.virConnect, cUsageType, cUsageID)

	if cSecret == nil {
		err := LastError()
		return Secret{}, err
	}

	secret := Secret{
		log:       conn.log,
		virSecret: cSecret,
	}

	return secret, nil
}

// FindStoragePoolSources talks to a storage backend and attempts to
// auto-discover the set of available storage pool sources. e.g. For iSCSI this
// would be a set of iSCSI targets. For NFS this would be a list of exported
// paths. The "source" (optional for some storage pool types, e.g. local ones)
// is an instance of the storage pool's source element specifying where to look
// for the pools.
// "source" is not required for some types (e.g., those querying local storage
// resources only)
func (conn Connection) FindStoragePoolSources(typ string, source string) (string, error) {
	cType := C.CString(typ)
	defer C.free(unsafe.Pointer(cType))

	cSource := C.CString(source)
	defer C.free(unsafe.Pointer(cSource))

	conn.log.Printf("finding storage pool sources (type = %v)...\n", typ)
	cSources := C.virConnectFindStoragePoolSources(conn.virConnect, cType, cSource, 0)

	if cSources == nil {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return "", err
	}
	defer C.free(unsafe.Pointer(cSources))

	sources := C.GoString(cSources)
	conn.log.Printf("sources XML length: %v runes\n", utf8.RuneCountInString(sources))

	return sources, nil
}

// ListStoragePools collects the list of storage pools, and allocates an array
// to store those objects.
// Normally, all storage pools are returned; however, "flags" can be used to
// filter the results for a smaller list of targeted pools. The valid flags are
// divided into groups, where each group contains bits that describe mutually
// exclusive attributes of a pool, and where all bits within a group describe
// all possible pools.
// The first group of "flags" is PoolListActive (online) and PoolListInactive
// (offline) to filter the pools by state.
// The second group of "flags" is PoolListPersistent (defined) and
// PoolListTransient (running but not defined), to filter the pools by whether
// they have persistent config or not.
// The third group of "flags" is PoolListAutostart and PoolListNoAutostart, to
// filter the pools by whether they are marked as autostart or not.
// The last group of "flags" is provided to filter the pools by the types, the
// flags include: PoolListDir, PoolListFS, PoolListNetFS, PoolListLogical,
// PoolListDisk, PoolListISCSI, PoolListSCSI, PoolListMPath, PoolListRBD,
// PoolListSheepdog.
func (conn Connection) ListStoragePools(flags StoragePoolListFlag) ([]StoragePool, error) {
	var cStoragePools []C.virStoragePoolPtr
	cStoragePoolsSH := (*reflect.SliceHeader)(unsafe.Pointer(&cStoragePools))

	conn.log.Printf("reading storage pools (flags = %v)...\n", flags)
	cRet := C.virConnectListAllStoragePools(conn.virConnect, (**C.virStoragePoolPtr)(unsafe.Pointer(&cStoragePoolsSH.Data)), C.uint(flags))
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		conn.log.Printf("an error occurred: %v\n", err)
		return nil, err
	}
	defer C.free(unsafe.Pointer(cStoragePoolsSH.Data))

	cStoragePoolsSH.Cap = int(ret)
	cStoragePoolsSH.Len = int(ret)

	storagePools := make([]StoragePool, ret)
	for i, cPool := range cStoragePools {
		storagePools[i] = StoragePool{
			log:            conn.log,
			virStoragePool: cPool,
		}
	}

	conn.log.Printf("pools count: %v\n", ret)

	return storagePools, nil
}
