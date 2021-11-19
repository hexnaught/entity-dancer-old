package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"
)

func main() {

	clientCnt := 505
	packetsCnt := 555

	fmt.Printf("Starting Client Emulator - Clients: %d | Packets %d\n", clientCnt, packetsCnt)

	raddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8989")
	if err != nil {
		fmt.Println("FAILED TO INIT CONN")
		fmt.Println(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for i := 0; i < clientCnt; i++ {
		// time.Sleep(time.Millisecond * 1)
		packetSender, err := net.DialUDP("udp", nil, raddr)
		defer packetSender.Close()
		if err != nil {
			fmt.Println("FAILED TO DIAL CONN")
			fmt.Println(err)
		}
		go func(packetSender *net.UDPConn, i int) {
			for l := 0; l < packetsCnt; l++ {
				packet := fmt.Sprintf("Client: %d | Packet %d", i, l)
				fmt.Println(fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05.999999"), packet))
				packetSender.Write([]byte(
					packet,
				))
			}
			PrintMemUsage()
		}(packetSender, i)
	}

	<-sc

	fmt.Printf("Client Emulator Shutting Down...")

	time.Sleep(5 * time.Second)
	debug.FreeOSMemory()
	time.Sleep(5 * time.Second)
}

// PrintMemUsage ...
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// func main() {

// 	for i := 0; i < 1_000_000; i++ {
// 		go func(i int) {
// 			for l := 0; l < 10_000_000; l++ {
// 				raddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8989")
// 				if err != nil {
// 					fmt.Println("FAILED TO INIT CONN")
// 					fmt.Println(err)
// 					fmt.Println(fmt.Sprintf("Client: %d | Packet %d", i, l))
// 					break
// 				}
// 				packetSender, err := net.DialUDP("udp", nil, raddr)
// 				if err != nil {
// 					fmt.Println("FAILED TO DIAL CONN")
// 					fmt.Println(err)
// 					fmt.Println(fmt.Sprintf("Client: %d | Packet %d", i, l))
// 					break
// 				}
// 				packet := fmt.Sprintf("Client: %d | Packet %d", i, l)
// 				fmt.Println(packet)
// 				packetSender.Write([]byte(
// 					packet,
// 				))
// 			}
// 		}(i)
// 	}

// 	sc := make(chan os.Signal, 1)
// 	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
// 	<-sc
// }

// func main() {

// 	for i := 0; i < 1_000_000; i++ {
// 		raddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8989")
// 		if err != nil {
// 			fmt.Println(err)
// 			break
// 		}
// 		packetSender, err := net.DialUDP("udp", nil, raddr)
// 		if err != nil {
// 			fmt.Println(err)
// 			break
// 		}
// 		packetSender.Write([]byte(
// 			fmt.Sprintf("Packet %d\n", i),
// 		))
// 	}

// 	sc := make(chan os.Signal, 1)
// 	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
// 	<-sc
// }
