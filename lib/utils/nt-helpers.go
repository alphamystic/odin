//go:build windows
// +build windows

package utils

import (
  "os"
)

// stollen from https://github.com/mauri870/ransomware
func GetDrives() (letters []string) {
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		_, err := os.Open(string(drive) + ":\\")
		if err == nil {
			letters = append(letters, string(drive)+":\\")
		}
	}
	return
}
