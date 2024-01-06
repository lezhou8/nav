//go:build !windows
// +build !windows

package main

import (
	"strings"
	"path/filepath"
)

func isHidden(f string) (bool, error) {
	return strings.HasPrefix(filepath.Base(f), "."), nil
}
