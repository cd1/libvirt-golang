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

// Abort requests that the in progress data transfer be cancelled abnormally
// before the end of the stream has been reached. For output streams this can be
// used to inform the driver that the stream is being terminated early. For
// input streams this can be used to inform the driver that it should stop
// sending data.
func (str Stream) Abort() error {
	str.log.Println("aborting stream...")
	cRet := C.virStreamAbort(str.virStream)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		str.log.Printf("an error occurred: %v\n", err)
		return err
	}

	str.log.Println("stream aborted")

	return nil
}

// Finish indicates that there is no further data to be transmitted on the
// stream. For output streams this should be called once all data has been
// written. For input streams this should be called once Recv returns
// end-of-file.
// This method is a synchronization point for all asynchronous errors, so if
// this returns a success code the application can be sure that all data has
// been successfully processed.
func (str Stream) Finish() error {
	str.log.Println("finishing stream...")
	cRet := C.virStreamFinish(str.virStream)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		str.log.Printf("an error occurred: %v\n", err)
		return err
	}

	str.log.Println("stream finished")

	return nil
}

// Ref increments the reference count on the stream. For each additional call to
// this method, there shall be a corresponding call to Free to release the
// reference count, once the caller no longer needs the reference to this
// object.
func (str Stream) Ref() error {
	str.log.Println("incrementing stream's reference count...")
	cRet := C.virStreamRef(str.virStream)
	ret := int32(cRet)

	if ret == -1 {
		err := LastError()
		str.log.Printf("an error occurred: %v\n", err)
		return err
	}

	str.log.Println("reference count incremented")

	return nil
}
