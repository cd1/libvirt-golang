package libvirt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
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
</domain>`

// Configuration variables. Feel free to change them.
var (
	testConnectionURI = "qemu:///session"
	testLogOutput     = ioutil.Discard
)

// These variables shouldn't be changed.
var (
	testDomainMetadataTmpl = template.Must(template.New("test-domain-metadata").Parse(testDomainMetadataXML))
	testDomainTmpl         = template.Must(template.New("test-domain").Parse(testDomainXML))
)

// testDomainData contains the data of a domain used for testing.
type testDomainData struct {
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

// testEnvironment represents the environment used for a test function. It is
// responsible for opening the connection to libvirt, creating test domains and
// other resources, and cleaning them up.
type testEnvironment struct {
	conn    *Connection
	dom     *Domain
	domData *testDomainData
	t       testing.TB
}

// newTestDomainData creates new data for a test domain. Some values are
// generated randomly every time this function is called.
func newTestDomainData() *testDomainData {
	var maxMemory uint64 = 1048576 // 1 MiB
	var maxVCPUs int32 = 4

	return &testDomainData{
		Name:              fmt.Sprintf("domain-%v", utils.RandomString()),
		MaxMemory:         maxMemory,
		MaxVCPUs:          maxVCPUs,
		Memory:            uint64(rand.Intn(int(maxMemory)) + 1),
		MetadataContent:   fmt.Sprintf("content-%v", utils.RandomString()),
		MetadataKey:       fmt.Sprintf("key-%v", utils.RandomString()),
		MetadataNamespace: fmt.Sprintf("ns-%v", utils.RandomString()),
		MetadataTag:       fmt.Sprintf("tag-%v", utils.RandomString()),
		OSType:            "hvm",
		Type:              "kvm",
		UUID:              uuid.New(),
		VCPUs:             int32(rand.Intn(int(maxVCPUs)) + 1),
	}
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

	_, err := env.conn.Close()
	if err != nil {
		env.t.Error(err)
	}
}

// withDomain defines a new test domain. The domain "dom" will not be started.
func (env *testEnvironment) withDomain() *testEnvironment {
	data := newTestDomainData()

	var xml bytes.Buffer

	if err := testDomainTmpl.Execute(&xml, data); err != nil {
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
