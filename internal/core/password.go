package core

import (
	"crypto/rand"
	"math/big"
)

const (
	charset    = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // human-friendly charset
	groupSize  = 5
	groupCount = 5
	separator  = "-"
)

// GeneratePassword generate a random password formatted as:
// AAAAA-BBBBB-CCCCC-DDDDD-EEEEE
func GeneratePassword() (string, error) {
	totalChars := groupSize * groupCount
	password := make([]byte, 0, totalChars+(groupCount-1))

	max := big.NewInt(int64(len(charset)))

	for i := 0; i < totalChars; i++ {
		// separator every n characters
		if i > 0 && i%groupSize == 0 {
			password = append(password, separator...)
		}

		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}

		password = append(password, charset[n.Int64()])
	}

	return string(password), nil
}
