// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package packaging

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter/v2"
)

// Used to Detect file paths but not actually copy them to a new location.
type NoOpFileGetter struct {
	fg getter.FileGetter
}

const MaxImportRecursionDepth = 10

var cacheDir string

var client *getter.Client

func init() {
	if err := initCacheDir(); err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}

	initGetterClient()
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

func initGetterClient() {
	gitGetter := &getter.GitGetter{
		Detectors: []getter.Detector{
			new(getter.GitHubDetector),
			new(getter.GitDetector),
			new(getter.BitBucketDetector),
			new(getter.GitLabDetector),
		},
	}

	fileGetter := new(NoOpFileGetter)
	getters := []getter.Getter{gitGetter, fileGetter}

	client = &getter.Client{
		Getters:         getters,
		Decompressors:   nil,
		DisableSymlinks: true,
	}
}

func fetchAndCachePackages(pwd string, urls []string) ([]string, error) {
	var dirs []string
	for _, src := range urls {
		dst, err := fetchAndCachePackage(pwd, src)
		if err != nil {
			return dirs, err
		}
		dirs = append(dirs, dst)
	}
	return dirs, nil
}

// Fetches and caches a yardl schema package directory from url
// pwd is the directory of the current schema package
// src is the path to the dependency
func fetchAndCachePackage(pwd string, src string) (string, error) {
	hash := md5.Sum([]byte(src))
	dst := filepath.Join(cacheDir, hex.EncodeToString(hash[:]))

	if stat, err := os.Stat(dst); err == nil && stat.IsDir() {
		log.Printf("Already cached: %v -> %v", src, dst)
		return dst, nil
	}

	req := &getter.Request{
		Src:     src,
		Dst:     dst,
		Pwd:     pwd,
		GetMode: getter.ModeDir,
	}

	log.Printf("Fetching %v", req.Src)
	res, err := client.Get(context.Background(), req)
	if err != nil {
		return "", err
	}
	if dst == res.Dst {
		log.Printf("Cached in: %v", res.Dst)
	}
	return res.Dst, nil
}

func (n *NoOpFileGetter) Get(context.Context, *getter.Request) error {
	return nil
}

func (n *NoOpFileGetter) GetFile(context.Context, *getter.Request) error {
	return nil
}

func (n *NoOpFileGetter) Mode(ctx context.Context, u *url.URL) (getter.Mode, error) {
	return n.fg.Mode(ctx, u)
}

// Overwrites req.Dst, setting it to the absolute path of req.Src
func (n *NoOpFileGetter) Detect(req *getter.Request) (bool, error) {
	match, err := n.fg.Detect(req)
	if err != nil {
		return match, err
	}
	if match {
		dst, err := filepath.Abs(req.Src)
		if err != nil {
			return match, err
		}
		req.Dst = dst
	}
	return match, nil
}
