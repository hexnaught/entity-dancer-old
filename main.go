package main

import (
	"fmt"
	"time"

	"github.com/hexnaught/entity-dancer/server"
)

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

// TypeHandler ...
func TypeHandler(pt PacketType) {
	switch pt {
	case Join:
		fmt.Println("Join Packet")
	case Part:
		fmt.Println("Part Packet")
	case KeepAlive:
		fmt.Println("Keep Alive Packet")
	case Move:
		fmt.Println("Move Packet")
	}
}

type Vector3 struct {
	x, y, z int16
}

type Quaternion struct {
	x, y, z, w int16
}

type EntityData struct {
	NWID      []byte     // Unique identifier
	Postition Vector3    // Pos
	Rotation  Quaternion // Rot
}

func main() {
	server.New("localhost", 8989).
		WithCleanupEvery(2 * time.Second).
		WithConnectionTimeoutLimit(10 * time.Second).
		Start(
			PacketHandler,
		)
}

// PacketHandler takes an argument of packet data from the UDP socket,
// returns the data to be written back to all clients
func PacketHandler(p []byte) *server.Response {
	// fmt.Printf("%v\n", b)
	fmt.Printf("%v\n", string(p))
	// fmt.Printf("Packet Type: %v\n", PacketType(p[0]).String())

	switch PacketType(p[0]) {
	case Join:
		fmt.Println("Join Packet")
		return nil
	case Part:
		fmt.Println("Part Packet")
		return nil
	case KeepAlive:
		fmt.Println("Keep Alive Packet")
		return nil
	case Move:
		fmt.Println("Move Packet")
		return &server.Response{
			ResponseType: server.All,
			Data:         []byte(p[1:]),
		}
	default:
		return &server.Response{
			ResponseType: server.Self,
			Data:         []byte("Unknown Packet Received"),
		}
	}
}
