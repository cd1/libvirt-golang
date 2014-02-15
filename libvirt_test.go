package libvirt

import (
	"testing"
)

const HYPERVISOR_URI = "qemu:///system"

func TestOpen(t *testing.T) {
	conn, err := Open(HYPERVISOR_URI)
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
}

func TestOpenReadOnly(t *testing.T) {
	conn, err := OpenReadOnly(HYPERVISOR_URI)
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
}

func TestOpenBadUri(t *testing.T) {
	if _, err := Open("xxx"); err == nil {
		t.Error("an error was not returned when connecting to a bad URI")
	}

	if _, err := OpenReadOnly("xxx"); err == nil {
		t.Error("an error was not returned when connecting (RO) to a bad URI")
	}
}

func BenchmarkConnection(b *testing.B) {
	for n := 0; n < b.N; n++ {
		conn, err := Open(HYPERVISOR_URI)
		if err != nil {
			b.Error(err)
		}

		if _, err := conn.Close(); err != nil {
			b.Error(err)
		}
	}
}
