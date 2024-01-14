package helper

// generate random reference
import (
	"math/rand"
	"time"
)

func GenerateRandomAlphabet() string {
	// Create a new random number generator
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	// Generate a random alphabet string of length 5
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, 5)
	for i := range result {
		result[i] = alphabet[r.Intn(len(alphabet))]
	}

	return string(result)
}
