package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"sync"
	"syscall"
	"time"
)

var m runtime.MemStats
var packetsRx int

// Type is the type of packet received
type Type byte

const (
	Join Type = iota
	Part
	Move
)

func (pt Type) String() string {
	return [...]string{"Join", "Part", "Move"}[pt]
}

// TypeHandler ...
func TypeHandler(pt Type) {
	switch pt {
	case Join:
		fmt.Println("Join Packet")
	case Part:
		fmt.Println("Part Packet")
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

type Connection struct {
	NetAddr    net.Addr
	LastUpdate time.Time
}

type ConnectedClients struct {
	Clients map[string]Connection
	Mux     *sync.Mutex
}

func (cc *ConnectedClients) AddClient(addr net.Addr) {
	_, exist := cc.Clients[addr.String()]

	if !exist {
		fmt.Println(addr.String())
		cc.Mux.Lock()
		cc.Clients[addr.String()] = Connection{
			NetAddr:    addr,
			LastUpdate: time.Now(),
		}
		cc.Mux.Unlock()
	}
}

// RemoveClient ...
func (cc *ConnectedClients) RemoveClient(key string) {
	cc.Mux.Lock()
	delete(cc.Clients, key)
	cc.Mux.Unlock()
}

// GetClients ...
func (cc *ConnectedClients) GetClients() map[string]Connection {
	rMap := make(map[string]Connection)
	cc.Mux.Lock()
	for key, value := range cc.Clients {
		rMap[key] = value
	}
	cc.Mux.Unlock()
	return rMap
}

// PrintMemUsage ...
func PrintMemUsage() {
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v KiB", bToKb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v KiB", bToKb(m.TotalAlloc))
	fmt.Printf("\tSys = %v KiB", bToKb(m.Sys))
	fmt.Printf("\tMallocs = %v", m.Mallocs)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)

	fmt.Printf("HeapInUse = %v KiB", bToKb(m.HeapInuse))
	fmt.Printf("\tHeapIdle = %v KiB", bToKb(m.HeapIdle))
	fmt.Printf("\tHeapReleased = %v KiB\n", bToKb(m.HeapReleased))

	debug.FreeOSMemory()
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func bToKb(b uint64) uint64 {
	return b / 1024
}

func main() {
	ctx, cancelFn := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	packetsRx = 0

	listener, err := net.ListenPacket("udp", fmt.Sprintf("%s:%d", "127.0.0.1", 8989))
	if err != nil {
		log.Fatalf("could not create udp listener\n%+v\n", err)
	}
	defer listener.Close()

	sendChannel := make(chan []byte)
	buffer := make([]byte, 1024)

	clients := &ConnectedClients{
		Clients: make(map[string]Connection),
		Mux:     &sync.Mutex{},
	}

	// Packet Sender
	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				fmt.Printf("* Packet Sender Go Fn Killed\n")
				return
			case packet := <-sendChannel:
				fmt.Printf("%b", packet)

				cl := clients.GetClients()
				for _, v := range cl {
					go func(v Connection) {
						_, err := listener.WriteTo(packet, v.NetAddr)
						if err != nil {
							fmt.Printf("%+v", err)
							return
						}
						packetsRx++
						fmt.Printf("packet-written: bytes=%d to=%s\n", len(packet), v.NetAddr.String())
					}(v)
				}
			}
		}
	}(wg)

	// Packet Receiver
	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				fmt.Printf("* Packet Receiver Go Fn Killed\n")
				return
			default:
				n, addr, err := listener.ReadFrom(buffer)
				if err != nil {
					continue
				}

				fmt.Printf("packet-received: bytes=%d from=%s packet=%s\n",
					n, addr.String(), string(buffer[:n]))

				clients.AddClient(addr)

				sendChannel <- buffer[:n]
			}
		}
	}(wg)

	// Connection Timeout Cleanup && Stats Reporting
	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				fmt.Printf("* Connection Cleanup Go Fn Killed\n")
				return
			default:
				time.Sleep(time.Second * 5)

				cl := clients.GetClients()
				currTime := time.Now()

				fmt.Printf("[%s] Current Clients: %d | Packets Written: %d\n", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"), len(cl), packetsRx)
				for k, v := range cl {
					if currTime.Sub(v.LastUpdate) > time.Second*5 {
						clients.RemoveClient(k)
					}
				}
				PrintMemUsage()
			}
		}
	}(wg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Printf("** Shutting Down Server **\n")

	cancelFn()
	listener.Close()
	wg.Wait()
	close(sendChannel)

	fmt.Printf("** Shut Down **\n")
}
