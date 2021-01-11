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

## License

[GNU Lesser General Public License 3.0](LICENSE)

## Support Disclaimer

Paid support can be provided upon request.

No other support will be provided.
