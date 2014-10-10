package libvirt

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cd1/utils-golang"
)

const (
	DomTestDevice1XMLFile    = "res/cdrom-log.xml"
	DomTestDevice2XMLFile    = "res/cdrom-tmp.xml"
	DomTestMaxMemory         = 131072 // KiB
	DomTestMemory            = 131072 // KiB
	DomTestMetadataContent   = "<message>Hello world</message>"
	DomTestMetadataKey       = "golang"
	DomTestMetadataNamespace = "code.google.com/p/libvirt-golang"
	DomTestName              = "golang-test"
	DomTestOSType            = "hvm"
	DomTestUUID              = "9652e5cd-15f1-49ad-af73-63a502a9e2b8"
	DomTestVCPUs             = 1
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
	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	if dom.Autostart() {
		t.Error("test domain should not have autostart enabled by default")
	}

	if err := dom.SetAutostart(true); err != nil {
		t.Fatal(err)
	}

	if !dom.Autostart() {
		t.Error("test domain should have autostart enabled after setting")
	}
}

func TestDomainHasCurrentSnapshot(t *testing.T) {
	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	if dom.HasCurrentSnapshot() {
		t.Error("test domain should not have current snapshot")
	}
}

func TestDomainIsActive(t *testing.T) {
	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	if dom.IsActive() {
		t.Error("test domain should not be active")
	}
}

func TestDomainIsPersistent(t *testing.T) {
	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	if !dom.IsPersistent() {
		t.Error("test domain should be persistent")
	}
}

func TestDomainIsUpdated(t *testing.T) {
	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	if dom.IsUpdated() {
		t.Error("test domain should not have been updated")
	}
}

func TestDomainOSType(t *testing.T) {
	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	os, err := dom.OSType()
	if err != nil {
		t.Fatal(err)
	}

	if os != DomTestOSType {
		t.Errorf("wrong test domain OS type; got=%s, want=%s", os, DomTestOSType)
	}
}

func TestDomainName(t *testing.T) {
	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	name := dom.Name()

	if name != DomTestName {
		t.Errorf("wrong test domain name; got=%s, want=%s", name, DomTestName)
	}
}

func TestDomainHostname(t *testing.T) {
	// Hostname is not supported by the "QEMU" driver
	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	if _, err := dom.Hostname(); err == nil {
		t.Error("Hostname should not be supported by the \"QEMU\" driver")
	}
}

func TestDomainID(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	_, err := dom.ID()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDomainUUID(t *testing.T) {
	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	uuid, err := dom.UUID()
	if err != nil {
		t.Fatal(err)
	}

	if uuid != DomTestUUID {
		t.Errorf("wrong test domain UUID; got=%s, want=%s", uuid, DomTestUUID)
	}
}

func TestDomainXML(t *testing.T) {
	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

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
	const newMetadata = `
        <messages>
            <m1>foo</m1>
            <m2>bar</m2>
        </messages>
    `

	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	if err := dom.SetMetadata(DomainMetadataType(99), "", "", "", DomAffectCurrent); err == nil {
		t.Error("an error was not returned when using an invalid type to set a domain metadata")
	}

	if err := dom.SetMetadata(DomMetaElement, "", "", "", DomAffectCurrent); err == nil {
		t.Error("an error was not returned when using an empty content to set a domain metadata")
	}

	if err := dom.SetMetadata(DomMetaElement, newMetadata, DomTestMetadataKey, DomTestMetadataNamespace, DomainModificationImpact(99)); err == nil {
		t.Error("an error was not returned when using an invalid impact config to set a domain metadata")
	}

	if _, err := dom.Metadata(99, "", DomAffectCurrent); err == nil {
		t.Error("an error was not returned when using an invalid type to get a domain metadata")
	}

	if _, err := dom.Metadata(DomMetaElement, utils.RandomString(), DomAffectCurrent); err == nil {
		t.Error("an error was not returned when using a non-existing metadata tag to get a domain metadata")
	}

	if _, err := dom.Metadata(DomMetaElement, "", 99); err == nil {
		t.Error("an error was not returned when using an invalid impact config to get a domain metadata")
	}

	metadata, err := dom.Metadata(DomMetaElement, DomTestMetadataNamespace, DomAffectCurrent)
	if err != nil {
		t.Fatal(err)
	}

	if metadata != DomTestMetadataContent {
		t.Errorf("wrong metadata content; got=\"%s\", want=\"%s\"", metadata, DomTestMetadataContent)
	}

	if err = dom.SetMetadata(DomMetaElement, newMetadata, DomTestMetadataKey, DomTestMetadataNamespace, DomAffectCurrent); err != nil {
		t.Fatal(err)
	}

	metadata, err = dom.Metadata(DomMetaElement, DomTestMetadataNamespace, DomAffectCurrent)
	if err != nil {
		t.Fatal(err)
	}

	if metadata != strings.TrimSpace(newMetadata) {
		t.Errorf("wrong metadata content; got=\"%s\", want=\"%s\"", metadata, newMetadata)
	}
}

func TestDomainReboot(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if err := dom.Reboot(DomainRebootFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid reboot flag")
	}

	if err := dom.Reboot(DomRebootDefault); err != nil {
		t.Error(err)
	}
}

func TestDomainReset(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if err := dom.Reset(); err != nil {
		t.Error(err)
	}
}

func TestDomainShutdown(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if err := dom.Shutdown(); err != nil {
		t.Error(err)
	}
}

func TestDomainSuspendResume(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	state, reason, err := dom.State()
	if err != nil {
		t.Error(err)
	}

	if state != DomStateRunning {
		t.Errorf("unexpected domain state; got=%d (reason %d), want=%d", state, reason, DomStateRunning)
	}

	if err = dom.Suspend(); err != nil {
		t.Error(err)
	}

	state, reason, err = dom.State()
	if err != nil {
		t.Error(err)
	}

	if state != DomStatePaused {
		t.Errorf("unexpected domain state; got=%d (reason %d), want=%d", state, reason, DomStatePaused)
	}

	if err = dom.Resume(); err != nil {
		t.Error(err)
	}

	state, reason, err = dom.State()
	if err != nil {
		t.Error(err)
	}

	if state != DomStateRunning {
		t.Errorf("unexpected domain state; got=%d (reason %d), want=%d", state, reason, DomStateRunning)
	}
}

func TestDomainCoreDump(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if err := dom.CoreDump(".", DomDumpLive); err == nil {
		t.Error("a core dump file should not be generated into a directory path")
	}

	dumpFile := DomTestName + ".core"

	if err := dom.CoreDump(dumpFile, DomainDumpFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid core dump flag")
	}

	if err := dom.CoreDump(dumpFile, DomDumpLive); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dumpFile)

	if _, err := os.Stat(dumpFile); os.IsNotExist(err) {
		t.Errorf("core dump file was not generated [%s]", err)
	}
}

func TestDomainRef(t *testing.T) {
	dom, conn := defineTestDomain(t)
	defer conn.Close()

	if err := dom.Undefine(DomUndefineDefault); err != nil {
		t.Fatal(err)
	}

	if err := dom.Ref(); err != nil {
		t.Fatal(err)
	}

	if err := dom.Free(); err != nil {
		t.Error(err)
	}

	if err := dom.Free(); err != nil {
		t.Error(err)
	}
}

func TestDomainMemory(t *testing.T) {
	const newMaxMemory = 1024 * 1024 * 10 // 10 GiB
	const newMemory = 1024 * 1024 * 3     // 3 GiB

	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	if err := dom.SetMemory(0, DomMemoryCurrent); err == nil {
		t.Error("an error was not returned when setting the domain memory to 0")
	}

	if err := dom.SetMemory(newMemory, DomainMemoryFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid flag to set the domain memory")
	}

	memory, err := dom.MaxMemory()
	if err != nil {
		t.Fatal(err)
	}

	if memory != DomTestMaxMemory {
		t.Errorf("wrong domain maximum memory; got=%d, want=%d", memory, DomTestMaxMemory)
	}

	if err = dom.SetMemory(newMaxMemory, DomMemoryMaximum); err != nil {
		t.Fatal(err)
	}

	memory, err = dom.MaxMemory()
	if err != nil {
		t.Fatal(err)
	}
	if memory != newMaxMemory {
		t.Errorf("wrong maximum memory; got=%d, want=%d", memory, newMaxMemory)
	}

	if err := dom.SetMemory(newMaxMemory+1, DomMemoryCurrent); err == nil {
		t.Error("an error was not returned when setting a memory value greater than the maximum allowed")
	}

	if err = dom.SetMemory(newMemory, DomMemoryCurrent); err != nil {
		t.Fatal(err)
	}
}

func TestDomainVCPUs(t *testing.T) {
	const newMaxVCPUs = 8
	const newVCPUs = 3

	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	if err := dom.SetVCPUs(0, DomVCPUsCurrent); err == nil {
		t.Error("an error was not returned when setting an invalid VCPU number")
	}

	if err := dom.SetVCPUs(newVCPUs, DomainVCPUsFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid flag to set VCPU")
	}

	if _, err := dom.VCPUs(DomainVCPUsFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid flag to get VCPU")
	}

	vcpus, err := dom.VCPUs(DomVCPUsCurrent)
	if err != nil {
		t.Fatal(err)
	}

	if vcpus != DomTestVCPUs {
		t.Errorf("wrong VCPUs number; got=%d, want=%d", vcpus, DomTestVCPUs)
	}

	if err = dom.SetVCPUs(newMaxVCPUs, DomVCPUsMaximum); err != nil {
		t.Fatal(err)
	}

	if err = dom.SetVCPUs(newMaxVCPUs+1, DomVCPUsCurrent); err == nil {
		t.Error("an error was not returned when setting a VCPU number greater than the maximum allowed")
	}

	if err = dom.SetVCPUs(newVCPUs, DomVCPUsCurrent); err != nil {
		t.Fatal(err)
	}

	vcpus, err = dom.VCPUs(DomVCPUsCurrent)
	if err != nil {
		t.Fatal(err)
	}
	if vcpus != newVCPUs {
		t.Errorf("wrong VCPUs count; got=%d, want=%d", vcpus, newVCPUs)
	}
}

func TestDomainInfo(t *testing.T) {
	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	state, err := dom.InfoState()
	if err != nil {
		t.Error(err)
	}
	otherState, _, err := dom.State()
	if err != nil {
		t.Error(err)
	}
	if state != otherState {
		t.Errorf("domain states obtained from different functions do not match; state1=%d, state2=%d", state, otherState)
	}

	maxMemory, err := dom.InfoMaxMemory()
	if err != nil {
		t.Error(err)
	}
	otherMaxMemory, err := dom.MaxMemory()
	if err != nil {
		t.Error(err)
	}
	if maxMemory != otherMaxMemory {
		t.Errorf("domain maximum memories obtained from different functions do not match; memory1=%d, memory2=%d", maxMemory, otherMaxMemory)
	}

	vcpus, err := dom.InfoVCPUs()
	if err != nil {
		t.Error(err)
	}
	otherVcpus, err := dom.VCPUs(DomVCPUsCurrent)
	if err != nil {
		t.Error(err)
	}
	if vcpus != uint16(otherVcpus) {
		t.Errorf("numbers of domain VCPUs obtained from different functions do not match; VCPUs1=%d, VCPUs2=%d", vcpus, otherVcpus)
	}

	memory, err := dom.InfoMemory()
	if err != nil {
		t.Error(err)
	}
	if memory != DomTestMemory {
		t.Errorf("wrong memory value; got=%d, want=%d", memory, DomTestMemory)
	}

	if _, err = dom.InfoCPUTime(); err != nil {
		t.Error(err)
	}
}

func TestDomainSaveRestore(t *testing.T) {
	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	if err := dom.Create(DomCreateDefault); err != nil {
		t.Fatal(err)
	}
	defer dom.Destroy(DomDestroyDefault)

	if err := dom.Save("", "", DomSaveDefault); err == nil {
		t.Error("an error was not returned when using an invalid file name")
	}

	file, ioerr := ioutil.TempFile("", "test-save-restore_")
	if ioerr != nil {
		t.Fatal(ioerr)
	}
	defer os.Remove(file.Name())

	if err := dom.Save(file.Name(), "", DomainSaveFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid save flag")
	}

	if err := dom.Save(file.Name(), "", DomSaveDefault); err != nil {
		t.Error(err)
	}

	state, reason, err := dom.State()
	if err != nil {
		t.Error(err)
	}
	if state != DomStateShutoff {
		t.Errorf("unexpected domain state; got=%d (reason %d), want=%d", state, reason, DomStateShutoff)
	}

	if err = conn.RestoreDomain(file.Name(), "", DomSaveDefault); err != nil {
		t.Error(err)
	}

	state, reason, err = dom.State()
	if err != nil {
		t.Error(err)
	}
	if state != DomStateRunning {
		t.Errorf("unexpected domain state; got=%d (reason %d), want=%d", state, reason, DomStateRunning)
	}
}

func TestDomainDevices(t *testing.T) {
	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	if err := dom.AttachDevice("", DomDeviceModifyCurrent); err == nil {
		t.Error("an error was not returned when attaching an empty XML")
	}

	if err := dom.DetachDevice("", DomDeviceModifyCurrent); err == nil {
		t.Error("an error was not returned when detaching an empty XML")
	}

	if err := dom.UpdateDevice("", DomDeviceModifyCurrent); err == nil {
		t.Error("an error was not returned when updating an empty XML")
	}

	xml, ioerr := ioutil.ReadFile(DomTestDevice1XMLFile)
	if ioerr != nil {
		t.Fatal(ioerr)
	}

	if err := dom.AttachDevice(string(xml), DomainDeviceModifyFlag(99)); err == nil {
		t.Error("an error was not returned when attaching a device with an invalid modify flag")
	}

	if err := dom.DetachDevice(string(xml), DomainDeviceModifyFlag(99)); err == nil {
		t.Error("an error was not returned when detaching a device with an invalid modify flag")
	}

	if err := dom.UpdateDevice(string(xml), DomainDeviceModifyFlag(99)); err == nil {
		t.Error("an error was not returned when updating a device with an invalid modify flag")
	}

	if err := dom.AttachDevice(string(xml), DomDeviceModifyCurrent); err != nil {
		t.Fatal(err)
	}

	xml, ioerr = ioutil.ReadFile(DomTestDevice2XMLFile)
	if ioerr != nil {
		t.Fatal(ioerr)
	}

	if err := dom.UpdateDevice(string(xml), DomDeviceModifyCurrent); err != nil {
		t.Error(err)
	}

	if err := dom.DetachDevice(string(xml), DomDeviceModifyCurrent); err != nil {
		t.Error(err)
	}
}

func TestDomainManagedSave(t *testing.T) {
	dom, conn := defineTestDomain(t)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	if err := dom.Create(DomCreateDefault); err != nil {
		t.Fatal(err)
	}
	// do not Destroy the domain - it will be already destroyed in the end

	if err := dom.ManagedSave(DomainSaveFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid save flag")
	}

	if dom.HasManagedSaveImage() {
		t.Error("the test domain should not have a managed save image initially")
	}

	if err := dom.ManagedSave(DomSaveDefault); err != nil {
		t.Error(err)
	}

	if !dom.HasManagedSaveImage() {
		t.Error("the test domain should have a managed save image after creating a managed save image")
	}

	state, reason, err := dom.State()
	if err != nil {
		t.Error(err)
	}
	if expectedState := DomStateShutoff; state != expectedState {
		t.Errorf("unexpected domain state; got=%d (reason %d), want=%d", state, reason, expectedState)
	}

	if err = dom.Create(DomCreateDefault); err != nil {
		t.Error(err)
	}

	if dom.HasManagedSaveImage() {
		t.Error("the test domain should not have a managed save image anymore after starting from an existing managed save image")
	}

	state, reason, err = dom.State()
	if err != nil {
		t.Error(err)
	}
	if expectedState := DomStateRunning; state != expectedState {
		t.Errorf("unexpected domain state; got=%d (reason %d), want=%d", state, reason, expectedState)
	}

	if err := dom.ManagedSave(DomSaveDefault); err != nil {
		t.Error(err)
	}

	state, reason, err = dom.State()
	if err != nil {
		t.Error(err)
	}
	if expectedState := DomStateShutoff; state != expectedState {
		t.Errorf("unexpected domain state; got=%d (reason %d), want=%d", state, reason, expectedState)
	}

	if !dom.HasManagedSaveImage() {
		t.Error("the test domain should have a managed save image after creating a managed save image")
	}

	if err := dom.ManagedSaveRemove(); err != nil {
		t.Error(err)
	}

	if dom.HasManagedSaveImage() {
		t.Error("the test domain should not have a managed save image anymore after removing it")
	}
}

func TestDomainSendKey(t *testing.T) {
	CtrlAltDel := []uint32{29, 56, 111}

	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if err := dom.SendKey(DomKeycodeSetLinux, time.Duration(50)*time.Millisecond, CtrlAltDel); err != nil {
		t.Error(err)
	}
}

func TestDomainSendProcessSignal(t *testing.T) {
	dom, conn := createTestDomain(t, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	if err := dom.SendProcessSignal(0, DomSIGNOP); err == nil {
		t.Error("cannot send a signal to the process 0")
	}

	if err := dom.SendProcessSignal(1, DomSIGNOP); err == nil {
		t.Error("the function \"SendProcessSignal\" should not be supported yet")
	}
}

func BenchmarkSuspendResume(b *testing.B) {
	dom, conn := createTestDomain(b, DomCreateStartAutodestroy)
	defer conn.Close()
	defer dom.Free()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := dom.Suspend(); err != nil {
			b.Error(err)
		}

		if err := dom.Resume(); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkSaveRestore(b *testing.B) {
	dom, conn := defineTestDomain(b)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	if err := dom.Create(DomCreateDefault); err != nil {
		b.Fatal(err)
	}
	defer dom.Destroy(DomDestroyDefault)

	file, ioerr := ioutil.TempFile("", "bench-save-restore_")
	if ioerr != nil {
		b.Fatal(ioerr)
	}
	defer os.Remove(file.Name())

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := dom.Save(file.Name(), "", DomSaveDefault); err != nil {
			b.Error(err)
		}

		if err := conn.RestoreDomain(file.Name(), "", DomSaveDefault); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkManagedSave(b *testing.B) {
	dom, conn := defineTestDomain(b)
	defer conn.Close()
	defer dom.Free()
	defer dom.Undefine(DomUndefineDefault)

	if err := dom.Create(DomCreateDefault); err != nil {
		b.Fatal(err)
	}
	defer dom.Destroy(DomDestroyDefault)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := dom.ManagedSave(DomSaveDefault); err != nil {
			b.Error(err)
		}

		if err := dom.Create(DomCreateDefault); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}
