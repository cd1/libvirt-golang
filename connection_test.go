package libvirt

import (
	"bytes"
	"testing"

	"github.com/cd1/utils-golang"
)

func TestConnectionOpenClose(t *testing.T) {
	if _, err := Open(utils.RandomString(), ReadWrite, testLogOutput); err == nil {
		t.Error("an error was not returned when connecting to a bad URI")
	}

	conn, err := Open(testConnectionURI, ReadWrite, testLogOutput)
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
	env := newTestEnvironment(t)

	if err := env.conn.Ref(); err != nil {
		t.Fatal(err)
	}

	ref, err := env.conn.Close()
	if err != nil {
		t.Error(err)
	}

	if ref != 1 {
		t.Errorf("unexpected connection reference count after closing connection for the first time; got=%v, want=1", ref)
	}

	ref, err = env.conn.Close()
	if err != nil {
		t.Error("could not close the connection for the second time after calling Ref")
	}

	if ref != 0 {
		t.Errorf("unexpected connection reference count after closing connection for the second time; got=%v, want=0", ref)
	}
}

func TestConnectionReadOnly(t *testing.T) {
	conn, err := Open(testConnectionURI, ReadOnly, testLogOutput)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	var xml bytes.Buffer
	data := newTestDomainData()

	if err = testDomainTmpl.Execute(&xml, data); err != nil {
		t.Fatal(err)
	}

	if _, err := conn.DefineDomain(xml.String()); err == nil {
		t.Error("a readonly libvirt connection should not allow defining domains")
	}

	if _, err := conn.CreateDomain(xml.String(), DomCreateDefault); err == nil {
		t.Error("a readonly libvirt connection should not allow creating domains")
	}
}

func TestConnectionInit(t *testing.T) {
	env := newTestEnvironment(t)
	defer env.cleanUp()

	if !env.conn.IsAlive() {
		t.Error("the libvirt connection was opened but it is not alive")
	}

	if env.conn.IsEncrypted() {
		t.Error("the libvirt connection is encrypted (but it should not)")
	}

	if !env.conn.IsSecure() {
		t.Error("the libvirt connection is not secure (but it should)")
	}

	if _, err := env.conn.Version(); err != nil {
		t.Error(err)
	}

	if _, err := env.conn.LibVersion(); err != nil {
		t.Error(err)
	}

	cap, err := env.conn.Capabilities()
	if err != nil {
		t.Error(err)
	}

	if len(cap) == 0 {
		t.Error("libvirt capabilities should not be empty")
	}

	hostname, err := env.conn.Hostname()
	if err != nil {
		t.Error(err)
	}

	if len(hostname) == 0 {
		t.Error("libvirt hostname should not be empty")
	}

	_, err = env.conn.Sysinfo()
	if err != nil {
		t.Error(err)
	}

	_, err = env.conn.Type()
	if err != nil {
		t.Error(err)
	}

	uri, err := env.conn.URI()
	if err != nil {
		t.Error(err)
	}

	if uri != testConnectionURI {
		t.Errorf("libvirt URI should be the same used to open the connection; got=%v, want=%v", uri, testConnectionURI)
	}

	if _, err = env.conn.CPUModelNames(utils.RandomString()); err == nil {
		t.Error("an error was not returned when getting CPU model names from invalid arch")
	}

	models, err := env.conn.CPUModelNames("x86_64")
	if err != nil {
		t.Error(err)
	}

	if len(models) == 0 {
		t.Error("libvirt CPU model names should not be empty")
	}

	if _, err = env.conn.MaxVCPUs(utils.RandomString()); err == nil {
		t.Error("an error was not returned when getting maximum VCPUs from invalid type")
	}

	vcpus, err := env.conn.MaxVCPUs("kvm")
	if err != nil {
		t.Fatal(err)
	}

	if vcpus < 0 {
		t.Error("libvirt maximum VCPU should be a positive number")
	}
}

func TestConnectionListDomains(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	domains, err := env.conn.ListDomains(DomListAll)
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
	env := newTestEnvironment(t)
	defer env.cleanUp()

	if _, err := env.conn.CreateDomain("", DomCreateDefault); err == nil {
		t.Error("an error was not returned when creating a domain with empty XML descriptor")
	}

	var xml bytes.Buffer
	data := newTestDomainData()

	if err := testDomainTmpl.Execute(&xml, data); err != nil {
		t.Fatal(err)
	}

	if _, err := env.conn.CreateDomain(xml.String(), DomainCreateFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid create flag")
	}

	dom, err := env.conn.CreateDomain(xml.String(), DomCreateDefault)
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
	env := newTestEnvironment(t)
	defer env.cleanUp()

	if _, err := env.conn.DefineDomain(""); err == nil {
		t.Error("an error was not returned when defining a domain with empty XML descriptor")
	}

	var xml bytes.Buffer
	data := newTestDomainData()

	if err := testDomainTmpl.Execute(&xml, data); err != nil {
		t.Fatal(err)
	}

	dom, err := env.conn.DefineDomain(xml.String())
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
	// TODO: if a domain is created with "<Domain>.Create" after
	// "<Connection>.Define", it doesn't see to get an ID. as a workaround, we
	// create it directly with "<Connection>.CreateDomain" because then it works.
	env := newTestEnvironment(t)
	defer env.cleanUp()

	data := newTestDomainData()

	var xml bytes.Buffer

	if err := testDomainTmpl.Execute(&xml, data); err != nil {
		t.Fatal(err)
	}

	dom, err := env.conn.CreateDomain(xml.String(), DomCreateAutodestroy)
	if err != nil {
		t.Fatal(err)
	}

	// ByID
	if _, err := env.conn.LookupDomainByID(99); err == nil {
		t.Error("an error was not returned when looking up a non-existing domain ID")
	}

	expectedID, err := dom.ID()
	if err != nil {
		t.Error(err)
	}

	dom, err = env.conn.LookupDomainByID(expectedID)
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
	if _, err = env.conn.LookupDomainByName(utils.RandomString()); err == nil {
		t.Error("an error was not returned when looking up a non-existing domain name")
	}

	dom, err = env.conn.LookupDomainByName(data.Name)
	if err != nil {
		t.Error(err)
	}
	defer dom.Free()

	if name := dom.Name(); name != data.Name {
		t.Errorf("looked up domain with unexpected name; got=%v, want=%v", name, data.Name)
	}

	// ByUUID
	if _, err := env.conn.LookupDomainByUUID(utils.RandomString()); err == nil {
		t.Error("an error was not returned when looking up a non-existing domain UUID")
	}

	dom, err = env.conn.LookupDomainByUUID(data.UUID)
	if err != nil {
		t.Error(err)
	}
	defer dom.Free()

	uuid, err := dom.UUID()
	if err != nil {
		t.Error(err)
	}

	if uuid != data.UUID {
		t.Errorf("looked up domain with unexpected UUID; got=%v, want=%v", uuid, data.UUID)
	}
}

func BenchmarkConnectionOpenRW(b *testing.B) {
	for n := 0; n < b.N; n++ {
		conn, err := Open(testConnectionURI, ReadWrite, testLogOutput)
		if err != nil {
			b.Error(err)
		}

		if _, err := conn.Close(); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkConnectionCreateDomain(b *testing.B) {
	conn, err := Open(testConnectionURI, ReadWrite, testLogOutput)
	if err != nil {
		b.Fatal(err)
	}

	var xml bytes.Buffer
	data := newTestDomainData()

	if err = testDomainTmpl.Execute(&xml, data); err != nil {
		b.Fatal(err)
	}
	xmlStr := xml.String()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dom, err := conn.CreateDomain(xmlStr, DomCreateDefault)
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
	conn, err := Open(testConnectionURI, ReadWrite, testLogOutput)
	if err != nil {
		b.Fatal(err)
	}

	var xml bytes.Buffer
	data := newTestDomainData()

	if err = testDomainTmpl.Execute(&xml, data); err != nil {
		b.Fatal(err)
	}
	xmlStr := xml.String()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dom, err := conn.DefineDomain(xmlStr)
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
