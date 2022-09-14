package utils

import (
	"io"
	"net/http"
	"os"
)

func DownloadFile(file, url string) error {
	// Remove any existing file
	_ = os.Remove(file)
	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
