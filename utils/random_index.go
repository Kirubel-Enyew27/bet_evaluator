package utils

import "math/rand"

func GetRandomIndex(length int) int {
	return rand.Intn(length)
}
