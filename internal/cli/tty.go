package cli

import (
	"errors"
	"fmt"
	"os"
	"pdfsec/internal/core"
	"runtime"

	"golang.org/x/term"
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
		pw, err := core.GeneratePassword()
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
