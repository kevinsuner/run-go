/*
	SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
	SPDX-License-Identifier: MIT
*/
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const APP_DIR string = ".run-go"
const SNIPPETS_DIR string = "snippets"
const GOS_DIR string = ".gos"

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	if err := os.Setenv("RUNGO_HOME", home); err != nil {
		log.Fatalln(err)
	}

	appDir := fmt.Sprintf("%s/%s", home, APP_DIR)
	_, err = os.ReadDir(appDir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(appDir, 0755); err != nil {
			log.Fatalln(err)
		}
	}

	snippetsDir := fmt.Sprintf("%s/%s", appDir, SNIPPETS_DIR)
	_, err = os.ReadDir(snippetsDir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(snippetsDir, 0755); err != nil {
			log.Fatalln(err)
		}
	}

	gosDir := fmt.Sprintf("%s/%s", appDir, GOS_DIR)
	_, err = os.ReadDir(gosDir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(gosDir, 0755); err != nil {
			log.Fatalln(err)
		}
	}

	if err := getOSAndArch(); err != nil {
		log.Fatalln(err)
	}

	goVersion, err := checkLatestGoVersion()
	if err != nil {
		log.Fatalln(err)
	}

	_, err = os.ReadDir(fmt.Sprintf("%s/%s", gosDir, goVersion))
	if os.IsNotExist(err) {
		if err := downloadGoTarball(goVersion, appDir); err != nil {
			log.Fatalln(err)
		}

		if err := untar(
			fmt.Sprintf("%s/%s.tar.gz", appDir, goVersion),
			gosDir,
		); err != nil {
			log.Fatalln(err)
		}

		if err := os.Rename(
			fmt.Sprintf("%s/%s", gosDir, "go"),
			fmt.Sprintf("%s/%s", gosDir, goVersion),
		); err != nil {
			log.Fatalln(err)
		}

		if err := os.Remove(fmt.Sprintf("%s/%s.tar.gz",
			appDir,
			goVersion,
		)); err != nil {
			log.Fatalln(err)
		}
	}

	if err := os.Setenv("RUNGO_GOBIN", fmt.Sprintf("%s/%s/bin/go",
		gosDir,
		goVersion,
	)); err != nil {
		log.Fatalln(err)
	}
}

func getOSAndArch() error {
	var goos, goarch string

	unsupportedOS := errors.New("unsupported operating system")
	unsupportedArch := errors.New("unsupported processor architecture")
	
	switch runtime.GOOS {
	case "darwin":
		goos = runtime.GOOS
		if runtime.GOARCH == "arm64" {
			goarch = runtime.GOARCH
		} else if runtime.GOARCH == "amd64" {
			goarch = runtime.GOARCH
		} else {
			return unsupportedArch
		}
	case "linux":
		goos = runtime.GOOS
		if runtime.GOARCH == "amd64" {
			goarch = runtime.GOARCH
		} else {
			return unsupportedArch
		}
	case "windows":
		goos = runtime.GOOS
		if runtime.GOARCH == "amd64" {
			goarch = runtime.GOARCH
		} else {
			return unsupportedArch
		}
	default:
		return unsupportedOS
	}

	if err := os.Setenv("RUNGO_OS", goos); err != nil {
		return err
	}

	if err := os.Setenv("RUNGO_ARCH", goarch); err != nil {
		return err
	}

	return nil
}

func checkLatestGoVersion() (string, error) {
	resp, err := http.Get("https://go.dev/VERSION?m=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s-%s",
		strings.Split(string(body), "\n")[0],
		os.Getenv("RUNGO_OS"),
		os.Getenv("RUNGO_ARCH"),
	), nil	
}

func desktopLayout(
	tabs *container.AppTabs, 
	shortcuts, about, version *widget.Button,
) *fyne.Container {
	return container.NewBorder(
		nil,
		container.NewPadded(container.NewGridWithColumns(3,
			container.NewGridWithColumns(2,
				shortcuts,
				about,
			),
			layout.NewSpacer(),
			version,
		)),
		nil,
		nil,
		container.NewPadded(tabs),
	)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("RunGo")
	
	tabs := container.NewAppTabs(
		container.NewTabItem("Tab 1", container.NewGridWithColumns(2,
			widget.NewEntry(),
			widget.NewLabel("Output"),
		)),
		container.NewTabItem("Tab 2", container.NewGridWithColumns(2,
			widget.NewEntry(),
			widget.NewLabel("Output"),
		)),
	)

	goVersionPopUp := goVersionPopUp(myWindow.Canvas())

	shortcuts := widget.NewButton("Shortcuts", nil)
	about := widget.NewButton("About RunGo", nil)
	version := widget.NewButton("Go 1.21.4", func() {
		goVersionPopUp.Resize(fyne.NewSize(440, 540))
		goVersionPopUp.Show()
	})

	myWindow.SetContent(desktopLayout(tabs, shortcuts, about, version))
	myWindow.Resize(fyne.NewSize(1280, 720))
	myWindow.ShowAndRun()
}
