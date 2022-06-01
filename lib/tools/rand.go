package tools

import (
	"math/rand"
	"time"
)

func Rand(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randKey := rand.Intn(max) + min
	return randKey
}
