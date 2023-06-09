package iocommon

import (
	"bytes"
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

// Writes the given contents to the file at the given path, unless the file already
// exists and its contents already match the given contents.
func WriteFileIfNeeded(filename string, contents []byte, perm os.FileMode) error {
	existingContents, err := ioutil.ReadFile(filename)
	if err == nil && bytes.Equal(existingContents, contents) {
		return nil
	}

	return ioutil.WriteFile(filename, contents, perm)
}

func CopyEmbeddedStaticFiles(destinationDir string, symlink bool, embeddedFiles embed.FS) error {
	if !symlink {
		return copyEmbeddedDir(".", destinationDir, embeddedFiles)
	}

	entries, err := embeddedFiles.ReadDir(".")
	if err != nil {
		return err
	}

	if len(entries) != 1 || !entries[0].IsDir() {
		panic("expected a single embedded directory")
	}

	_, callerFilename, _, _ := runtime.Caller(1)
	symLinkTarget := path.Join(filepath.Dir(callerFilename), entries[0].Name())
	relativeTargetDir, _ := filepath.Rel(path.Dir(destinationDir), symLinkTarget)

	stat, err := os.Lstat(destinationDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if stat.Mode()&os.ModeSymlink == 0 {
			return fmt.Errorf("static headers destination dir %s exists and is not a symlink", destinationDir)
		}

		currentTarget, err := os.Readlink(destinationDir)
		if err != nil {
			return err
		}
		if currentTarget == relativeTargetDir {
			return nil
		}

		err = os.Remove(destinationDir)
		if err != nil {
			return err
		}
	}

	return os.Symlink(relativeTargetDir, destinationDir)
}

func copyEmbeddedDir(sourceDir, destDir string, embeddedFiles embed.FS) error {
	entries, err := embeddedFiles.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	if sourceDir == "." {
		if len(entries) != 1 || !entries[0].IsDir() {
			panic("expected a single embedded directory")
		}
		return copyEmbeddedDir(entries[0].Name(), destDir, embeddedFiles)
	}

	if len(entries) == 0 {
		return nil
	}

	err = os.MkdirAll(destDir, 0775)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {

			err = copyEmbeddedDir(path.Join(sourceDir, entry.Name()), path.Join(destDir, entry.Name()), embeddedFiles)
			if err != nil {
				return err
			}
		} else {
			content, err := embeddedFiles.ReadFile(path.Join(sourceDir, entry.Name()))
			if err != nil {
				return err
			}
			err = WriteFileIfNeeded(path.Join(destDir, entry.Name()), content, 0664)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
