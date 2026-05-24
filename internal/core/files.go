package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// checkFileExists returns an error if the file does not exist or if the path points to a folder
func CheckFileExists(path string) error {
	filePointer, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}

	if filePointer.IsDir() {
		return fmt.Errorf("expected a single pdf file not a folder.")
	}
	return err
}

// resolveOutputFile determines the output filename, adding a suffix if none is provided.
func ResolveOutputFile(input, output string) string {
	if output != "" {
		return output
	}
	ext := filepath.Ext(input)
	name := strings.TrimSuffix(input, ext)
	return name + FileSuffix + ext
}
