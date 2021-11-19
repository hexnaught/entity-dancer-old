package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

type S struct {
	PType int8
	X     uint16
	Y     uint16
	Z     uint16
}

func (s *S) toBin() bytes.Buffer {
	buf := &bytes.Buffer{}

	binary.Write(buf, binary.LittleEndian, s.PType)
	binary.Write(buf, binary.LittleEndian, s.X)
	binary.Write(buf, binary.LittleEndian, s.Y)
	binary.Write(buf, binary.LittleEndian, s.Z)

	return *buf
}

func toStruct(b *bytes.Buffer) S {
	return S{
		PType: b.Next(1),
		X:     b[1:2],
		Y:     b[3:4],
		Z:     b[5:6],
	}
}

func main() {
	p := make([]byte, 2048)
	conn, err := net.Dial("udp", "127.0.0.1:8989")

	if err != nil {
		fmt.Printf("Some error %v", err)
		return

	}

	buf := &bytes.Buffer{}

	myS := S{3, 8899, 666, 333}

	err = binary.Write(buf, binary.LittleEndian, myS.PType)
	fmt.Printf("Error: %v\n", err)
	err = binary.Write(buf, binary.LittleEndian, myS.X)
	fmt.Printf("Error: %v\n", err)
	err = binary.Write(buf, binary.LittleEndian, myS.Y)
	fmt.Printf("Error: %v\n", err)
	err = binary.Write(buf, binary.LittleEndian, myS.Z)
	fmt.Printf("Error: %v\n", err)

	mySS := &S{}
	err = binary.Read(buf, binary.LittleEndian, mySS)
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("GO oBuf: %v\n", mySS)

	log.Printf("SEND: %v\n", buf.Bytes())

	// fmt.Fprintf(conn, "%b", buf)
	conn.Write(buf.Bytes())
	n, err := bufio.NewReader(conn).Read(p)

	if err == nil {
		log.Printf("READ: %+v\n", p[:n])
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()
}
