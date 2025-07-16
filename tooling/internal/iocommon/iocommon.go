package iocommon

import (
	"bytes"
	"embed"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

// Writes the given contents to the file at the given path, unless the file already
// exists and its contents already match the given contents.
func WriteFileIfNeeded(filename string, contents []byte, perm os.FileMode) error {
	existingContents, err := os.ReadFile(filename)
	if err == nil && bytes.Equal(existingContents, contents) {
		return nil
	}

	return os.WriteFile(filename, contents, perm)
}

func CopyEmbeddedStaticFiles(destinationDir string, symlink bool, embeddedFiles embed.FS) error {
	if symlink {
		entries, err := embeddedFiles.ReadDir(".")
		if err != nil {
			return err
		}

		if len(entries) != 1 || !entries[0].IsDir() {
			panic("expected a single embedded directory")
		}

		_, callerFilename, _, _ := runtime.Caller(1)
		symLinkTarget := path.Join(filepath.Dir(callerFilename), entries[0].Name())
		// relativeTargetDir, _ := filepath.Rel(path.Dir(destinationDir), symLinkTarget)

		return symLinkEmbeddedDir(".", destinationDir, symLinkTarget, embeddedFiles)
	}

	return copyEmbeddedDir(".", destinationDir, embeddedFiles)
}

func symLinkEmbeddedDir(emdeddedSourceDir, destDir string, relativeTargetDir string, embeddedFiles embed.FS) error {
	entries, err := embeddedFiles.ReadDir(emdeddedSourceDir)
	if err != nil {
		return err
	}

	if emdeddedSourceDir == "." {
		if len(entries) != 1 || !entries[0].IsDir() {
			panic("expected a single embedded directory")
		}
		return symLinkEmbeddedDir(entries[0].Name(), destDir, relativeTargetDir, embeddedFiles)
	}

	stat, err := os.Lstat(destDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		mode := stat.Mode()
		if mode&os.ModeSymlink != 0 {
			currentTarget, err := os.Readlink(destDir)
			if err != nil {
				return err
			}
			if currentTarget == relativeTargetDir {
				return nil
			}

			err = os.Remove(destDir)
			if err != nil {
				return err
			}
		} else if mode&os.ModeDir != 0 {
			for _, entry := range entries {
				if entry.IsDir() {
					err = symLinkEmbeddedDir(path.Join(emdeddedSourceDir, entry.Name()), path.Join(emdeddedSourceDir, entry.Name()), path.Join(relativeTargetDir, entry.Name()), embeddedFiles)
					if err != nil {
						return err
					}
				} else {
					err = os.Symlink(path.Join(relativeTargetDir, entry.Name()), path.Join(destDir, entry.Name()))
					if err != nil {
						if _, ok := err.(*os.LinkError); ok {
							continue // TODO: full check
						}
						return err
					}
				}
			}

			return nil
		}
	}

	return os.Symlink(relativeTargetDir, destDir)
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
