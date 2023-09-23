package utils

import (
	"config-service/about"
	"os"
	"path/filepath"
)

func ExistFileOrDir(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func CheckExistSSL(path, pemFile, keyFile string) bool {
	ssldir := filepath.Join(path, about.SSL_Directory)
	pemPath := filepath.Join(ssldir, pemFile)
	if _, err := os.Stat(pemPath); err != nil {
		return false
	}

	keyPath := filepath.Join(ssldir, keyFile)
	if _, err := os.Stat(keyPath); err != nil {
		return false
	}
	return true
}
