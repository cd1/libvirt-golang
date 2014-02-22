package libvirt

import (
	"testing"
)

func TestNewError(t *testing.T) {
	nilError := NewError(nil)
	if nilError != nil {
		t.Error("creating an error with a nil value should return nil")
	}
}
