package libvirt

import (
	"bytes"
	"os"
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

	ref, err := conn.Close()
	if err != nil {
		t.Fatal(err)
	}

	if ref != 0 {
		t.Errorf("unexpected connection reference count after closing connection; got=%v, want=0", ref)
	}
}

func TestConnectionOpenDefault(t *testing.T) {
	conn, err := OpenDefault()
	if err != nil {
		t.Fatal(err)
	}

	ref, err := conn.Close()
	if err != nil {
		t.Error(err)
	}

	if ref != 0 {
		t.Errorf("unexpected connection reference count after closing connection; got=%v, want=0", ref)
	}
}

func TestConnectionRef(t *testing.T) {
	env := newTestEnvironment(t)
	defer env.cleanUp()

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
}

func TestConnectionReadOnly(t *testing.T) {
	roConn, err := Open(testConnectionURI, ReadOnly, testLogOutput)
	if err != nil {
		t.Fatal(err)
	}
	defer roConn.Close()

	rwConn, err := Open(testConnectionURI, ReadWrite, testLogOutput)
	if err != nil {
		t.Fatal(err)
	}
	defer rwConn.Close()

	var xml bytes.Buffer
	domData, err := newTestDomainData(rwConn)
	if err != nil {
		t.Fatal(err)
	}
	defer domData.cleanUp(rwConn)

	if err = testDomainTmpl.Execute(&xml, domData); err != nil {
		t.Error(err)
	}

	if _, err = roConn.DefineDomain(xml.String()); err == nil {
		t.Error("a readonly libvirt connection should not allow defining domains")
	}

	if _, err = roConn.CreateDomain(xml.String(), DomCreateDefault); err == nil {
		t.Error("a readonly libvirt connection should not allow creating domains")
	}

	xml.Reset()
	secData := newTestSecretData()

	if err = testSecretTmpl.Execute(&xml, secData); err != nil {
		t.Error(err)
	}

	if _, err = roConn.DefineSecret(xml.String()); err == nil {
		t.Error("a readonly libvirt connection should not allow defining secrets")
	}

	xml.Reset()
	poolData, err := newTestStoragePoolData()
	if err != nil {
		t.Fatal(err)
	}
	defer poolData.cleanUp()

	if err = testStoragePoolTmpl.Execute(&xml, poolData); err != nil {
		t.Error(err)
	}

	if _, err = roConn.DefineStoragePool(xml.String()); err == nil {
		t.Error("a readonly libvirt connection should not allow defining storage pools")
	}

	if _, err = roConn.CreateStoragePool(xml.String()); err == nil {
		t.Error("a readonly libvirt connection should not allow creating storage pools")
	}
}

func TestConnectionInit(t *testing.T) {
	env := newTestEnvironment(t)
	defer env.cleanUp()

	alive, err := env.conn.IsAlive()
	if err != nil {
		t.Error(err)
	}
	if !alive {
		t.Error("the libvirt connection was opened but it is not alive")
	}

	encrypted, err := env.conn.IsEncrypted()
	if err != nil {
		t.Error(err)
	}
	if encrypted {
		t.Error("the libvirt connection is encrypted (but it should not)")
	}

	secure, err := env.conn.IsSecure()
	if err != nil {
		t.Error(err)
	}
	if !secure {
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
	// XXX: the function "<Connection>.Sysinfo" returns an error when the connection URI is "qemu:///session"
	if testConnectionURI == "qemu:///session" {
		if err == nil {
			t.Error("the function \"<Connection>.Sysinfo\" isn't supported when the connection URI is \"qemu:///session\",",
				"so it should always return an error")
		}
	} else {
		if err != nil {
			t.Error(err)
		}
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
	data, err := newTestDomainData(*env.conn)
	if err != nil {
		t.Fatal(err)
	}
	defer data.cleanUp(*env.conn)

	if err = testDomainTmpl.Execute(&xml, data); err != nil {
		t.Fatal(err)
	}

	if _, err = env.conn.CreateDomain(xml.String(), DomainCreateFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid create flag")
	}

	dom, err := env.conn.CreateDomain(xml.String(), DomCreateDefault)
	if err != nil {
		t.Fatal(err)
	}
	defer dom.Free()

	active, err := dom.IsActive()
	if err != nil {
		t.Error(err)
	}
	if !active {
		t.Error("domain should be active after being created")
	}

	persistent, err := dom.IsPersistent()
	if err != nil {
		t.Error(err)
	}
	if persistent {
		t.Error("domain should not be persistent after being created")
	}

	if err = dom.Destroy(DomainDestroyFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid destroy flag")
	}

	if err = dom.Destroy(DomDestroyDefault); err != nil {
		t.Error(err)
	}
}

func TestConnectionDefineUndefineDomain(t *testing.T) {
	env := newTestEnvironment(t)
	defer env.cleanUp()

	if _, err := env.conn.DefineDomain(""); err == nil {
		t.Error("an error was not returned when defining a domain with empty XML descriptor")
	}

	var xml bytes.Buffer
	data, err := newTestDomainData(*env.conn)
	if err != nil {
		t.Fatal(err)
	}
	defer data.cleanUp(*env.conn)

	if err = testDomainTmpl.Execute(&xml, data); err != nil {
		t.Fatal(err)
	}

	dom, err := env.conn.DefineDomain(xml.String())
	if err != nil {
		t.Fatal(err)
	}
	defer dom.Free()

	active, err := dom.IsActive()
	if err != nil {
		t.Error(err)
	}
	if active {
		t.Error("domain should not be active after being defined")
	}

	persistent, err := dom.IsPersistent()
	if err != nil {
		t.Error(err)
	}
	if !persistent {
		t.Error("domain should be persistent after being defined")
	}

	if err = dom.Create(DomCreateDefault); err != nil {
		t.Error(err)
	}

	active, err = dom.IsActive()
	if err != nil {
		t.Error(err)
	}
	if !active {
		t.Error("domain should be active after being defined and created")
	}

	persistent, err = dom.IsPersistent()
	if err != nil {
		t.Error(err)
	}
	if !persistent {
		t.Error("domain should still be persistent after being defined and created")
	}

	if err = dom.Destroy(DomDestroyDefault); err != nil {
		t.Error(err)
	}

	active, err = dom.IsActive()
	if err != nil {
		t.Error(err)
	}
	if active {
		t.Error("domain should not be active after being defined and destroyed")
	}

	persistent, err = dom.IsPersistent()
	if !persistent {
		t.Error("domain should be persistent after being defined and destroyed")
	}

	if err = dom.Undefine(DomainUndefineFlag(99)); err == nil {
		t.Error("an error was not return when using an invalid undefine flag")
	}

	if err = dom.Undefine(DomUndefineDefault); err != nil {
		t.Error(err)
	}
}

func TestConnectionLookupDomain(t *testing.T) {
	// TODO: if a domain is created with "<Domain>.Create" after
	// "<Connection>.Define", it doesn't see to get an ID. as a workaround, we
	// create it directly with "<Connection>.CreateDomain" because then it works.
	env := newTestEnvironment(t)
	defer env.cleanUp()

	data, err := newTestDomainData(*env.conn)
	if err != nil {
		t.Fatal(err)
	}
	defer data.cleanUp(*env.conn)

	var xml bytes.Buffer

	if err = testDomainTmpl.Execute(&xml, data); err != nil {
		t.Fatal(err)
	}

	dom, err := env.conn.CreateDomain(xml.String(), DomCreateAutodestroy)
	if err != nil {
		t.Fatal(err)
	}
	defer dom.Free()

	// ByID
	if _, err = env.conn.LookupDomainByID(99); err == nil {
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

	name, err := dom.Name()
	if err != nil {
		t.Error(err)
	}

	if name != data.Name {
		t.Errorf("looked up domain with unexpected name; got=%v, want=%v", name, data.Name)
	}

	// ByUUID
	if _, err = env.conn.LookupDomainByUUID(utils.RandomString()); err == nil {
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

func TestConnectionListSecrets(t *testing.T) {
	env := newTestEnvironment(t).withSecret()
	defer env.cleanUp()

	if _, err := env.conn.ListSecrets(SecretListFlag(99)); err == nil {
		t.Error("an error was not returned when using an invalid flag")
	}

	secrets, err := env.conn.ListSecrets(SecListAll)
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range secrets {
		if err = s.Free(); err != nil {
			t.Error(err)
		}
	}
}

func TestConnectionDefineUndefineSecret(t *testing.T) {
	env := newTestEnvironment(t)
	defer env.cleanUp()

	if _, err := env.conn.DefineSecret(""); err == nil {
		t.Error("an error was not returned when using an empty XML descriptor")
	}

	var xml bytes.Buffer

	data := newTestSecretData()

	if err := testSecretTmpl.Execute(&xml, data); err != nil {
		t.Fatal(err)
	}

	sec, err := env.conn.DefineSecret(xml.String())
	if err != nil {
		t.Fatal(err)
	}
	defer sec.Free()

	if err = sec.Undefine(); err != nil {
		t.Error(err)
	}
}

func TestConnectionLookupSecret(t *testing.T) {
	env := newTestEnvironment(t).withSecret()
	defer env.cleanUp()

	if _, err := env.conn.LookupSecretByUUID(utils.RandomString()); err == nil {
		t.Error("an error was not returned when looking up a non-existing secret UUID")
	}

	if _, err := env.conn.LookupSecretByUsage(SecretUsageType(99), ""); err == nil {
		t.Error("an error was not returned when looking up a secret with an invalid usage flag")
	}

	if _, err := env.conn.LookupSecretByUsage(SecUsageTypeNone, ""); err == nil {
		t.Error("an error was not returned when looking up a secret with an empty ID")
	}

	sec, err := env.conn.LookupSecretByUUID(env.secData.UUID)
	if err != nil {
		t.Fatal(err)
	}
	defer sec.Free()

	uuid, err := sec.UUID()
	if err != nil {
		t.Error(err)
	}

	if uuid != env.secData.UUID {
		t.Errorf("wrong secret UUID; got=%v, want=%v", uuid, env.secData.UUID)
	}

	sec, err = env.conn.LookupSecretByUsage(env.secData.UsageType, env.secData.UsageName)
	if err != nil {
		t.Fatal(err)
	}
	defer sec.Free()

	usageType, err := sec.UsageType()
	if err != nil {
		t.Error(err)
	}

	if usageType != env.secData.UsageType {
		t.Errorf("wrong secret usage type; got=%v, want=%v", usageType, env.secData.UsageType)
	}

	usageID, err := sec.UsageID()
	if err != nil {
		t.Error(err)
	}

	if usageID != env.secData.UsageName {
		t.Errorf("wrong secret usage ID; got=%v, want=%v", usageID, env.secData.UsageName)
	}
}

func TestConnectionFindStoragePoolSources(t *testing.T) {
	env := newTestEnvironment(t)
	defer env.cleanUp()

	if _, err := env.conn.FindStoragePoolSources("", ""); err == nil {
		t.Error("an error was not returned when using an empty storage pool type")
	}

	// TODO: how to test this function on success?
}

func TestConnectionListStoragePools(t *testing.T) {
	env := newTestEnvironment(t).withStoragePool()
	defer env.cleanUp()

	if _, err := env.conn.ListStoragePools(StoragePoolListFlag(^uint32(0))); err == nil {
		t.Error("an error was not returned when using an invalid flag")
	}

	storagePools, err := env.conn.ListStoragePools(PoolListAll)
	if err != nil {
		t.Fatal(err)
	}

	for _, pool := range storagePools {
		if err = pool.Free(); err != nil {
			t.Error(err)
		}
	}
}

func TestConnectionDefineUndefineStoragePool(t *testing.T) {
	env := newTestEnvironment(t)
	defer env.cleanUp()

	if _, err := env.conn.DefineStoragePool(""); err == nil {
		t.Error("an error was not returned when defining a storage pool with an empty XML descriptor")
	}

	var xml bytes.Buffer

	data, err := newTestStoragePoolData()
	if err != nil {
		t.Fatal(err)
	}
	defer data.cleanUp()

	if err = testStoragePoolTmpl.Execute(&xml, data); err != nil {
		t.Fatal(err)
	}

	pool, err := env.conn.DefineStoragePool(xml.String())
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Free()

	active, err := pool.IsActive()
	if err != nil {
		t.Error(err)
	}
	if active {
		t.Error("storage pool should not be active after defining it")
	}

	persistent, err := pool.IsPersistent()
	if err != nil {
		t.Error(err)
	}
	if !persistent {
		t.Error("storage pool should be persistent after defining it")
	}

	if err = pool.Create(); err != nil {
		t.Error(err)
	}

	active, err = pool.IsActive()
	if err != nil {
		t.Error(err)
	}
	if !active {
		t.Error("storage pool should be active after starting it")
	}

	if err = pool.Destroy(); err != nil {
		t.Error(err)
	}

	active, err = pool.IsActive()
	if err != nil {
		t.Error(err)
	}
	if active {
		t.Error("storage pool should not be active after destroying it")
	}

	if err = pool.Delete(StoragePoolDeleteFlag(^uint32(0))); err == nil {
		t.Error("an error was not returned when using an invalid delete flag")
	}

	if err = pool.Delete(PoolDeleteNormal); err != nil {
		t.Error(err)
	}

	if err = pool.Undefine(); err != nil {
		t.Error(err)
	}
}

func TestConnectionCreateDestroyStoragePool(t *testing.T) {
	env := newTestEnvironment(t)
	defer env.cleanUp()

	if _, err := env.conn.CreateStoragePool(""); err == nil {
		t.Error("an error was not returned when creating a storage pool with an empty XML descriptor")
	}

	var xml bytes.Buffer

	data, err := newTestStoragePoolData()
	if err != nil {
		t.Fatal(err)
	}
	defer data.cleanUp()

	if err = testStoragePoolTmpl.Execute(&xml, data); err != nil {
		t.Fatal(err)
	}

	pool, err := env.conn.CreateStoragePool(xml.String())
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Free()

	persistent, err := pool.IsPersistent()
	if err != nil {
		t.Error(err)
	}
	if persistent {
		t.Error("storage pool should not be persistent after creating it")
	}

	if err = pool.Destroy(); err != nil {
		t.Error(err)
	}

	// as the pool is not persistent, there's no way to call "Delete" on it
	if err = os.Remove(data.TargetPath); err != nil {
		t.Log(err)
	}
}

func TestConnectionLookupStoragePool(t *testing.T) {
	env := newTestEnvironment(t).withStoragePool()
	defer env.cleanUp()

	if _, err := env.conn.LookupStoragePoolByName(utils.RandomString()); err == nil {
		t.Error("an error was not returned when using a non-existing pool name")
	}

	if _, err := env.conn.LookupStoragePoolByUUID(utils.RandomString()); err == nil {
		t.Error("an error was not returned when using a non-existing pool UUID")
	}

	pool, err := env.conn.LookupStoragePoolByName(env.poolData.Name)
	if err != nil {
		t.Error(err)
	}
	defer pool.Free()

	name, err := pool.Name()
	if err != nil {
		t.Error(err)
	}

	if name != env.poolData.Name {
		t.Errorf("looked up storage pool with unexpected name; got=%v, want=%v", name, env.poolData.Name)
	}

	pool, err = env.conn.LookupStoragePoolByUUID(env.poolData.UUID)
	if err != nil {
		t.Error(err)
	}
	defer pool.Free()

	uuid, err := pool.UUID()
	if err != nil {
		t.Error(err)
	}

	if uuid != env.poolData.UUID {
		t.Errorf("looked up storage pool with unexpected UUID; got=%v, want=%v", uuid, env.poolData.UUID)
	}
}

func TestConnectionLookupVolume(t *testing.T) {
	env := newTestEnvironment(t).withStorageVolume()
	defer env.cleanUp()

	if _, err := env.conn.LookupStorageVolumeByKey(utils.RandomString()); err == nil {
		t.Error("an error was not returned when looking up a storage volume with an invalid key")
	}

	if _, err := env.conn.LookupStorageVolumeByPath(utils.RandomString()); err == nil {
		t.Error("an error was not returned when looking up a storage volume with an invalid path")
	}

	key, err := env.vol.Key()
	if err != nil {
		t.Error(err)
	}

	vol, err := env.conn.LookupStorageVolumeByKey(key)
	if err != nil {
		t.Error(err)
	}
	defer vol.Free()

	lookedUpKey, err := env.vol.Key()
	if err != nil {
		t.Error(err)
	}

	if lookedUpKey != key {
		t.Errorf("unexpected looked up storage volume key; got=%v, want=%v", lookedUpKey, key)
	}

	path, err := env.vol.Path()
	if err != nil {
		t.Error(err)
	}

	vol, err = env.conn.LookupStorageVolumeByPath(path)
	if err != nil {
		t.Error(err)
	}
	defer vol.Free()

	lookedUpPath, err := vol.Path()
	if err != nil {
		t.Error(err)
	}

	if lookedUpPath != path {
		t.Errorf("unexpected looked up storage volume path; got=%v, want=%v", lookedUpPath, path)
	}
}

func TestConnectionNewStream(t *testing.T) {
	env := newTestEnvironment(t)
	defer env.cleanUp()

	str, err := env.conn.NewStream(StrDefault)
	if err != nil {
		t.Fatal(err)
	}

	if err = str.Free(); err != nil {
		t.Error(err)
	}
}

func TestConnectionListInterfaces(t *testing.T) {
	env := newTestEnvironment(t)
	defer env.cleanUp()

	if _, err := env.conn.ListInterfaces(InterfaceListFlag(^uint32(0))); err == nil {
		t.Error("an error was not returned when using an invalid flag")
	}

	interfaces, err := env.conn.ListInterfaces(IfaceListAll)
	if err != nil {
		t.Fatal(err)
	}

	for _, iface := range interfaces {
		if err = iface.Free(); err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkConnectionOpenClose(b *testing.B) {
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
	env := newTestEnvironment(b)
	defer env.cleanUp()

	var xml bytes.Buffer
	data, err := newTestDomainData(*env.conn)
	if err != nil {
		b.Fatal(err)
	}
	defer data.cleanUp(*env.conn)

	if err := testDomainTmpl.Execute(&xml, data); err != nil {
		b.Fatal(err)
	}
	xmlStr := xml.String()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dom, err := env.conn.CreateDomain(xmlStr, DomCreateDefault)
		if err != nil {
			b.Error(err)
		}
		defer dom.Free()

		if err = dom.Destroy(DomDestroyDefault); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkConnectionDefineDomain(b *testing.B) {
	env := newTestEnvironment(b)
	defer env.cleanUp()

	var xml bytes.Buffer
	data, err := newTestDomainData(*env.conn)
	if err != nil {
		b.Fatal(err)
	}
	defer data.cleanUp(*env.conn)

	if err := testDomainTmpl.Execute(&xml, data); err != nil {
		b.Fatal(err)
	}
	xmlStr := xml.String()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dom, err := env.conn.DefineDomain(xmlStr)
		if err != nil {
			b.Error(err)
		}
		defer dom.Free()

		if err = dom.Undefine(DomUndefineDefault); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkConnectionDefineSecret(b *testing.B) {
	env := newTestEnvironment(b)
	defer env.cleanUp()

	var xml bytes.Buffer

	data := newTestSecretData()

	if err := testSecretTmpl.Execute(&xml, data); err != nil {
		b.Fatal(err)
	}
	xmlStr := xml.String()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sec, err := env.conn.DefineSecret(xmlStr)
		if err != nil {
			b.Error(err)
		}
		defer sec.Free()

		if err = sec.Undefine(); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkConnectionDefinePool(b *testing.B) {
	env := newTestEnvironment(b)
	defer env.cleanUp()

	var xml bytes.Buffer

	data, err := newTestStoragePoolData()
	if err != nil {
		b.Fatal(err)
	}
	defer data.cleanUp()

	if err = testStoragePoolTmpl.Execute(&xml, data); err != nil {
		b.Fatal(err)
	}
	xmlStr := xml.String()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pool, err := env.conn.DefineStoragePool(xmlStr)
		if err != nil {
			b.Error(err)
		}
		defer pool.Free()

		if err = pool.Undefine(); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkConnectionCreatePool(b *testing.B) {
	env := newTestEnvironment(b)
	defer env.cleanUp()

	var xml bytes.Buffer

	data, err := newTestStoragePoolData()
	if err != nil {
		b.Fatal(err)
	}
	defer data.cleanUp()

	if err = testStoragePoolTmpl.Execute(&xml, data); err != nil {
		b.Fatal(err)
	}
	xmlStr := xml.String()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pool, err := env.conn.CreateStoragePool(xmlStr)
		if err != nil {
			b.Error(err)
		}
		defer pool.Free()

		if err = pool.Destroy(); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkConnectionNewStream(b *testing.B) {
	env := newTestEnvironment(b)
	defer env.cleanUp()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		str, err := env.conn.NewStream(StrDefault)
		if err != nil {
			b.Error(err)
		}
		defer str.Free()
	}
	b.StopTimer()
}
