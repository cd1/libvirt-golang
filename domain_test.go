package libvirt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cd1/utils-golang"
)

var testCtrlAltDel = []uint32{29, 56, 111}

func TestDomainInit(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if env.dom.IsUpdated() {
		t.Error("test domain should not have been updated initially")
	}

	os, err := env.dom.OSType()
	if err != nil {
		t.Error(err)
	}

	if os != env.domData.OSType {
		t.Errorf("wrong test domain OS type; got=%v, want=%v", os, env.domData.OSType)
	}

	name := env.dom.Name()

	if name != env.domData.Name {
		t.Errorf("wrong test domain name; got=%v, want=%v", name, env.domData.Name)
	}

	if _, err = env.dom.Hostname(); err == nil {
		t.Error("\"Hostname\" should not be supported by the \"QEMU\" driver")
	}

	uuid, err := env.dom.UUID()
	if err != nil {
		t.Error(err)
	}

	if uuid != env.domData.UUID {
		t.Errorf("wrong test domain UUID; got=%v, want=%v", uuid, env.domData.UUID)
	}
}

func TestDomainAutostart(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if env.dom.Autostart() {
		t.Error("test domain should not have autostart enabled by default")
	}

	if err := env.dom.SetAutostart(true); err != nil {
		t.Fatal(err)
	}

	if !env.dom.Autostart() {
		t.Error("test domain should have autostart enabled after setting")
	}
}

func TestDomainID(t *testing.T) {
	// TODO: if a domain is created with "<Domain>.Create" after
	// "<Connection>.Define", it doesn't see to get an ID. as a workaround, we
	// create it directly with "<Connection>.CreateDomain" because then it works.
	env := newTestEnvironment(t)
	defer env.cleanUp()

	data := newTestDomainData(t)
	defer data.cleanUp()

	var xml bytes.Buffer

	if err := testDomainTmpl.Execute(&xml, data); err != nil {
		t.Fatal(err)
	}

	dom, err := env.conn.CreateDomain(xml.String(), DomCreateAutodestroy)
	if err != nil {
		t.Fatal(err)
	}
	defer dom.Destroy(DomDestroyDefault)

	_, err = dom.ID()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDomainXML(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if _, err := env.dom.XML(99); err == nil {
		t.Error("an error was not returned when using an invalid flag")
	}

	xml, err := env.dom.XML(DomXMLDefault)
	if err != nil {
		t.Fatal(err)
	}

	if len(xml) == 0 {
		t.Error("empty domain XML")
	}
}

func TestDomainMetadata(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if err := env.dom.SetMetadata(DomainMetadataType(99), "", "", "", DomAffectCurrent); err == nil {
		t.Error("an error was not returned when using an invalid type to set a domain metadata")
	}

	if err := env.dom.SetMetadata(DomMetaElement, "", "", "", DomAffectCurrent); err == nil {
		t.Error("an error was not returned when using an empty content to set a domain metadata")
	}

	var originalMetadataBuf bytes.Buffer

	if err := testDomainMetadataTmpl.Execute(&originalMetadataBuf, env.domData); err != nil {
		t.Fatal(err)
	}
	originalMetadata := originalMetadataBuf.String()

	if err := env.dom.SetMetadata(DomMetaElement, originalMetadata, env.domData.MetadataKey, env.domData.MetadataNamespace, DomainModificationImpact(99)); err == nil {
		t.Error("an error was not returned when using an invalid impact config to set a domain metadata")
	}

	if err := env.dom.SetMetadata(DomMetaElement, originalMetadata, env.domData.MetadataKey, env.domData.MetadataNamespace, DomAffectCurrent); err != nil {
		t.Fatal(err)
	}

	if _, err := env.dom.Metadata(99, "", DomAffectCurrent); err == nil {
		t.Error("an error was not returned when using an invalid type to get a domain metadata")
	}

	if _, err := env.dom.Metadata(DomMetaElement, utils.RandomString(), DomAffectCurrent); err == nil {
		t.Error("an error was not returned when using a non-existing metadata tag to get a domain metadata")
	}

	if _, err := env.dom.Metadata(DomMetaElement, "", 99); err == nil {
		t.Error("an error was not returned when using an invalid impact config to get a domain metadata")
	}

	metadata, err := env.dom.Metadata(DomMetaElement, env.domData.MetadataNamespace, DomAffectCurrent)
	if err != nil {
		t.Fatal(err)
	}

	if metadata != strings.TrimSpace(originalMetadata) {
		t.Errorf("wrong metadata content; got=\"%v\", want=\"%v\"", metadata, originalMetadata)
	}
}

func TestDomainReboot(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if err := env.dom.Reboot(DomRebootDefault); err == nil {
		t.Error("an error was not returned when trying to reboot an offline domain")
	}

	if err := env.dom.Create(DomCreateAutodestroy); err != nil {
		t.Fatal(err)
	}

	if err := env.dom.Reboot(DomainRebootFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid reboot flag")
	}

	if err := env.dom.Reboot(DomRebootDefault); err != nil {
		t.Error(err)
	}
}

func TestDomainReset(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if err := env.dom.Reset(); err == nil {
		t.Error("an error was not returned when trying to reset an offline domain")
	}

	if err := env.dom.Create(DomCreateAutodestroy); err != nil {
		t.Fatal(err)
	}

	if err := env.dom.Reset(); err != nil {
		t.Error(err)
	}
}

func TestDomainShutdown(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if err := env.dom.Shutdown(); err == nil {
		t.Error("an error was not returned when trying to shutdown an offline domain")
	}

	if err := env.dom.Create(DomCreateAutodestroy); err != nil {
		t.Fatal(err)
	}

	if err := env.dom.Shutdown(); err != nil {
		t.Error(err)
	}

	state, reason, err := env.dom.State()
	if err != nil {
		t.Error(err)
	}

	if state != DomStateShutoff && DomainShutoffReason(reason) != DomShutoffReasonShutdown {
		t.Errorf("unexpected domain state; got=%v (reason %v), want=%v (reason %v)", state, reason, DomStateShutoff, DomShutoffReasonShutdown)
	}
}

func TestDomainSuspendResume(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if err := env.dom.Create(DomCreateDefault); err != nil {
		t.Fatal(err)
	}

	if err := env.dom.Suspend(); err != nil {
		t.Error(err)
	}

	state, reason, err := env.dom.State()
	if err != nil {
		t.Error(err)
	}

	if state != DomStatePaused || DomainPausedReason(reason) != DomPausedReasonUser {
		t.Errorf("unexpected domain state; got=%v (reason %v), want=%v (reason %v)", state, reason, DomStatePaused, DomPausedReasonUser)
	}

	if err = env.dom.Resume(); err != nil {
		t.Error(err)
	}

	state, reason, err = env.dom.State()
	if err != nil {
		t.Error(err)
	}

	if state != DomStateRunning || DomainRunningReason(reason) != DomRunningReasonUnpaused {
		t.Errorf("unexpected domain state; got=%v (reason %v), want=%v (reason %v)", state, reason, DomStateRunning, DomRunningReasonUnpaused)
	}
}

func TestDomainCoreDump(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if err := env.dom.Create(DomCreateDefault); err != nil {
		t.Fatal(err)
	}

	if err := env.dom.CoreDump(".", DomDumpLive); err == nil {
		t.Error("a core dump file should not be generated into a directory path")
	}

	dumpFile, ioerr := ioutil.TempFile("", fmt.Sprintf("%v-coredump_", env.domData.Name))
	if ioerr != nil {
		t.Fatal(ioerr)
	}
	defer os.Remove(dumpFile.Name())

	if err := env.dom.CoreDump(dumpFile.Name(), DomainDumpFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid core dump flag")
	}

	if err := env.dom.CoreDump(dumpFile.Name(), DomDumpLive); err != nil {
		t.Fatal(err)
	}

	stat, err := dumpFile.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if stat.Size() == 0 {
		t.Error("core dump file was not generated (empty size)")
	}
}

func TestDomainRef(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if err := env.dom.Undefine(DomUndefineDefault); err != nil {
		t.Fatal(err)
	}

	if err := env.dom.Ref(); err != nil {
		t.Fatal(err)
	}

	if err := env.dom.Free(); err != nil {
		t.Error(err)
	}

	if err := env.dom.Free(); err != nil {
		t.Error(err)
	}

	env.dom = nil
}

func TestDomainMemory(t *testing.T) {
	newMaxMemory := uint64(rand.Intn(1024))
	newMemory := uint64(rand.Intn(int(newMaxMemory)))

	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if err := env.dom.SetMemory(0, DomMemoryCurrent); err == nil {
		t.Error("an error was not returned when setting the domain memory to 0")
	}

	if err := env.dom.SetMemory(newMemory, DomainMemoryModifyFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid flag to set the domain memory")
	}

	maxMemory, err := env.dom.MaxMemory()
	if err != nil {
		t.Fatal(err)
	}

	if maxMemory != env.domData.MaxMemory {
		t.Errorf("wrong domain maximum memory; got=%v, want=%v", maxMemory, env.domData.MaxMemory)
	}

	if err = env.dom.SetMemory(newMaxMemory, DomMemoryMaximum); err != nil {
		t.Fatal(err)
	}

	maxMemory, err = env.dom.MaxMemory()
	if err != nil {
		t.Fatal(err)
	}
	if maxMemory != newMaxMemory {
		t.Errorf("wrong maximum memory; got=%v, want=%v", maxMemory, newMaxMemory)
	}

	if err = env.dom.SetMemory(newMaxMemory+1, DomMemoryCurrent); err == nil {
		t.Error("an error was not returned when setting a memory value greater than the maximum allowed")
	}

	if err = env.dom.SetMemory(newMemory, DomMemoryCurrent); err != nil {
		t.Fatal(err)
	}
}

func TestDomainVCPUs(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	newMaxVCPUs := int32(rand.Intn(int(env.domData.MaxVCPUs)) + 1)
	newVCPUs := int32(rand.Intn(int(newMaxVCPUs)) + 1)

	if err := env.dom.SetVCPUs(0, DomVCPUsCurrent); err == nil {
		t.Error("an error was not returned when setting an invalid VCPU number")
	}

	if err := env.dom.SetVCPUs(uint32(newVCPUs), DomainVCPUsFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid flag to set VCPU")
	}

	if _, err := env.dom.VCPUs(DomainVCPUsFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid flag to get VCPU")
	}

	vcpus, err := env.dom.VCPUs(DomVCPUsCurrent)
	if err != nil {
		t.Fatal(err)
	}
	if vcpus != env.domData.VCPUs {
		t.Errorf("wrong VCPUs number; got=%v, want=%v", vcpus, env.domData.VCPUs)
	}

	maxVCPUs, err := env.dom.VCPUs(DomVCPUsMaximum)
	if err != nil {
		t.Fatal(err)
	}
	if maxVCPUs != env.domData.MaxVCPUs {
		t.Errorf("wrong maximum VCPUs number; got=%v, want=%v", maxVCPUs, env.domData.MaxVCPUs)
	}

	if err = env.dom.SetVCPUs(uint32(newVCPUs), DomVCPUsCurrent); err != nil {
		t.Fatal(err)
	}

	vcpus, err = env.dom.VCPUs(DomVCPUsCurrent)
	if err != nil {
		t.Fatal(err)
	}
	if vcpus != newVCPUs {
		t.Errorf("wrong VCPUs count; got=%v, want=%v", vcpus, newVCPUs)
	}

	if err = env.dom.SetVCPUs(uint32(newMaxVCPUs), DomVCPUsMaximum); err != nil {
		t.Fatal(err)
	}

	maxVCPUs, err = env.dom.VCPUs(DomVCPUsMaximum)
	if err != nil {
		t.Fatal(err)
	}
	if maxVCPUs != newMaxVCPUs {
		t.Errorf("wrong new maximum VCPUs number; got=%v, want=%v", maxVCPUs, newMaxVCPUs)
	}

	if err = env.dom.SetVCPUs(uint32(newMaxVCPUs+1), DomVCPUsCurrent); err == nil {
		t.Error("an error was not returned when setting a VCPU number greater than the maximum allowed")
	}
}

func TestDomainInfo(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if err := env.dom.Create(DomCreateDefault); err != nil {
		t.Fatal(err)
	}

	state, err := env.dom.InfoState()
	if err != nil {
		t.Error(err)
	}
	otherState, _, err := env.dom.State()
	if err != nil {
		t.Error(err)
	}
	if state != otherState {
		t.Errorf("domain states obtained from different functions do not match; state1=%v, state2=%v", state, otherState)
	}

	maxMemory, err := env.dom.InfoMaxMemory()
	if err != nil {
		t.Error(err)
	}
	otherMaxMemory, err := env.dom.MaxMemory()
	if err != nil {
		t.Error(err)
	}
	if maxMemory != otherMaxMemory {
		t.Errorf("domain maximum memories obtained from different functions do not match; memory1=%v, memory2=%v", maxMemory, otherMaxMemory)
	}

	vcpus, err := env.dom.InfoVCPUs()
	if err != nil {
		t.Error(err)
	}
	otherVcpus, err := env.dom.VCPUs(DomVCPUsCurrent)
	if err != nil {
		t.Error(err)
	}
	if vcpus != uint16(otherVcpus) {
		t.Errorf("numbers of domain VCPUs obtained from different functions do not match; VCPUs1=%v, VCPUs2=%v", vcpus, otherVcpus)
	}

	memory, err := env.dom.InfoMemory()
	if err != nil {
		t.Error(err)
	}
	if memory != env.domData.Memory {
		t.Errorf("wrong memory value; got=%v, want=%v", memory, env.domData.Memory)
	}

	if _, err = env.dom.InfoCPUTime(); err != nil {
		t.Error(err)
	}
}

func TestDomainSaveRestore(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if err := env.dom.Create(DomCreateDefault); err != nil {
		t.Fatal(err)
	}

	if err := env.dom.Save("", "", DomSaveDefault); err == nil {
		t.Error("an error was not returned when using an invalid file name")
	}

	file, ioerr := ioutil.TempFile("", fmt.Sprintf("%v-save-restore_", env.domData.Name))
	if ioerr != nil {
		t.Fatal(ioerr)
	}
	defer os.Remove(file.Name())

	if err := env.dom.Save(file.Name(), "", DomainSaveFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid save flag")
	}

	if err := env.dom.Save(file.Name(), "", DomSaveDefault); err != nil {
		t.Error(err)
	}

	state, reason, err := env.dom.State()
	if err != nil {
		t.Error(err)
	}
	if state != DomStateShutoff && DomainShutoffReason(reason) != DomShutoffReasonSaved {
		t.Errorf("unexpected domain state; got=%v (reason %v), want=%v (reason %v)", state, reason, DomStateShutoff, DomShutoffReasonSaved)
	}

	if err = env.conn.RestoreDomain(file.Name(), "", DomSaveDefault); err != nil {
		t.Error(err)
	}

	state, reason, err = env.dom.State()
	if err != nil {
		t.Error(err)
	}
	if state != DomStateRunning && DomainRunningReason(reason) != DomRunningReasonRestored {
		t.Errorf("unexpected domain state; got=%v (reason %v), want=%v (reason %v)", state, reason, DomStateRunning, DomRunningReasonRestored)
	}
}

func TestDomainDevices(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if err := env.dom.AttachDevice("", DomDeviceModifyCurrent); err == nil {
		t.Error("an error was not returned when attaching an empty XML")
	}

	if err := env.dom.DetachDevice("", DomDeviceModifyCurrent); err == nil {
		t.Error("an error was not returned when detaching an empty XML")
	}

	if err := env.dom.UpdateDevice("", DomDeviceModifyCurrent); err == nil {
		t.Error("an error was not returned when updating an empty XML")
	}

	if err := env.dom.AttachDevice(testDeviceLogXML, DomainDeviceModifyFlag(99)); err == nil {
		t.Error("an error was not returned when attaching a device with an invalid modify flag")
	}

	if err := env.dom.DetachDevice(testDeviceLogXML, DomainDeviceModifyFlag(99)); err == nil {
		t.Error("an error was not returned when detaching a device with an invalid modify flag")
	}

	if err := env.dom.UpdateDevice(testDeviceLogXML, DomainDeviceModifyFlag(99)); err == nil {
		t.Error("an error was not returned when updating a device with an invalid modify flag")
	}

	if err := env.dom.AttachDevice(testDeviceLogXML, DomDeviceModifyCurrent); err != nil {
		t.Fatal(err)
	}

	if err := env.dom.UpdateDevice(testDeviceTmpXML, DomDeviceModifyCurrent); err != nil {
		t.Error(err)
	}

	if err := env.dom.DetachDevice(testDeviceTmpXML, DomDeviceModifyCurrent); err != nil {
		t.Error(err)
	}
}

func TestDomainManagedSave(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if err := env.dom.Create(DomCreateDefault); err != nil {
		t.Fatal(err)
	}
	// do not Destroy the domain - it will be already destroyed in the end

	if err := env.dom.ManagedSave(DomainSaveFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid save flag")
	}

	if env.dom.HasManagedSaveImage() {
		t.Error("the test domain should not have a managed save image initially")
	}

	if err := env.dom.ManagedSave(DomSaveDefault); err != nil {
		t.Error(err)
	}

	if !env.dom.HasManagedSaveImage() {
		t.Error("the test domain should have a managed save image after creating a managed save image")
	}

	state, reason, err := env.dom.State()
	if err != nil {
		t.Error(err)
	}
	if expectedState := DomStateShutoff; state != expectedState && DomainShutoffReason(reason) != DomShutoffReasonSaved {
		t.Errorf("unexpected domain state; got=%v (reason %v), want=%v (reason %v)", state, reason, expectedState, DomShutoffReasonSaved)
	}

	if err = env.dom.Create(DomCreateDefault); err != nil {
		t.Error(err)
	}

	if env.dom.HasManagedSaveImage() {
		t.Error("the test domain should not have a managed save image anymore after starting from an existing managed save image")
	}

	state, reason, err = env.dom.State()
	if err != nil {
		t.Error(err)
	}
	if expectedState := DomStateRunning; state != expectedState && DomainRunningReason(reason) != DomRunningReasonRestored {
		t.Errorf("unexpected domain state; got=%v (reason %v), want=%v (reason %v)", state, reason, expectedState, DomRunningReasonRestored)
	}

	if err := env.dom.ManagedSave(DomSaveDefault); err != nil {
		t.Error(err)
	}

	state, reason, err = env.dom.State()
	if err != nil {
		t.Error(err)
	}
	if expectedState := DomStateShutoff; state != expectedState && DomainShutoffReason(reason) != DomShutoffReasonSaved {
		t.Errorf("unexpected domain state; got=%v (reason %v), want=%v (reason %v)", state, reason, expectedState, DomShutoffReasonSaved)
	}

	if !env.dom.HasManagedSaveImage() {
		t.Error("the test domain should have a managed save image after creating a managed save image")
	}

	if err := env.dom.ManagedSaveRemove(); err != nil {
		t.Error(err)
	}

	if env.dom.HasManagedSaveImage() {
		t.Error("the test domain should not have a managed save image anymore after removing it")
	}
}

func TestDomainSendKey(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if err := env.dom.Create(DomCreateDefault); err != nil {
		t.Fatal(err)
	}

	if err := env.dom.SendKey(DomKeycodeSetLinux, time.Duration(50)*time.Millisecond, testCtrlAltDel); err != nil {
		t.Error(err)
	}
}

func TestDomainSendProcessSignal(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if err := env.dom.Create(DomCreateDefault); err != nil {
		t.Fatal(err)
	}

	if err := env.dom.SendProcessSignal(0, DomSIGNOP); err == nil {
		t.Error("cannot send a signal to the process 0")
	}

	if err := env.dom.SendProcessSignal(1, DomSIGNOP); err == nil {
		t.Error("the function \"SendProcessSignal\" should not be supported yet")
	}
}

func TestDomainListSnapshots(t *testing.T) {
	env := newTestEnvironment(t).withSnapshot()
	defer env.cleanUp()

	snapshots, err := env.dom.ListSnapshots(SnapListAll)
	if err != nil {
		t.Fatal(err)
	}

	for _, snap := range snapshots {
		if err := snap.Free(); err != nil {
			t.Error(err)
		}
	}
}

func TestDomainCreateAndDeleteSnapshot(t *testing.T) {
	env := newTestEnvironment(t).withDomain()
	defer env.cleanUp()

	if env.dom.HasCurrentSnapshot() {
		t.Error("test domain should not have a current snapshot initially")
	}

	if _, err := env.dom.CreateSnapshot("", SnapCreateDefault); err == nil {
		t.Error("an error was not returned when using an empty XML descriptor")
	}

	var xml bytes.Buffer
	data := newTestSnapshotData()

	if err := testSnapshotTmpl.Execute(&xml, data); err != nil {
		t.Fatal(err)
	}

	if _, err := env.dom.CreateSnapshot(xml.String(), SnapshotCreateFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid create flag")
	}

	snap, err := env.dom.CreateSnapshot(xml.String(), SnapCreateDefault)
	if err != nil {
		t.Fatal(err)
	}
	defer snap.Free()

	if !env.dom.HasCurrentSnapshot() {
		t.Error("test domain should have a current snapshot")
	}

	if !snap.IsCurrent() {
		t.Error("snapshot should be current (but it's not)")
	}

	if err := snap.Delete(SnapDeleteDefault); err != nil {
		t.Error(err)
	}
}

func TestDomainLookupSnapshot(t *testing.T) {
	env := newTestEnvironment(t).withSnapshot()
	defer env.cleanUp()

	if _, err := env.dom.LookupSnapshotByName(utils.RandomString()); err == nil {
		t.Error("an error was not returned when looking up an invalid snapshot name")
	}

	snap, err := env.dom.LookupSnapshotByName(env.snapData.Name)
	if err != nil {
		t.Error(err)
	}

	if newName := snap.Name(); newName != env.snapData.Name {
		t.Errorf("wrong snapshot name; got=%v, want=%v", newName, env.snapData.Name)
	}
}

func BenchmarkDomainSuspendResume(b *testing.B) {
	env := newTestEnvironment(b).withDomain()
	defer env.cleanUp()

	if err := env.dom.Create(DomCreateDefault); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := env.dom.Suspend(); err != nil {
			b.Error(err)
		}

		if err := env.dom.Resume(); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkDomainSaveRestore(b *testing.B) {
	env := newTestEnvironment(b).withDomain()
	defer env.cleanUp()

	if err := env.dom.Create(DomCreateDefault); err != nil {
		b.Fatal(err)
	}

	file, ioerr := ioutil.TempFile("", fmt.Sprintf("%v-bench-save-restore_", env.domData.Name))
	if ioerr != nil {
		b.Fatal(ioerr)
	}
	defer os.Remove(file.Name())

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := env.dom.Save(file.Name(), "", DomSaveDefault); err != nil {
			b.Error(err)
		}

		if err := env.conn.RestoreDomain(file.Name(), "", DomSaveDefault); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkDomainManagedSave(b *testing.B) {
	env := newTestEnvironment(b).withDomain()
	defer env.cleanUp()

	if err := env.dom.Create(DomCreateDefault); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := env.dom.ManagedSave(DomSaveDefault); err != nil {
			b.Error(err)
		}

		if err := env.dom.Create(DomCreateDefault); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkDomainCreateSnapshot(b *testing.B) {
	env := newTestEnvironment(b).withDomain()
	defer env.cleanUp()

	var xml bytes.Buffer
	data := newTestSnapshotData()

	if err := testSnapshotTmpl.Execute(&xml, data); err != nil {
		b.Fatal(err)
	}

	xmlStr := xml.String()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		snap, err := env.dom.CreateSnapshot(xmlStr, SnapCreateDefault)
		if err != nil {
			b.Error(err)
		}

		if err := snap.Delete(SnapDeleteDefault); err != nil {
			b.Error(err)
		}

		if err := snap.Free(); err != nil {
			b.Error(err)
		}
	}
}
