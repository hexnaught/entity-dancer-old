package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

type SomeStruct struct {
	A int8
	B uint16
	C uint32
	D uint16
}

func main() {
	var myStruct SomeStruct
	myStruct.A = 1
	myStruct.B = 2
	myStruct.C = 3
	myStruct.D = 4

	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, myStruct)

	fmt.Printf("Error: %v\n", err)

	fmt.Printf("Sizeof myStruct: %d, Sizeof buf: %d, Len of buf: %d\n", unsafe.Sizeof(myStruct), unsafe.Sizeof(buf), buf.Len())

	newStruct := SomeStruct{}
	err = binary.Read(buf, binary.LittleEndian, &newStruct)
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("%+v", newStruct)
}
