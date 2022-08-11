package blog

import (
	"math/rand"
	"time"
)

func GenerateRandomNumbers() int {
	x1 := rand.NewSource(time.Now().UnixNano())
	y1 := rand.New(x1)
	min := 1
	max := 100000
	return y1.Intn(max-min) + min
}
