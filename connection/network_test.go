package connection

import "testing"

func TestTCPInterface(t *testing.T) {
	var _ Connection = new(Network)
}
