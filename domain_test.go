package libvirt

import (
	"testing"
)

const DomTestXMLFile = "res/dom-test.xml"

func openTestDomain(t testing.TB) (Domain, Connection) {
	conn := openTestConnection(t)
	domains, err := conn.ListDomains(DomAll)
	if err != nil {
		t.Fatal(err)
	}

	if len(domains) == 0 {
		t.Skip("there is no available domain to test")
	}

	// free every domain except the first one, which will be returned
	for i, d := range domains {
		if i > 0 {
			if err := d.Free(); err != nil {
				t.Error(err)
			}
		}
	}

	return domains[0], conn
}

func TestDomainAutostart(t *testing.T) {
	dom, conn := openTestDomain(t)
	defer conn.Close()
	defer dom.Free()

	_ = dom.Autostart()
}

func TestDomainHasCurrentSnapshot(t *testing.T) {
	dom, conn := openTestDomain(t)
	defer conn.Close()
	defer dom.Free()

	_ = dom.HasCurrentSnapshot()
}

func TestDomainHasManagedSaveImage(t *testing.T) {
	dom, conn := openTestDomain(t)
	defer conn.Close()
	defer dom.Free()

	_ = dom.HasManagedSaveImage()
}

func TestDomainIsActive(t *testing.T) {
	dom, conn := openTestDomain(t)
	defer conn.Close()
	defer dom.Free()

	_ = dom.IsActive()
}

func TestDomainIsPersistent(t *testing.T) {
	dom, conn := openTestDomain(t)
	defer conn.Close()
	defer dom.Free()

	_ = dom.IsPersistent()
}

func TestDomainIsUpdated(t *testing.T) {
	dom, conn := openTestDomain(t)
	defer conn.Close()
	defer dom.Free()

	_ = dom.IsUpdated()
}

func TestDomainOSType(t *testing.T) {
	dom, conn := openTestDomain(t)
	defer conn.Close()
	defer dom.Free()

	os, err := dom.OSType()
	if err != nil {
		t.Fatal(err)
	}

	if len(os) == 0 {
		t.Error("empty domain OS type")
	}
}

func TestDomainName(t *testing.T) {
	dom, conn := openTestDomain(t)
	defer conn.Close()
	defer dom.Free()

	name := dom.Name()

	if len(name) == 0 {
		t.Error("empty domain name")
	}
}

func TestDomainHostname(t *testing.T) {
	// Hostname is not supported by the "QEMU" driver
	dom, conn := openTestDomain(t)
	defer conn.Close()
	defer dom.Free()

	if _, err := dom.Hostname(); err == nil {
		t.Error("Hostname should not be supported by the \"QEMU\" driver")
	}
}

func TestDomainID(t *testing.T) {
	dom, conn := openTestDomain(t)
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
	dom, conn := openTestDomain(t)
	defer conn.Close()
	defer dom.Free()

	uuid, err := dom.UUID()
	if err != nil {
		t.Fatal(err)
	}

	if len(uuid) == 0 {
		t.Error("empty domain UUID")
	}
}

func TestDomainXML(t *testing.T) {
	dom, conn := openTestDomain(t)
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
	dom, conn := openTestDomain(t)
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
}
