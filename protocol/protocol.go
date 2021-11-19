package protocol

// PacketType is the type of packet received
type PacketType byte

const (
	Join PacketType = iota
	Part
	KeepAlive
	Move
)

func (pt PacketType) String() string {
	return [...]string{"join", "part", "keep-alive", "move"}[pt]
}
