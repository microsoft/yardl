// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package packaging

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	ParamRef = "ref"
	ParamDir = "dir"
)

var cacheDir string

func init() {
	if err := initCacheDir(); err != nil {
		log.Fatal().Msgf("Failed to initialize cache: %v", err)
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
// src can be a local file path for URL to a git repository
func fetchAndCachePackage(src string) (string, error) {
	u, err := url.Parse(src)
	if err != nil {
		return "", err
	}

	if u.Scheme == "" {
		u.Scheme = "file"
	}

	if u.Path == "" {
		return u.String(), fmt.Errorf("invalid path '%s'", src)
	}

	log.Info().Msgf("Fetching %s ", u.String())

	switch u.Scheme {
	case "file":
		abs, err := filepath.Abs(u.Path)
		if err != nil {
			return u.Path, err
		}
		return abs, nil
	case "git", "https":
		return fetchGit(u)
	default:
		return u.Path, fmt.Errorf("scheme '%s' not yet supported", u.Scheme)
	}
}

func fetchGit(url *url.URL) (string, error) {
	q := url.Query()

	ref := ""
	dir := ""
	if q.Has(ParamRef) {
		ref = q.Get(ParamRef)
		q.Del(ParamRef)
		url.RawQuery = q.Encode()
	}
	if q.Has(ParamDir) {
		dir = q.Get(ParamDir)
		q.Del(ParamDir)
		url.RawQuery = q.Encode()
	}

	dst := filepath.Join(cacheDir, path.Join(url.Hostname(), url.EscapedPath()))

	if ref == "" {
		ref = "remotes/origin/HEAD"
		dst = filepath.Join(dst, "HEAD")
	} else {
		dst = filepath.Join(dst, ref)
	}

	if stat, err := os.Stat(dst); err == nil {
		if !stat.IsDir() {
			return dst, fmt.Errorf("cache target '%s' is not a directory", dst)
		}

		log.Info().Msgf("Updating cached repo %v in %v", url, dst)
		if err := runGit("-C", dst, "fetch", "--all"); err != nil {
			return dst, err
		}
	} else if os.IsNotExist(err) {
		log.Info().Msgf("Cloning %s into %s", url, dst)
		if err := runGit("clone", url.String(), dst); err != nil {
			return dst, err
		}
	} else {
		return dst, err
	}

	log.Info().Msgf("Checking out ref %s", ref)
	if err := runGit("-C", dst, "checkout", ref); err != nil {
		return dst, err
	}

	if dir != "" {
		dst = filepath.Join(dst, dir)
		stat, err := os.Stat(dst)
		if (err != nil && os.IsNotExist(err)) || !stat.IsDir() {
			return dst, fmt.Errorf("git repository '%s' does not contain a dir named '%s'", url, dir)
		}
		if err != nil {
			return dst, err
		}
	}

	log.Info().Msgf("Cached %s", dst)
	return dst, nil
}

func runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	var capture strings.Builder
	cmd.Stdout = &capture
	cmd.Stderr = &capture

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed with error: %w\n\tcommand: %s\n\toutput: %s", err, cmd, capture.String())
	}

	return nil
}
