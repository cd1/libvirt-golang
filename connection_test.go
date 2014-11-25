package libvirt

import (
	"io/ioutil"
	"testing"

	"github.com/cd1/utils-golang"
)

const (
	qemuSystemURI  = "qemu:///system"
	testDefaultURI = "test:///default"
)

var testLog = ioutil.Discard

func openTestConnection(t testing.TB) Connection {
	conn, err := Open(qemuSystemURI, ReadWrite, testLog)
	if err != nil {
		t.Fatal(err)
	}

	return conn
}

func TestConnectionOpenClose(t *testing.T) {
	if _, err := Open(utils.RandomString(), ReadWrite, testLog); err == nil {
		t.Error("an error was not returned when connecting to a bad URI")
	}

	conn, err := Open(qemuSystemURI, ReadWrite, testLog)
	if err != nil {
		t.Fatal(err)
	}

	_, err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestConnectionOpenDefault(t *testing.T) {
	conn, err := OpenDefault()
	if err != nil {
		t.Fatal(err)
	}

	if _, err = conn.Close(); err != nil {
		t.Error(err)
	}
}

func TestConnectionRef(t *testing.T) {
	conn := openTestConnection(t)

	if err := conn.Ref(); err != nil {
		t.Fatal(err)
	}

	ref, err := conn.Close()
	if err != nil {
		t.Error(err)
	}

	if ref != 1 {
		t.Errorf("unexpected connection reference count after closing connection for the first time; got=%v, want=1", ref)
	}

	ref, err = conn.Close()
	if err != nil {
		t.Error("could not close the connection for the second time after calling Ref")
	}

	if ref != 0 {
		t.Errorf("unexpected connection reference count after closing connection for the second time; got=%v, want=0", ref)
	}
}

func TestConnectionReadOnly(t *testing.T) {
	conn, err := Open(qemuSystemURI, ReadOnly, testLog)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if _, err := conn.DefineDomain(DomTestXML); err == nil {
		t.Error("a readonly libvirt connection should not allow defining domains")
	}

	if _, err := conn.CreateDomain(DomTestXML, DomCreateDefault); err == nil {
		t.Error("a readonly libvirt connection should not allow creating domains")
	}
}

func TestConnectionInit(t *testing.T) {
	conn := openTestConnection(t)
	defer conn.Close()

	if !conn.IsAlive() {
		t.Error("the libvirt connection was opened but it is not alive")
	}

	if conn.IsEncrypted() {
		t.Error("the libvirt connection is encrypted (but it should not)")
	}

	if !conn.IsSecure() {
		t.Error("the libvirt connection is not secure (but it should)")
	}

	if _, err := conn.Version(); err != nil {
		t.Error(err)
	}

	if _, err := conn.LibVersion(); err != nil {
		t.Error(err)
	}

	cap, err := conn.Capabilities()
	if err != nil {
		t.Error(err)
	}

	if len(cap) == 0 {
		t.Error("libvirt capabilities should not be empty")
	}

	hostname, err := conn.Hostname()
	if err != nil {
		t.Error(err)
	}

	if len(hostname) == 0 {
		t.Error("libvirt hostname should not be empty")
	}

	sysinfo, err := conn.Sysinfo()
	if err != nil {
		t.Error(err)
	}

	if len(sysinfo) == 0 {
		t.Error("libvirt sysinfo should not be empty")
	}

	typ, err := conn.Type()
	if err != nil {
		t.Error(err)
	}

	if len(typ) == 0 {
		t.Error("libvirt type should not be empty")
	}

	uri, err := conn.URI()
	if err != nil {
		t.Error(err)
	}

	if uri != qemuSystemURI {
		t.Errorf("libvirt URI should be the same used to open the connection; got=%v, want=%v", uri, qemuSystemURI)
	}

	if _, err = conn.CPUModelNames(utils.RandomString()); err == nil {
		t.Error("an error was not returned when getting CPU model names from invalid arch")
	}

	models, err := conn.CPUModelNames("x86_64")
	if err != nil {
		t.Error(err)
	}

	if len(models) == 0 {
		t.Error("libvirt CPU model names should not be empty")
	}

	if _, err = conn.MaxVCPUs(utils.RandomString()); err == nil {
		t.Error("an error was not returned when getting maximum VCPUs from invalid type")
	}

	vcpus, err := conn.MaxVCPUs("kvm")
	if err != nil {
		t.Fatal(err)
	}

	if vcpus < 0 {
		t.Error("libvirt maximum VCPU should be a positive number")
	}
}

func TestConnectionListDomains(t *testing.T) {
	conn := openTestConnection(t)
	defer conn.Close()

	domains, err := conn.ListDomains(DomListAll)
	if err != nil {
		t.Fatal(err)
	}

	for _, d := range domains {
		if err := d.Free(); err != nil {
			t.Error(err)
		}
	}
}

func TestConnectionCreateDestroyDomain(t *testing.T) {
	conn := openTestConnection(t)
	defer conn.Close()

	if _, err := conn.CreateDomain("", DomCreateDefault); err == nil {
		t.Error("an error was not returned when creating a domain with empty XML descriptor")
	}

	if _, err := conn.CreateDomain(DomTestXML, DomainCreateFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid create flag")
	}

	dom, err := conn.CreateDomain(DomTestXML, DomCreateDefault)
	if err != nil {
		t.Fatal(err)
	}
	defer dom.Free()

	if !dom.IsActive() {
		t.Error("domain should be active after being created")
	}

	if dom.IsPersistent() {
		t.Error("domain should not be persistent after being created")
	}

	if err := dom.Destroy(DomainDestroyFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid destroy flag")
	}

	if err := dom.Destroy(DomDestroyDefault); err != nil {
		t.Error(err)
	}

	if dom.IsActive() {
		t.Error("domain should not be active after being destroyed")
	}

	if dom.IsPersistent() {
		t.Error("domain should still not be persistent after being created and destroyed")
	}
}

func TestConnectionDefineUndefineDomain(t *testing.T) {
	conn := openTestConnection(t)
	defer conn.Close()

	if _, err := conn.DefineDomain(""); err == nil {
		t.Error("an error was not returned when defining a domain with empty XML descriptor")
	}

	dom, err := conn.DefineDomain(DomTestXML)
	if err != nil {
		t.Fatal(err)
	}
	defer dom.Free()

	if dom.IsActive() {
		t.Error("domain should not be active after being defined")
	}

	if !dom.IsPersistent() {
		t.Error("domain should be persistent after being defined")
	}

	if err := dom.Create(DomCreateDefault); err != nil {
		t.Fatal(err)
	}

	if !dom.IsActive() {
		t.Error("domain should be active after being defined and created")
	}

	if !dom.IsPersistent() {
		t.Error("domain should still be persistent after being defined and created")
	}

	if err := dom.Destroy(DomDestroyDefault); err != nil {
		t.Fatal(err)
	}

	if dom.IsActive() {
		t.Error("domain should not be active after being defined and destroyed")
	}

	if !dom.IsPersistent() {
		t.Error("domain should be persistent after being defined and destroyed")
	}

	if err := dom.Undefine(DomainUndefineFlag(99)); err == nil {
		t.Error("an error was not return when using an invalid undefine flag")
	}

	if err := dom.Undefine(DomUndefineDefault); err != nil {
		t.Error(err)
	}

	if dom.IsActive() {
		t.Error("domain should not be active after being undefined")
	}

	if dom.IsPersistent() {
		t.Error("domain should not be persistent after being undefined")
	}
}

func TestConnectionLookupDomain(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateAutodestroy)
	defer conn.Close()
	defer dom.Free()

	// ByID
	if _, err := conn.LookupDomainByID(99); err == nil {
		t.Error("an error was not returned when looking up a non-existing domain ID")
	}

	expectedID, err := dom.ID()
	if err != nil {
		t.Error(err)
	}

	dom, err = conn.LookupDomainByID(expectedID)
	if err != nil {
		t.Error(err)
	}
	defer dom.Free()

	id, err := dom.ID()
	if err != nil {
		t.Error(err)
	}

	if id != expectedID {
		t.Errorf("looked up domain with unexpected id; got=%v, want=%v", id, expectedID)
	}

	// ByName
	if _, err = conn.LookupDomainByName(utils.RandomString()); err == nil {
		t.Error("an error was not returned when looking up a non-existing domain name")
	}

	dom, err = conn.LookupDomainByName(DomTestName)
	if err != nil {
		t.Error(err)
	}
	defer dom.Free()

	if name := dom.Name(); name != DomTestName {
		t.Errorf("looked up domain with unexpected name; got=%v, want=%v", name, DomTestName)
	}

	// ByUUID
	if _, err := conn.LookupDomainByUUID(utils.RandomString()); err == nil {
		t.Error("an error was not returned when looking up a non-existing domain UUID")
	}

	dom, err = conn.LookupDomainByUUID(DomTestUUID)
	if err != nil {
		t.Error(err)
	}
	defer dom.Free()

	uuid, err := dom.UUID()
	if err != nil {
		t.Error(err)
	}

	if uuid != DomTestUUID {
		t.Errorf("looked up domain with unexpected UUID; got=%v, want=%v", uuid, DomTestUUID)
	}
}

func BenchmarkConnectionQEMU(b *testing.B) {
	for n := 0; n < b.N; n++ {
		conn, err := Open(qemuSystemURI, ReadWrite, testLog)
		if err != nil {
			b.Error(err)
		}

		if _, err := conn.Close(); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkConnectionTest(b *testing.B) {
	for n := 0; n < b.N; n++ {
		conn, err := Open(testDefaultURI, ReadWrite, testLog)
		if err != nil {
			b.Error(err)
		}

		if _, err := conn.Close(); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkConnectionCreateDomain(b *testing.B) {
	conn, err := Open(qemuSystemURI, ReadWrite, testLog)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dom, err := conn.CreateDomain(DomTestXML, DomCreateDefault)
		if err != nil {
			b.Error(err)
		}

		if err := dom.Destroy(DomDestroyDefault); err != nil {
			b.Error(err)
		}

		if err := dom.Free(); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()

	if _, err := conn.Close(); err != nil {
		b.Error(err)
	}
}

func BenchmarkConnectionDefineDomain(b *testing.B) {
	conn, err := Open(qemuSystemURI, ReadWrite, testLog)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dom, err := conn.DefineDomain(DomTestXML)
		if err != nil {
			b.Error(err)
		}

		if err := dom.Undefine(DomUndefineDefault); err != nil {
			b.Error(err)
		}

		if err := dom.Free(); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()

	if _, err := conn.Close(); err != nil {
		b.Error(err)
	}
}
