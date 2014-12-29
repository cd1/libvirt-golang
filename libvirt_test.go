package libvirt

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"text/template"

	"code.google.com/p/go-uuid/uuid"
	"github.com/cd1/utils-golang"
)

const testDeviceLogXML = `
<disk type="dir" device="cdrom">
    <driver name="qemu" type="raw" />
    <source dir="/var/log/" />
    <target dev="hdc" />
    <readonly />
</disk>`

const testDeviceTmpXML = `
<disk type="dir" device="cdrom">
    <driver name="qemu" type="raw" />
    <source dir="/var/tmp/" />
    <target dev="hdc" />
    <readonly />
</disk>`

const testDomainMetadataXML = `
<{{.MetadataTag}}>
    {{.MetadataContent}}
</{{.MetadataTag}}>`

const testDomainXML = `
<domain type='{{.Type}}'>
    <name>{{.Name}}</name>
    <uuid>{{.UUID}}</uuid>
    <memory>{{.MaxMemory}}</memory>
    <currentMemory>{{.Memory}}</currentMemory>
    <vcpu current="{{.VCPUs}}">{{.MaxVCPUs}}</vcpu>
    <os>
        <type>{{.OSType}}</type>
    </os>
    <devices>
        <disk type="file">
            <source file="{{.DiskPath}}" />
            <driver name="qemu" type="{{.DiskFormat}}" />
            <target dev="{{.DiskTarget}}" />
        </disk>
    </devices>
</domain>`

const testSecretXML = `
<secret>
    <uuid>{{.UUID}}</uuid>
    <usage type="{{.UsageTypeString}}">
        <name>{{.UsageName}}</name>
    </usage>
</secret>
`

const testSnapshotXML = `
<domainsnapshot>
    <name>{{.Name}}</name>
</domainsnapshot>`

const testStoragePoolXML = `
<pool type="{{.Type}}">
    <name>{{.Name}}</name>
    <target>
        <path>{{.TargetPath}}</path>
    </target>
</pool>`

// Configuration variables. Feel free to change them.
var (
	testConnectionURI = "qemu:///session"
	testLogOutput     = ioutil.Discard
)

// These variables shouldn't be changed.
var (
	testDomainMetadataTmpl = template.Must(template.New("test-domain-metadata").Parse(testDomainMetadataXML))
	testDomainTmpl         = template.Must(template.New("test-domain").Parse(testDomainXML))
	testSecretTmpl         = template.Must(template.New("test-secret").Parse(testSecretXML))
	testSnapshotTmpl       = template.Must(template.New("test-snapshot").Parse(testSnapshotXML))
	testStoragePoolTmpl    = template.Must(template.New("test-storagepool").Parse(testStoragePoolXML))
)

// testDomainData contains the data of a domain used for testing.
type testDomainData struct {
	DiskFormat        string
	DiskPath          string
	DiskSize          int
	DiskTarget        string
	Name              string
	MaxMemory         uint64
	MaxVCPUs          int32
	Memory            uint64
	MetadataContent   string
	MetadataKey       string
	MetadataNamespace string
	MetadataTag       string
	OSType            string
	Type              string
	UUID              string
	VCPUs             int32
}

// testSecretData contains the data of a secret used for testing.
type testSecretData struct {
	UUID            string
	UsageName       string
	UsageType       SecretUsageType
	UsageTypeString string
	Value           string
}

// testSnapshotData contains the data of a snapshot used for testing.
type testSnapshotData struct {
	Name string
}

// testStoragePoolData contains the data of a storage pool used for testing.
type testStoragePoolData struct {
	Name       string
	TargetPath string
	Type       string
}

// testEnvironment represents the environment used for a test function. It is
// responsible for opening the connection to libvirt, creating test domains and
// other resources, and cleaning them up.
type testEnvironment struct {
	conn     *Connection
	dom      *Domain
	domData  *testDomainData
	pool     *StoragePool
	poolData *testStoragePoolData
	sec      *Secret
	secData  *testSecretData
	snap     *Snapshot
	snapData *testSnapshotData
	t        testing.TB
}

// newTestDomainData creates new data for a test domain. Some values are
// generated randomly every time this function is called.
func newTestDomainData() (*testDomainData, error) {
	data := &testDomainData{
		DiskFormat:        "qcow2",
		DiskSize:          rand.Intn(1048576) + 1, // <= 1 MiB
		DiskTarget:        "vda",
		Name:              fmt.Sprintf("domain-%v", utils.RandomString()),
		MaxMemory:         1048576, // 1 MiB
		MaxVCPUs:          4,
		MetadataContent:   fmt.Sprintf("content-%v", utils.RandomString()),
		MetadataKey:       fmt.Sprintf("key-%v", utils.RandomString()),
		MetadataNamespace: fmt.Sprintf("ns-%v", utils.RandomString()),
		MetadataTag:       fmt.Sprintf("tag-%v", utils.RandomString()),
		OSType:            "hvm",
		Type:              "kvm",
		UUID:              uuid.New(),
	}

	// TODO: this path can be looked up only once instead of for every domain data.
	qemuImgPath, err := exec.LookPath("qemu-img")
	if err != nil {
		return nil, err
	}

	diskPath := filepath.Join(os.TempDir(), fmt.Sprintf("%v-%v.%v", data.Name, data.DiskTarget, data.DiskFormat))

	cmd := exec.Command(qemuImgPath, "create", diskPath, "-f", data.DiskFormat, strconv.Itoa(data.DiskSize))
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	data.DiskPath = diskPath
	data.Memory = uint64(rand.Intn(int(data.MaxMemory)) + 1)
	data.VCPUs = int32(rand.Intn(int(data.MaxVCPUs)) + 1)

	return data, nil
}

// cleanUp cleans up the domain data values, like temporary files.
func (data *testDomainData) cleanUp() error {
	return os.Remove(data.DiskPath)
}

// newTestSecretData creates new data for a test secret. The values are
// generated randomly every time this function is called.
func newTestSecretData() *testSecretData {
	var value bytes.Buffer

	encoder := base64.NewEncoder(base64.StdEncoding, &value)
	encoder.Write([]byte(utils.RandomString()))
	encoder.Close()

	return &testSecretData{
		UsageName:       fmt.Sprintf("name-%v", utils.RandomString()),
		UsageType:       SecUsageTypeCeph,
		UsageTypeString: "ceph",
		UUID:            uuid.New(),
		Value:           value.String(),
	}
}

// newTestSnapshotData creates new data for a test snapshot. The values are
// generated randomly every time this function is called.
func newTestSnapshotData() *testSnapshotData {
	return &testSnapshotData{
		Name: fmt.Sprintf("snapshot-%v", utils.RandomString()),
	}
}

// newTestStoragePoolData creates new data for a test storage pool.
// The values are generated randomly every time this function is called.
func newTestStoragePoolData() (*testStoragePoolData, error) {
	path, err := ioutil.TempDir("", "storagepool-")
	if err != nil {
		return nil, err
	}

	data := &testStoragePoolData{
		Name:       fmt.Sprintf("name-%v", utils.RandomString()),
		TargetPath: path,
		Type:       "dir",
	}

	return data, nil
}

// cleanUp cleans up the storage pool data values, like temporary files.
func (data *testStoragePoolData) cleanUp() error {
	return os.Remove(data.TargetPath)
}

// newTestEnvironment creates a new test environment. Basically it opens a
// connection to libvirt.
func newTestEnvironment(t testing.TB) *testEnvironment {
	conn, err := Open(testConnectionURI, ReadWrite, testLogOutput)
	if err != nil {
		t.Fatal(err)
	}

	return &testEnvironment{
		conn: &conn,
		t:    t,
	}
}

// cleanUp cleans up the test environment. The domain "dom" is undefined, if it
// exists, and the connection to libvirt is closed.
func (env *testEnvironment) cleanUp() {
	if env.dom != nil {
		if env.snap != nil {
			if err := env.snap.Delete(SnapDeleteDefault); err != nil {
				env.t.Error(err)
			}

			if err := env.snap.Free(); err != nil {
				env.t.Error(err)
			}
		}

		state, _, err := env.dom.State()
		if err != nil {
			env.t.Error(err)
		}

		if state != DomStateShutoff {
			if err := env.dom.Destroy(DomDestroyDefault); err != nil {
				env.t.Error(err)
			}
		}

		if err := env.dom.Undefine(DomUndefineDefault); err != nil {
			env.t.Error(err)
		}

		if err := env.dom.Free(); err != nil {
			env.t.Error(err)
		}
	}

	if env.domData != nil {
		if err := env.domData.cleanUp(); err != nil {
			env.t.Error(err)
		}
	}

	if env.sec != nil {
		if err := env.sec.Undefine(); err != nil {
			env.t.Error(err)
		}

		if err := env.sec.Free(); err != nil {
			env.t.Error(err)
		}
	}

	if env.pool != nil {
		if err := env.pool.Undefine(); err != nil {
			env.t.Error(err)
		}

		if err := env.pool.Free(); err != nil {
			env.t.Error(err)
		}
	}

	if env.poolData != nil {
		if err := env.poolData.cleanUp(); err != nil {
			env.t.Error(err)
		}
	}

	_, err := env.conn.Close()
	if err != nil {
		env.t.Error(err)
	}
}

// withDomain defines a new test domain. The domain "dom" will not be started.
func (env *testEnvironment) withDomain() *testEnvironment {
	data, err := newTestDomainData()
	if err != nil {
		env.t.Fatal(err)
	}

	var xml bytes.Buffer

	if err = testDomainTmpl.Execute(&xml, data); err != nil {
		env.t.Fatal(err)
	}

	dom, err := env.conn.DefineDomain(xml.String())
	if err != nil {
		env.t.Fatal(err)
	}

	env.domData = data
	env.dom = &dom

	return env
}

func (env *testEnvironment) withSnapshot() *testEnvironment {
	if env.dom == nil {
		env.withDomain()
	}

	data := newTestSnapshotData()

	var xml bytes.Buffer

	if err := testSnapshotTmpl.Execute(&xml, data); err != nil {
		env.t.Fatal(err)
	}

	snap, err := env.dom.CreateSnapshot(xml.String(), SnapCreateDefault)
	if err != nil {
		env.t.Fatal(err)
	}

	env.snapData = data
	env.snap = &snap

	return env
}

// withSecret defines a new test secret.
func (env *testEnvironment) withSecret() *testEnvironment {
	data := newTestSecretData()

	var xml bytes.Buffer

	if err := testSecretTmpl.Execute(&xml, data); err != nil {
		env.t.Fatal(err)
	}

	sec, err := env.conn.DefineSecret(xml.String())
	if err != nil {
		env.t.Fatal(err)
	}

	env.secData = data
	env.sec = &sec

	return env
}

// withStoragePool defines a new test storage pool. The pool "pool" will remain
// inactive.
func (env *testEnvironment) withStoragePool() *testEnvironment {
	data, err := newTestStoragePoolData()
	if err != nil {
		env.t.Fatal(err)
	}

	var xml bytes.Buffer

	if err = testStoragePoolTmpl.Execute(&xml, data); err != nil {
		env.t.Fatal(err)
	}

	pool, err := env.conn.DefineStoragePool(xml.String())
	if err != nil {
		env.t.Fatal(err)
	}

	env.poolData = data
	env.pool = &pool

	return env
}
