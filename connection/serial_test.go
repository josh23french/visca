package connection

import "testing"

func TestSerialInterface(t *testing.T) {
	var _ Connection = new(Serial)
}
