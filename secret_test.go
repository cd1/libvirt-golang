package libvirt

import (
	"testing"
)

func TestSecretInit(t *testing.T) {
	env := newTestEnvironment(t).withSecret()
	defer env.cleanUp()

	uuid, err := env.sec.UUID()
	if err != nil {
		t.Error(err)
	}

	if uuid != env.secData.UUID {
		t.Errorf("wrong test secret UUID; got=%v, want=%v", uuid, env.secData.UUID)
	}

	xml, err := env.sec.XML()
	if err != nil {
		t.Error(err)
	}

	if len(xml) == 0 {
		t.Error("empty secret XML")
	}

	usageID, err := env.sec.UsageID()
	if err != nil {
		t.Error(err)
	}

	if usageID != env.secData.UsageName {
		t.Errorf("wrong test secret ID; got=%v, want=%v", usageID, env.secData.UsageName)
	}

	usageType, err := env.sec.UsageType()
	if err != nil {
		t.Error(err)
	}

	if usageType != env.secData.UsageType {
		t.Errorf("wrong test secret usage type; got=%v, want=%v", usageType, env.secData.UsageType)
	}
}

func TestSecretValue(t *testing.T) {
	env := newTestEnvironment(t).withSecret()
	defer env.cleanUp()

	if err := env.sec.SetValue(env.secData.Value); err != nil {
		t.Fatal(err)
	}

	value, err := env.sec.Value()
	if err != nil {
		t.Fatal(err)
	}

	if value != env.secData.Value {
		t.Errorf("wrong secret value; got=%v, want=%v", value, env.secData.Value)
	}
}
