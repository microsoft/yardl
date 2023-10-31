// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package packaging

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

var cacheDir string

func init() {
	if err := initCacheDir(); err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}
}

func initCacheDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	cacheDir = filepath.Join(homeDir, ".yardl", "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}
	return nil
}

func fetchAndCachePackages(pwd string, urls []string) ([]string, error) {
	curLoc, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if err := os.Chdir(pwd); err != nil {
		return nil, err
	}
	defer os.Chdir(curLoc)

	var dirs []string
	for _, src := range urls {
		dst, err := fetchAndCachePackage(src)
		if err != nil {
			return dirs, err
		}
		dirs = append(dirs, dst)
	}
	return dirs, nil
}

// Fetches and caches a yardl package directory from src url
// NOTE: Currently only supports local file paths
func fetchAndCachePackage(src string) (string, error) {
	u, err := preprocessUrl(src)
	if err != nil {
		return u.String(), err
	}

	log.Printf("Fetching %s (%s)", src, u.String())

	switch u.Scheme {
	case "file":
		return u.Path, nil
	default:
		return u.Path, fmt.Errorf("scheme '%s' not yet supported", u.Scheme)
	}

	hash := md5.Sum([]byte(u.Path))
	dst := filepath.Join(cacheDir, hex.EncodeToString(hash[:]))

	if stat, err := os.Stat(dst); err == nil && stat.IsDir() {
		log.Printf("Already cached: %v -> %v", src, dst)
		return dst, nil
	}

	log.Printf("Fetching %s", u.String())
	if err := doFetch(u, dst); err != nil {
		return "", err
	}
	log.Printf("Cached in: %v", dst)
	return dst, nil
}

func doFetch(src *url.URL, dst string) error {
	panic("fetch not yet supported")
	return nil
}

func preprocessUrl(src string) (*url.URL, error) {
	u, err := url.Parse(src)
	if err != nil {
		return u, err
	}

	if u.Scheme == "" {
		u.Scheme = "file"
	}

	if u.Path == "" {
		return u, fmt.Errorf("invalid path '%s'", src)
	}

	abs, err := filepath.Abs(u.Path)
	if err != nil {
		return u, err
	}
	u.Path = abs

	return u, nil
}
