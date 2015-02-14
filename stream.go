package libvirt

// #include <libvirt/libvirt.h>
import "C"
import (
	"log"
)

// StreamFlag defines how a stream should be created.
type StreamFlag uint32

// Possible values for StreamFlag.
const (
	StrDefault  StreamFlag = 0
	StrNonBlock StreamFlag = C.VIR_STREAM_NONBLOCK
)

// Stream holds a libvirt stream. There are no exported fields.
type Stream struct {
	log       *log.Logger
	virStream C.virStreamPtr
}

// Free decrements the reference count on a stream, releasing the stream object
// if the reference count has hit zero.
// There must not be an active data transfer in progress when releasing the
// stream. If a stream needs to be disposed of prior to end of stream being
// reached, then the Abort function should be called first.
func (str Stream) Free() error {
	str.log.Println("freeing stream object...")
	cRet := C.virStreamFree(str.virStream)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		str.log.Printf("an error occurred: %v\n", err)
		return err
	}

	str.log.Println("stream freed")

	return nil
}
