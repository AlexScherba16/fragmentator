package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomDataChunk(length uint32) []byte {
	src := rand.New(rand.NewSource(time.Now().UnixNano()))
	data := make([]byte, length)
	for i := range data {
		data[i] = byte(src.Intn(256))
	}
	return data
}
