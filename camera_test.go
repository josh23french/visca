package visca

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCamera(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Errorf("Error starting test TCP server: %v", err)
		return
	}
	go func() {
		defer ln.Close()
		ln.Accept()
	}()
	cam, err := NewCamera("Cam6", "tcp://"+ln.Addr().String())
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, "Cam6", cam.Name, "name should be Cam6")
	assert.NotNil(t, cam.conn, "conn should not be nil")
}

func TestNewCameraWithConnectionError(t *testing.T) {
	cam, err := NewCamera("Cam6", "tcp://192.0.2.2:33444") // RFC 5737, TEST-NET-1
	assert.Nil(t, err)
	assert.NotNil(t, cam)
}
