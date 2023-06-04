# visca

This is a low-level VISCA library for Golang.

High-level interfaces are planned for 1.0, but currently incomplete and subject to breaking changes.

The low-level interfaces are Packet, Message, Connection, and Scanner.

The main high-level interface is the Controller. Command and Inquiry interfaces provide an abstraction of the low-level Message type, allowing for automatic ack/completion response pairing in conjunction with a Connection receive queue.

## Supported Protocols

This only supports the VISCA protocol, however it supports that protocol over several connection types:

* Serial
* Network (tcp/udp)

## Low-Level Usage

```golang
import "github.com/josh23french/visca"

// Create a packet to send
// Either parse a Packet from bytes directly:
pkt, err := visca.PacketFromBytes([]byte{
     0x81, 0x01, 0x04, 0x18, 0x01, 0xFF,
})
// Or create it manually from source, dest, message:
pkt, err = visca.NewPacket(0, 1, []byte{0x01, 0x04, 0x18, 0x01})

if err != nil {
    // Deal with invalid packets or other errors
    panic(err)
}

// Open a new Connection
conn, err := visca.NewSerialConnection("/dev/ttyUSB0")
if err != nil {
    // Deal with it
    panic(err)
}

// Start it up
queue := make(chan *visca.Packet) // need a channel to put received packets in
conn.SetReceiveQueue(queue)
err := conn.Start() // this starts the received data processing
if err != nil {
    // Deal with it
    panic(err)
}

// Send the packet we created
err := conn.Send(pkt)
if err != nil {
    // Deal with it!
    panic(err)
}

// get a response packet from the channel
resPkt := <- queue
```

## Higher-Level Usage

### **(not yet implemented!)**

This is an idea of what could be... subject to change. Basically a substitute for a draft design doc.

```golang
import (
    "github.com/josh23french/visca"
    "github.com/josh23french/visca/commands"
)

conn, err := NewConnectionFromString("/dev/ttyUSB0")
if err != nil {
    // Do something about it
    panic(err)
}

command := commands.PanTiltUp(9) // tilt up with given speed

conn.Send(command)
response := <- command.response // wait for a completion/ack/error (if it's an ack, just wait again for a completion/error)
```

## Controller Usage

### **(also not yet implemented)**

```golang
import (
    "github.com/josh23french/visca/controller"
    "time"
)

ctrl := controller.NewController()

cam1, err := NewConnectionFromString("tcp://10.1.2.7:1239")
if err != nil {
    // Do something besides panic
    panic(err)
}

cam2, err := NewConnectionFromString("tcp://10.1.2.8:1239")
if err != nil {
    // Do something besides panic, like try again to add it later
    panic(err)
}

ctrl.AddCamera(1, cam1) // add the camera to the controller
ctrl.AddCamera(2, cam2)

ctrl.PanTiltUp(9) // Tilt up with speed 9
time.Sleep(100 * time.Millisecond)
ctrl.PanTiltStop()

ctrl.SetCamera(2)

ctrl.PanTiltUpLeft(8, 2) // Pan/tilt up+left with speeds 8 and 2, respectively
time.Sleep(250 * time.Millisecond)
ctrl.PanTiltStop()

```

## License

[GNU Lesser General Public License 3.0](LICENSE)

## Support Disclaimer

Paid support can be provided upon request.

No other support will be provided.

If you find bugs, please submit an issue anyway.
