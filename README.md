# visca

This is a low-level VISCA library for Golang.

## Supported Protocols

This does not directly support communications.

## Usage

```golang
import "github.com/josh23french/visca"

pkt, err := visca.PacketFromBytes(bytes)
if err != nil {
  // Deal with it
}

if pkt.Destination() == 1 {
  // Packet is for camera 1
  sendAck(pkt.Source)
  processMessage(pkt.Message)
}
```

## Automatic Address Translation

Take for example an IP camera with address 2. IP-based PTZ protocols generally force the camera address to be 1, even if the camera has a different address on the serial bus. So messages require translation.

When the controller refers to camera 2, the gateway must open a connection to the camera over IP and manage sending packets destined to camera 2 over that connection.

Messages to this camera (with `destination=2`) are sent with `destination=1`, and return messages from the camera (with `source=1`) are translated to show they came from `source=2`.

## License

[GNU Lesser General Public License 3.0](LICENSE)

## Support Disclaimer

Paid support can be provided upon request.

No other support will be provided.
