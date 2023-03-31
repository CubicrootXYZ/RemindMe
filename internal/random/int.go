package random

import (
	"crypto/rand"
	"math/big"
	mrand "math/rand"
)

// Intn generates a random integer in [0, max). Will attempt to generate with
// secure crypto/rand, if that fails falls back to weak math/rand.
func Intn(max int) int {
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(urlSaveRunes))))
	if err == nil {
		return int(index.Int64())
	}

	return mrand.Intn(len(urlSaveRunes))
}
