// Package main loads optional .env for PAYARA_APP_ID / PAYARA_APP_SECRET.
package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func loadEnv() {
	for _, path := range []string{".env", ".env.local"} {
		if tryLoadEnv(path) {
			return
		}
		if exe, err := os.Executable(); err == nil {
			dir := filepath.Dir(exe)
			if tryLoadEnv(filepath.Join(dir, path)) || tryLoadEnv(filepath.Join(dir, "..", path)) {
				return
			}
		}
	}
}

func tryLoadEnv(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, "=")
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, `"`)
		if key != "" && os.Getenv(key) == "" {
			_ = os.Setenv(key, val)
		}
	}
	return true
}
