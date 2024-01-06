//go:build windows
// +build windows

package main

import (
	"syscall"
	"path/filepath"
)

func isHidden(f string) (bool, error) {
	ptr, err := syscall.UTF16PtrFromString(`\\?\` + f)
	if err != nil {
		return false, err
	}
	attrs, err := syscall.GetFileAttributes(ptr)
	if err != nil {
		return false, err
	}
	return attrs&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
}
