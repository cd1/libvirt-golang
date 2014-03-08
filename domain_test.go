package libvirt

import (
	"io/ioutil"
	"testing"
)

const (
	DomTestMetadataNamespace = "code.google.com/p/libvirt-golang"
	DomTestName              = "golang-test"
	DomTestOSType            = "hvm"
	DomTestUUID              = "9652e5cd-15f1-49ad-af73-63a502a9e2b8"
	DomTestXMLFile           = "res/dom-test.xml"
)

func createTestDomain(t testing.TB, flags DomainCreateFlag) (Domain, Connection) {
	conn := openTestConnection(t)

	xml, ioerr := ioutil.ReadFile(DomTestXMLFile)
	if ioerr != nil {
		t.Fatal(ioerr)
	}

	dom, err := conn.CreateDomain(string(xml), flags)
	if err != nil {
		t.Fatal(err)
	}

	return dom, conn
}

func defineTestDomain(t testing.TB) (Domain, Connection) {
	conn := openTestConnection(t)

	xml, ioerr := ioutil.ReadFile(DomTestXMLFile)
	if ioerr != nil {
		t.Fatal(ioerr)
	}

	dom, err := conn.DefineDomain(string(xml))
	if err != nil {
		t.Fatal(err)
	}

	return dom, conn
}

func TestDomainAutostart(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if dom.Autostart() {
		t.Error("test domain should not have autostart enabled")
	}
}

func TestDomainHasCurrentSnapshot(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if dom.HasCurrentSnapshot() {
		t.Error("test domain should not have current snapshot")
	}
}

func TestDomainHasManagedSaveImage(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if dom.HasManagedSaveImage() {
		t.Error("test domain should not have managed save image")
	}
}

func TestDomainIsActive(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if !dom.IsActive() {
		t.Error("test domain should be active")
	}
}

func TestDomainIsPersistent(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if dom.IsPersistent() {
		t.Error("test domain should be transient")
	}
}

func TestDomainIsUpdated(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if dom.IsUpdated() {
		t.Error("test domain should not have been updated")
	}
}

func TestDomainOSType(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	os, err := dom.OSType()
	if err != nil {
		t.Fatal(err)
	}

	if os != DomTestOSType {
		t.Errorf("wrong test domain OS type; got=%s, want=%", os, DomTestOSType)
	}
}

func TestDomainName(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	name := dom.Name()

	if name != DomTestName {
		t.Errorf("wrong test domain name; got=%s, want=%s", name, DomTestName)
	}
}

func TestDomainHostname(t *testing.T) {
	// Hostname is not supported by the "QEMU" driver
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if _, err := dom.Hostname(); err == nil {
		t.Error("Hostname should not be supported by the \"QEMU\" driver")
	}
}

func TestDomainID(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	id, err := dom.ID()
	if err != nil {
		t.Fatal(err)
	}

	if id < 0 {
		t.Error("domain ID should be a positive number")
	}
}

func TestDomainUUID(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	uuid, err := dom.UUID()
	if err != nil {
		t.Fatal(err)
	}

	if uuid != DomTestUUID {
		t.Errorf("wrong test domain UUID; got=%s, want=%s", uuid, DomTestUUID)
	}
}

func TestDomainXML(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if _, err := dom.XML(99); err == nil {
		t.Error("an error was not returned when using an invalid flag")
	}

	xml, err := dom.XML(DomXMLDefault)
	if err != nil {
		t.Fatal(err)
	}

	if len(xml) == 0 {
		t.Error("empty domain XML")
	}
}

func TestDomainMetadata(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if _, err := dom.Metadata(99, "", DomAffectCurrent); err == nil {
		t.Error("an error was not returned when using an invalid type")
	}

	if _, err := dom.Metadata(DomMetaElement, "xxx", DomAffectCurrent); err == nil {
		t.Error("an error was not returned when using a non-existing metadata tag")
	}

	if _, err := dom.Metadata(DomMetaElement, "", 99); err == nil {
		t.Error("an error was not returned when using an invalid impact config")
	}

	metadata, err := dom.Metadata(DomMetaElement, DomTestMetadataNamespace, DomAffectCurrent)
	if err != nil {
		t.Fatal(err)
	}

	if len(metadata) == 0 {
		t.Error("empty test domain metadata")
	}
}
