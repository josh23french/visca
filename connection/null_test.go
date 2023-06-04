package connection

import "testing"

func TestNullInterface(t *testing.T) {
	var _ Connection = new(NullIface)
}
