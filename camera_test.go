package visca

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTCPConnectionFromString(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		defer ln.Close()
		ln.Accept()
	}()
	conn, err := connectionFromString(ln.Addr().String(), 10*time.Millisecond)
	// fmt.Printf("Error: %v\n", err)
	assert.Equal(t, nil, err, "does not error")
	assert.NotEqual(t, nil, conn, "conn is not nil")
}

func TestUDPConnectionFromString(t *testing.T) {
	sAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		t.Errorf("Error getting UDP addr: %v", err)
		return
	}
	ln, err := net.ListenUDP("udp", sAddr)
	if err != nil {
		t.Errorf("ListenUDP had an error: %v", err)
		return
	}
	defer ln.Close()
	t.Logf("addr: %v", ln.LocalAddr().String())
	connString := strings.Join([]string{"udp://", ln.LocalAddr().String()}, "")
	conn, err := connectionFromString(connString, 10*time.Millisecond)
	if err != nil {
		t.Errorf("Unexpected err: %v", err)
	}
	assert.NotNil(t, conn, "conn should not be nil")
}

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
	assert.NotNil(t, err, "err should NOT be nil")
	assert.Nil(t, cam, "cam should be nil")
}

func TestConnType(t *testing.T) {
	cType := ConnectionType(2)
	assert.Equal(t, "UDP", cType.String(), "should equal")
	assert.True(t, cType.isValid())

	cType = ConnectionType(1)
	assert.Equal(t, "TCP", cType.String(), "should equal")
	assert.True(t, cType.isValid())

	cType = ConnectionType(3)
	assert.Equal(t, "CharDev", cType.String(), "should equal")
	assert.True(t, cType.isValid())

	cType = ConnectionType(4)
	assert.Equal(t, "Unix", cType.String(), "should equal")
	assert.True(t, cType.isValid())

	cType = ConnectionType(99)
	assert.Equal(t, "ConnectionType(99)", cType.String(), "should equal")
	assert.False(t, cType.isValid())
}
