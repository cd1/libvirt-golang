package libvirt

// #include <stdlib.h>
// #include <libvirt/libvirt.h>
import "C"
import (
	"io"
	"log"
	"unsafe"
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

// Write writes a series of bytes to the stream. This method may block the
// calling application for an arbitrary amount of time. Once an application has
// finished sending data it should call Finish to wait for successful
// confirmation from the driver, or detect any error.
// This method may not be used if a stream source has been registered.
// Errors are not guaranteed to be reported synchronously with the call, but may
// instead be delayed until a subsequent call.
// This function is equivalent to the libvirt function "Send" but it has been
// renamed to "Write" in order to implement the standard interface io.Writer.
func (str Stream) Write(data []byte) (int, error) {
	cData := C.CString(string(data))
	defer C.free(unsafe.Pointer(cData))

	l := len(data)

	str.log.Printf("sending %v bytes to stream...\n", l)
	cRet := C.virStreamSend(str.virStream, cData, C.size_t(l))
	ret := int32(cRet)

	if ret < 0 {
		err := LastError()
		str.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	str.log.Printf("%v bytes sent\n", ret)

	return int(ret), nil
}

// Read reads a series of bytes from the stream. This method may block the
// calling application for an arbitrary amount of time.
// Errors are not guaranteed to be reported synchronously with the call, but may
// instead be delayed until a subsequent call.
// This function is equivalent to the libvirt function "Recv" but it has been
// renamed to "Read" in order to implement the standard interface io.Reader. And
// due to that interface requirement, this function now returns (0, io.EOF)
// instead of (0, nil) when there's nothing left to be read from the stream.
func (str Stream) Read(data []byte) (int, error) {
	dataLen := len(data)

	cData := (*C.char)(C.malloc(C.size_t(dataLen)))
	defer C.free(unsafe.Pointer(cData))

	str.log.Printf("receiving %v bytes from stream...\n", dataLen)
	cRet := C.virStreamRecv(str.virStream, (*C.char)(unsafe.Pointer(cData)), C.size_t(dataLen))
	ret := int32(cRet)

	if ret < 0 {
		err := LastError()
		str.log.Printf("an error occurred: %v\n", err)
		return 0, err
	}

	str.log.Printf("%v bytes received\n", ret)

	if ret == 0 && dataLen > 0 {
		return 0, io.EOF
	}

	newData := C.GoStringN(cData, cRet)
	copy(data, newData)

	return int(ret), nil
}
