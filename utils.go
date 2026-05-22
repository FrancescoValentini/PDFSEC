package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"runtime"

	"golang.org/x/term"
)

const (
	charset    = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // human-friendly charset
	groupSize  = 5
	groupCount = 5
	separator  = "-"
)

// ReadPasswordFromTTY securely reads a password from the terminal with confirmation.
// It operates directly on the TTY device, allowing it to work even when stdin/stdout are redirected.
func ReadPasswordFromTTY(pwdRecipient string) (string, bool, error) {
	tty, err := openTTY()
	if err != nil {
		return "", false, err
	}
	defer tty.Close()

	fd := int(tty.Fd())

	fmt.Fprintf(os.Stderr, "Enter %s passphrase (Leave blank to generate one): ", pwdRecipient)

	pass1, err := term.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", false, err
	}

	if len(pass1) == 0 {
		pw, err := GeneratePassword()
		return pw, true, err // random password
	}

	fmt.Fprintf(os.Stderr, "Confirm %s passphrase: ", pwdRecipient)
	pass2, err := term.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", false, err
	}

	if string(pass1) != string(pass2) {
		return "", false, errors.New("passphrases do not match")
	}

	return string(pass1), false, nil
}

// openTTY opens the system's terminal device for direct I/O operations.
// It provides cross-platform support for accessing the physical terminal
// independent of stdin/stdout redirection.
func openTTY() (*os.File, error) {
	if runtime.GOOS == "windows" {
		return os.OpenFile("CONIN$", os.O_RDWR, 0)
	}
	return os.OpenFile("/dev/tty", os.O_RDWR, 0)
}

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
