package libvirt

import (
	"testing"
)

func TestStorageVolumeInit(t *testing.T) {
	env := newTestEnvironment(t).withStorageVolume()
	defer env.cleanUp()

	key, err := env.vol.Key()
	if err != nil {
		t.Error(err)
	}

	if l := len(key); l == 0 {
		t.Error("empty storage volume key")
	}

	name, err := env.vol.Name()
	if err != nil {
		t.Error(err)
	}

	if name != env.volData.Name {
		t.Errorf("unexpected storage volume name; got=%v, want=%v", name, env.volData.Name)
	}

	path, err := env.vol.Path()
	if err != nil {
		t.Error(err)
	}

	if l := len(path); l == 0 {
		t.Error("empty storage volume path")
	}

	xml, err := env.vol.XML()
	if err != nil {
		t.Error(err)
	}

	if l := len(xml); l == 0 {
		t.Error("empty XML descriptor")
	}

	typ, err := env.vol.InfoType()
	if err != nil {
		t.Error(err)
	}

	if typ != VolTypeFile {
		t.Errorf("unexpected storage volume type; got=%v want=%v", typ, VolTypeFile)
	}

	_, err = env.vol.InfoCapacity()
	if err != nil {
		t.Error(err)
	}

	_, err = env.vol.InfoAllocation()
	if err != nil {
		t.Error(err)
	}
}
