package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	salt := fmt.Sprintf("%010d", r.Int31())
	fmt.Println(salt)
}
