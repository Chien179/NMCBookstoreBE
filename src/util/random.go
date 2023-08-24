package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijkemnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt32 generates a random interger between min and max
func RandomInt32(min, max int32) int32 {
	return min + rand.Int31n(max-min+1)
}

// RandomInt64 generates a random interger between min and max
func RandomInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomFloat generates a random interger between min and max
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandomString generates a random string of length m
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomSex generates a random sex
func RandomSex() string {
	genders := []string{
		"male",
		"female",
		"other",
	}
	n := rand.Int() % len(genders)

	return genders[n]
}

// RandomOwner generates a random user name
func RandomUser() string {
	return RandomString(6)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@gamil.com", RandomString(6))
}
