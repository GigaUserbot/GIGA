package utils

import (
	"errors"
	"runtime"
)

func GetSupportedOS() (string, error) {
	switch runtime.GOOS {
	case "windows", "darwin", "linux":
		return runtime.GOOS, nil
	default:
		return "", errors.New("unsupported OS")
	}
}

func GetSupportedARCH() (string, error) {
	switch runtime.GOARCH {
	case "amd64", "arm64":
		return runtime.GOARCH, nil
	default:
		return "", errors.New("unsupported ARCH")
	}
}
