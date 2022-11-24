package iocommon

import (
	"bytes"
	"io/ioutil"
	"os"
)

// Writes the given contents to the file at the given path, unless the file already
// exists and its contents already match the given contents.
func EnsureFileContents(filename string, contents []byte, perm os.FileMode) error {
	existingContents, err := ioutil.ReadFile(filename)
	if err == nil && bytes.Equal(existingContents, contents) {
		return nil
	}

	return ioutil.WriteFile(filename, contents, perm)
}
