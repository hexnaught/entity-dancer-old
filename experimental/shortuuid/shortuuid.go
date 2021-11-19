package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lithammer/shortuuid/v3"
	"github.com/oklog/ulid"
)

func main() {
	u := shortuuid.New() // Cekw67uyMpBGZLRP2HFVbe
	fmt.Println(u)

	ExampleULID()
}

func ExampleULID() {
	t := time.Now()
	r := rand.New(rand.NewSource(t.UnixNano()))
	entropy := ulid.Monotonic(r, 0)

	fmt.Println(ulid.MustNew(ulid.Timestamp(t), entropy))
	fmt.Println(ulid.MustNew(ulid.Timestamp(t), entropy))
	fmt.Println(ulid.MustNew(ulid.Timestamp(t), entropy))
	// Output: 0000XSNJG0MQJHBF4QX1EFD6Y3
}
