package libvirt

import (
	"testing"
)

func TestErrorNew(t *testing.T) {
	if nilError := NewError(nil); nilError != nil {
		t.Error("creating an error with a nil value should return nil")
	}
}
