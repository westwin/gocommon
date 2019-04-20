package util

import (
	crand "crypto/rand"
	"encoding/base32"
	"io"
	"math/rand"
	"strings"
	"time"
)

// RandString returns a random string with specified length
func RandString(length int, choices string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = choices[rand.Int63()%int64(len(choices))]
	}

	return string(b)
}

var encoding = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")

// RandomID returns a random string which can be used as an ID
func RandomID() string {
	buff := make([]byte, 16) // 128 bit random ID.
	if _, err := io.ReadFull(crand.Reader, buff); err != nil {
		panic(err)

	}
	// Avoid the identifier to begin with number and trim padding
	return string(buff[0]%26+'a') + strings.TrimRight(encoding.EncodeToString(buff[1:]), "=")
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
