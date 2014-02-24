package libvirt

import (
	"testing"
)

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
