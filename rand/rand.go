package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)

	nRead, err := rand.Read(b)

	if err != nil {
		return nil, fmt.Errorf("rand: %w", err)
	}

	if nRead < n {
		return nil, fmt.Errorf("rand: %w", err)
	}

	return b, nil
}


// n is the number of bytes to generate a random string
func String(n int) (string, error) {
	b, err := Bytes(n)

	if err != nil {
		return "", fmt.Errorf("rand: %w", err)
	}

	// return a base64 encoded string
	return base64.URLEncoding.EncodeToString(b), nil
}

const SessionTokenBytes = 32

func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}