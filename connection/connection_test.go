package connection

import (
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPFromString(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		defer ln.Close()
		ln.Accept()
	}()
	conn, err := FromString(ln.Addr().String())
	// fmt.Printf("Error: %v\n", err)
	assert.Nil(t, err)
	assert.NotNil(t, conn, "conn is not nil")

	err = conn.Start()

	assert.Nil(t, err)
	conn.Stop()
}

func TestFromStringError(t *testing.T) {
	_, err := FromString("°Asdf")
	assert.NotNil(t, err, "err is nil")
}

func TestUDPFromString(t *testing.T) {
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
	conn, err := FromString(connString)
	if err != nil {
		t.Errorf("Unexpected err: %v", err)
	}
	assert.NotNil(t, conn, "conn should not be nil")
}
