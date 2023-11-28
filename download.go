package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const GO_DOWNLOADS_URL string = "https://go.dev/dl"

func downloadGoTarball(version, dst string) error {
	downloadURL := fmt.Sprintf("%s/%s.tar.gz", GO_DOWNLOADS_URL, version)
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	out, err := os.Create(
		fmt.Sprintf("%s/%s.tar.gz", dst, version),
	)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
