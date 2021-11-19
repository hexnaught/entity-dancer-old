package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hexnaught/entity-dancer/protocol"
	"github.com/hexnaught/entity-dancer/server"
	"github.com/oklog/ulid"
)

var startTime time.Time
var randSeed *rand.Rand

func main() {
	startTime = time.Now()
	randSeed = rand.New(rand.NewSource(startTime.UnixNano()))

	server.New("localhost", 8989).
		WithCleanupEvery(5 * time.Second).
		WithConnectionTimeoutLimit(15 * time.Second).
		Start(
			PacketHandler,
		)
}

// PacketHandler takes an argument of packet data from the UDP socket,
// returns the data to be written back to all clients
func PacketHandler(p []byte) *server.Response {
	// log.Printf("%v\n", b)
	log.Printf("[Main.PacketHandler] Packet Bytes: %v\n", p)
	log.Printf("[Main.PacketHandler] P0: %v\n", p[0])
	// log.Printf("Packet Type: %v\n", PacketType(p[0]).String())

	switch protocol.PacketType(p[0]) {
	case protocol.Join:
		log.Println("[RX] Join Packet")
		nwid, err := HandleJoin(p[1:])

		if err != nil {
			return &server.Response{
				ResponseType: server.Interest,
				Data:         []byte(err.Error()),
			}
		}

		return &server.Response{
			ResponseType: server.Interest,
			Data:         []byte(fmt.Sprintf("[SERVER] %s", nwid)),
		}
	case protocol.Part:
		log.Println("[RX] Part Packet")

		return &server.Response{
			ResponseType: server.Interest,
			Data:         []byte("[SERVER] Client Parting"),
		}
	case protocol.KeepAlive:
		log.Println("[RX] Keep Alive Packet")
		return &server.Response{
			ResponseType: server.Self,
			Data:         []byte("[SERVER] Pong"),
		}
		return nil
	case protocol.Move:
		log.Println("[RX] Move Packet")
		log.Println(p)

		return &server.Response{
			ResponseType: server.All,
			Data:         []byte("[SERVER] Moving"),
		}
	default:
		return &server.Response{
			ResponseType: server.Self,
			Data:         []byte("[SERVER] Unknown Packet Received"),
		}
	}
}

func HandleJoin(packetData []byte) (string, error) {
	log.Printf("Got Data: %s\n", string(packetData))
	gulid, err := GenerateULID()

	if err != nil {
		log.Printf("Error Generating ULID: %+v", err)
		return "", err
	}

	bulid, err := gulid.MarshalBinary()
	if err != nil {
		log.Printf("Error Generating ULID: %+v", err)
		return "", err
	}

	log.Printf("Generated ULID: %s | String: %s | Bin: %b\n", gulid, gulid.String(), bulid)
	return gulid.String(), nil
}

func GenerateULID() (ulid.ULID, error) {
	entropy := ulid.Monotonic(randSeed, 0)
	return ulid.New(ulid.Timestamp(startTime), entropy)
}
