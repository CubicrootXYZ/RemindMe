package random

import (
	"crypto/rand"
	"math/big"
	mrand "math/rand"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
)

// String returns a random string of length n
func String(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		var randomInt int

		randomBigInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			randomInt = mrand.Intn(len(letters))
			log.Error("Can not generate random integers, fallback to insecure method! " + err.Error())
		}
		randomInt = int(randomBigInt.Int64())

		s[i] = letters[randomInt]
	}
	return string(s)
}
