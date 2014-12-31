package libvirt

import (
	"testing"
)

func TestStoragePoolInit(t *testing.T) {
	env := newTestEnvironment(t).withStoragePool()
	defer env.cleanUp()

	name, err := env.pool.Name()
	if err != nil {
		t.Error(err)
	}
	if name != env.poolData.Name {
		t.Errorf("unexpected storage pool name; got=%v, want=%v", name, env.poolData.Name)
	}

	uuid, err := env.pool.UUID()
	if err != nil {
		t.Error(err)
	}

	if uuid != env.poolData.UUID {
		t.Errorf("unexpected storage pool UUID; got=%v, want=%v", uuid, env.poolData.UUID)
	}

	if _, err = env.pool.XML(StorageXMLFlag(^uint32(0))); err == nil {
		t.Error("an error was not returned when using an invalid XML flag")
	}

	xml, err := env.pool.XML(StorageXMLDefault)
	if err != nil {
		t.Error(err)
	}

	if l := len(xml); l == 0 {
		t.Error("empty storage pool XML descriptor")
	}

	state, err := env.pool.InfoState()
	if err != nil {
		t.Error(err)
	}

	if state != PoolStateInactive {
		t.Errorf("unexpected initial storage pool state; got=%v, want=%v", state, PoolStateInactive)
	}

	if err = env.pool.Create(); err != nil {
		t.Fatal(err)
	}

	state, err = env.pool.InfoState()
	if err != nil {
		t.Error(err)
	}

	if state != PoolStateRunning {
		t.Errorf("unexpected storage pool state after starting it; got=%v, want=%v", state, PoolStateRunning)
	}

	capacity, err := env.pool.InfoCapacity()
	if err != nil {
		t.Error(err)
	}

	if capacity == 0 {
		t.Errorf("storage pool capacity should not be zero; got=%v", capacity)
	}

	allocation, err := env.pool.InfoAllocation()
	if err != nil {
		t.Error(err)
	}

	if allocation == 0 {
		t.Errorf("storage pool allocated space should not be zero; got=%v", allocation)
	}

	available, err := env.pool.InfoAvailable()
	if err != nil {
		t.Error(err)
	}

	if sum := allocation + available; sum != capacity {
		t.Errorf("storage pool available space + allocated space should be the same as its total capacity; got=%v, want=%v", sum, capacity)
	}
}
