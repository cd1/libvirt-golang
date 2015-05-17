package libvirt

// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"
import (
	"log"
)

// InterfaceListFlag defines a filter when listing network interfaces.
type InterfaceListFlag uint32

// Possible values for InterfaceListFlag.
const (
	IfaceListAll      InterfaceListFlag = 0
	IfaceListActive   InterfaceListFlag = C.VIR_CONNECT_LIST_INTERFACES_ACTIVE
	IfaceListInactive InterfaceListFlag = C.VIR_CONNECT_LIST_INTERFACES_INACTIVE
)

// Interface holds a libvirt network interface. There are no exported fields.
type Interface struct {
	log          *log.Logger
	virInterface C.virInterfacePtr
}

// Free frees the interface object. The interface itself is unaltered. The data
// structure is freed and should not be used thereafter.
func (iface Interface) Free() error {
	iface.log.Println("freeing interface object...")
	cRet := C.virInterfaceFree(iface.virInterface)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		iface.log.Printf("an error occurred: %v\n", err)
		return err
	}

	iface.log.Println("interface freed")

	return nil
}
