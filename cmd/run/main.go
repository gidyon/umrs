package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
)

const num = 100000000

var ids = make(map[string]struct{}, num)

func main() {
	var (
		id     string
		exists uint32
	)

	fmt.Println("starting")

	for i := 0; i < num; i++ {
		id = uuid.New().String()
		_, ok := ids[id]
		if ok {
			exists++
		}
		ids[id] = struct{}{}
	}

	fmt.Println("exists: ", exists)
}

var h = sha256.New()

func genID() string {
	str := randomdata.RandStringRunes(32)
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
