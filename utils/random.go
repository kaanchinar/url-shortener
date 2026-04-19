package utils

import (
	"math/rand"
	"time"
)

func GenerateUniqueID() string {
	const chars = "ABCDEFGHIJKLMNOPRQSTUWXYZabcdefghijklmnoprqstuwxyz0123456789"
	const length = 7
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	out := make([]byte, length)
	for i := 0; i < length; i++ {
		out[i] = chars[rng.Intn(len(chars))]
	}

	return string(out)
}

