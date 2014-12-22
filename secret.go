package libvirt

// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"
import (
	"log"
)

// SecretListFlag defines a filter when listing secrets.
type SecretListFlag uint32

// Possible values for SecretListFlag.
const (
	SecListAll         SecretListFlag = 0
	SecListEphemeral   SecretListFlag = C.VIR_CONNECT_LIST_SECRETS_EPHEMERAL
	SecListNoEphemeral SecretListFlag = C.VIR_CONNECT_LIST_SECRETS_NO_EPHEMERAL
	SecListPrivate     SecretListFlag = C.VIR_CONNECT_LIST_SECRETS_PRIVATE
	SecListNoPrivate   SecretListFlag = C.VIR_CONNECT_LIST_SECRETS_NO_PRIVATE
)

// Secret holds a libvirt secret. There are no exported fields.
type Secret struct {
	log       *log.Logger
	virSecret C.virSecretPtr
}

// Free releases the secret handle. The underlying secret continues to exist.
func (sec Secret) Free() error {
	sec.log.Println("freeing secret...")
	cRet := C.virSecretFree(sec.virSecret)
	ret := int(cRet)

	if ret == -1 {
		err := LastError()
		sec.log.Printf("an error occurred: %v\n", err)
		return err
	}

	sec.log.Println("secret freed")

	return nil
}

// Undefine deletes the specified secret. This does not free the associated
// "Secret" object.
func (sec Secret) Undefine() error {
	sec.log.Println("undefining secret...")
	cRet := C.virSecretUndefine(sec.virSecret)
	ret := int(cRet)

	if ret == -1 {
		err := LastError()
		sec.log.Printf("an error occurred: %v\n", err)
		return err
	}

	sec.log.Println("secret undefined")

	return nil
}
