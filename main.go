/*
SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
SPDX-License-Identifier: MIT
*/
package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"run-go/widgets"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

const APP_DIR string = ".run-go"

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	appDir := fmt.Sprintf("%s/%s", home, APP_DIR)
	_, err = os.ReadDir(appDir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(appDir, 0755); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	snippetsDir := fmt.Sprintf("%s/%s", appDir, "snippets")
	_, err = os.ReadDir(snippetsDir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(snippetsDir, 0755); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	gosDir := fmt.Sprintf("%s/%s", appDir, "gos")
	_, err = os.ReadDir(gosDir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(gosDir, 0755); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	// Get the latest Go version
	goVersionURL := "https://go.dev/VERSION?m=text"
	resp, err := http.Get(goVersionURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, resp.Status)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	goVersion := strings.Split(string(body), "\n")[0]
	fmt.Println(goVersion)

	// Determine which os is the user running
	var goos, goarch string
	switch runtime.GOOS {
		case "darwin":
			goos = runtime.GOOS
			if runtime.GOARCH == "arm64" {
				goarch = runtime.GOARCH
			} else if runtime.GOARCH == "amd64" {
				goarch = runtime.GOARCH
			} else {
				fmt.Fprintln(os.Stderr, "invalid GOARCH")
				os.Exit(1)
			}
		case "linux":
			goos = runtime.GOOS
			if runtime.GOARCH == "amd64" {
				goarch = runtime.GOARCH
			} else {
				fmt.Fprintln(os.Stderr, "invalid GOARCH")
				os.Exit(1)
			}
		case "windows":
			goos = runtime.GOOS
			if runtime.GOARCH == "amd64" {
				goarch = runtime.GOARCH
			} else {
				fmt.Fprintln(os.Stderr, "invalid GOARCH")
				os.Exit(1)
			}
		default:
			fmt.Fprintln(os.Stderr, "invalid GOOS")
			os.Exit(1)
	}

	// Check if the latest version is already downloaded
	goLatestDir := fmt.Sprintf("%s/%s.%s-%s", gosDir, goVersion, goos, goarch)
	_, err = os.ReadDir(goLatestDir)
	if os.IsNotExist(err) {
		// Download latest Go version tarball in .run-go
		goDownloadURL := fmt.Sprintf("https://go.dev/dl/%s.%s-%s.tar.gz", goVersion, goos, goarch)
		resp, err = http.Get(goDownloadURL)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Fprintln(os.Stderr, resp.Status)
			os.Exit(1)
		}

		out, err := os.Create(fmt.Sprintf("%s/%s.%s-%s.tar.gz", appDir, goVersion, goos, goarch))
		if err != nil {
			fmt.Fprintln(os.Stderr, resp.Status)
			os.Exit(1)
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Untar the latest Go version tarball in .run-go/gos
		reader, err := os.Open(fmt.Sprintf("%s/%s.%s-%s.tar.gz", appDir, goVersion, goos, goarch))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer reader.Close()

		gzipReader, err := gzip.NewReader(reader)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer gzipReader.Close()

		tarReader := tar.NewReader(gzipReader)
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			target := filepath.Join(gosDir, header.Name)
			switch header.Typeflag {
				case tar.TypeDir:
					if _, err := os.Stat(target); err != nil {
						if err := os.MkdirAll(target, 0755); err != nil {
							fmt.Fprintln(os.Stderr, err)
							os.Exit(1)
						}
					}
				case tar.TypeReg:
					f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}

					if _, err := io.Copy(f, tarReader); err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}

					f.Close()
			}
		}

		// Rename the folder using the go version, os and architecture of the system
		if err := os.Rename(
			fmt.Sprintf("%s/%s", gosDir, "go"), 
			fmt.Sprintf("%s/%s", gosDir, 
				fmt.Sprintf("%s.%s-%s", goVersion, goos, goarch),
			),
		); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Remove the downloaded .tar.gz file
		if err := os.Remove(fmt.Sprintf("%s/%s.%s-%s.tar.gz", appDir, goVersion, goos, goarch)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	// Use the installed Go version
	if err := os.Setenv("GOPATH", fmt.Sprintf("%s/%s.%s-%s/bin/go", gosDir, goVersion, goos, goarch)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	ctrlT = &desktop.CustomShortcut{KeyName: fyne.KeyT, Modifier: fyne.KeyModifierControl}
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("RunGo")

	tabbar := widgets.NewTabBar(myWindow.Canvas())

	myWindow.Canvas().AddShortcut(ctrlT, tabbar.TypedShortcut)
	myWindow.Canvas().SetContent(tabbar.AppTabs)
	myWindow.Resize(fyne.NewSize(1024, 640))
	myWindow.ShowAndRun()
}
