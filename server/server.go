package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type HandlerFunc func([]byte) *Response

// Type is the type of packet received
type ResponseType byte

const (
	All ResponseType = iota
	Interest
	Self
	None
)

type Response struct {
	ResponseType ResponseType
	Data         []byte
}

type Server interface {
	Start(HandlerFunc)
	WithCleanupEvery(time.Duration) Server
	WithConnectionTimeoutLimit(time.Duration) Server
	WithBufferSize(uint16) Server
	WithTimeout(byte) Server
}

type server struct {
	IPAddress              string
	Port                   uint16
	MaxBufferSize          uint16
	UDPListener            net.PacketConn
	Timeout                time.Duration
	CleanupRateSeconds     time.Duration
	ConnectionTimeoutLimit time.Duration
	Clients                *ConnectedClients
	wg                     *sync.WaitGroup
}

type Client struct {
	NetAddr    net.Addr
	LastUpdate time.Time
}

type ConnectedClients struct {
	Clients map[string]Client
	Mux     *sync.Mutex
}

func (cc *ConnectedClients) AddClient(addr net.Addr) {
	_, exist := cc.Clients[addr.String()]

	if !exist {
		log.Printf("[entity-dancer] Adding Client: %s", addr.String())
	}

	cc.Mux.Lock()
	// TODO: Move to using NWID or other identifier as the key, generated on
	// join before adding as client or passed on via auth service
	cc.Clients[addr.String()] = Client{
		NetAddr:    addr,
		LastUpdate: time.Now(),
	}
	cc.Mux.Unlock()
}

// RemoveClient ...
func (cc *ConnectedClients) RemoveClient(key string) {
	cc.Mux.Lock()
	delete(cc.Clients, key)
	cc.Mux.Unlock()
}

// GetClients ...
func (cc *ConnectedClients) GetClients() map[string]Client {
	rMap := make(map[string]Client)
	cc.Mux.Lock()
	for key, value := range cc.Clients {
		rMap[key] = value
	}
	cc.Mux.Unlock()
	return rMap
}

func New(address string, port uint16) Server {
	listener, err := net.ListenPacket("udp", fmt.Sprintf("%s:%d", "127.0.0.1", 8989))
	if err != nil {
		log.Fatalf("could not create udp listener\n%+v\n", err)
	}

	return &server{
		UDPListener:            listener,
		IPAddress:              address,
		Port:                   port,
		MaxBufferSize:          1024,
		CleanupRateSeconds:     5 * time.Second,
		ConnectionTimeoutLimit: 5 * time.Second,
		Timeout:                5 * time.Second,
		Clients: &ConnectedClients{
			Clients: make(map[string]Client),
			Mux:     &sync.Mutex{},
		},
		wg: &sync.WaitGroup{},
	}
}

func (s *server) WithBufferSize(val uint16) Server {
	s.MaxBufferSize = uint16(val)
	return s
}

func (s *server) WithTimeout(val byte) Server {
	s.Timeout = (time.Duration(val) * time.Second)
	return s
}

func (s *server) WithCleanupEvery(val time.Duration) Server {
	s.CleanupRateSeconds = val
	return s
}

func (s *server) WithConnectionTimeoutLimit(val time.Duration) Server {
	s.ConnectionTimeoutLimit = val
	return s
}

func (s *server) Start(handler HandlerFunc) {
	log.Printf("[entity-dancer] Server Started on %s:%d\n", s.IPAddress, s.Port)
	log.Printf("[entity-dancer] Server Config %+v\n", s)

	ctx, cancelFn := context.WithCancel(context.Background())

	defer s.UDPListener.Close()

	buffer := make([]byte, s.MaxBufferSize)
	sendChannel := make(chan []byte)

	go s.packetReceiver(ctx, buffer, sendChannel, handler)
	go s.packetSender(ctx, sendChannel)
	go s.connectionCleaner(ctx)

	log.Printf("[entity-dancer] Server Is Running...\n")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	s.Stop(cancelFn, []chan []byte{sendChannel})
}

func (s *server) Stop(ctxCancel context.CancelFunc, channels []chan []byte) {

	log.Printf("[entity-dancer] ** Shutting Down Server **\n")

	ctxCancel()
	s.UDPListener.Close()
	s.wg.Wait()

	for _, c := range channels {
		close(c)
	}

	log.Printf("[entity-dancer] ** Shut Down **\n")
}

func (s *server) connectionCleaner(ctx context.Context) {
	s.wg.Add(1)
	defer s.wg.Done()

	for {
		time.Sleep(s.CleanupRateSeconds)

		select {
		case <-ctx.Done():
			log.Printf("[entity-dancer] * Connection Cleanup Go Fn Killed\n")
			return
		default:
			cl := s.Clients.GetClients()
			currTime := time.Now()

			log.Printf("[entity-dancer] Current Clients: %d\n", len(cl))
			for k, v := range cl {
				if currTime.Sub(v.LastUpdate) > s.ConnectionTimeoutLimit {
					s.Clients.RemoveClient(k)
				}
			}
		}
	}
}

func (s *server) packetReceiver(ctx context.Context, buffer []byte, packetSenderChannel chan []byte, handler HandlerFunc) {
	s.wg.Add(1)
	defer s.wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[entity-dancer] * Packet Receiver Go Fn Killed\n")
			return
		default:
			n, addr, err := s.UDPListener.ReadFrom(buffer)
			if err != nil {
				continue
			}

			// fmt.Printf("[entity-dancer] packet-received: bytes=%d from=%s packet=%s\n",
			// 	n, addr.String(), string(buffer[:n]))

			go s.packetHandler(addr, buffer[:n], packetSenderChannel, handler)
		}
	}
}

func (s *server) packetHandler(sendingClient net.Addr, buffer []byte, packetSenderChannel chan []byte, handler HandlerFunc) {
	s.Clients.AddClient(sendingClient)

	resp := handler(buffer)

	if resp == nil {
		return
	}

	switch resp.ResponseType {
	case All:
		packetSenderChannel <- resp.Data
		break
	case Interest:
		// TODO: Interest algorithm, should possibly be a passed in method to run against clients
		// 		 or the implementer uses GetClients and passes the clients that are 'of interest'
		// Interest response type will then only send to clients that are considered of interest
		return
	case Self:
		s.UDPListener.WriteTo(resp.Data, sendingClient)
		return
	case None:
		return
	}
}

func (s *server) packetSender(ctx context.Context, packetSenderChannel chan []byte) {
	s.wg.Add(1)
	defer s.wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[entity-dancer] * Packet Sender Go Fn Killed\n")
			return
		case packet := <-packetSenderChannel:
			if packet == nil {
				continue
			}

			log.Printf("[entity-dancer][SEND] Packet: %v\n", packet)

			cl := s.Clients.GetClients()
			for _, v := range cl {
				go func(v Client) {
					_, err := s.UDPListener.WriteTo([]byte(packet), v.NetAddr)
					if err != nil {
						fmt.Printf("ERROR: %+v\n", err)
						return
					}
					// log.Printf("[entity-dancer] packet-written: bytes=%d to=%s\n", len(packet), v.NetAddr.String())
				}(v)
			}
		}
	}
}
