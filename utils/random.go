package utils

import "math/rand"

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int) int {
	return rand.Intn(max - min + 1)
}
