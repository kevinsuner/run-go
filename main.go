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
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const APP_DIR string = ".run-go"
const SNIPPETS_DIR string = "snippets"
const GOS_DIR string = ".gos"

const ALT_T = "CustomDesktop:Alt+T"
const ALT_S = "CustomDesktop:Alt+S"
const ALT_O = "CustomDesktop:Alt+O"
const ALT_RETURN = "CustomDesktop:Alt+Return"

var altT = &desktop.CustomShortcut{KeyName: fyne.KeyT, Modifier: fyne.KeyModifierAlt}
var altS = &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierAlt}
var altO = &desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierAlt}
var altReturn = &desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierAlt}

var osAndArchGoVersion string
var bareGoVersion string

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

	osAndArchGoVersion, bareGoVersion, err = checkLatestGoVersion()
	if err != nil {
		log.Fatalln(err)
	}

	_, err = os.ReadDir(fmt.Sprintf("%s/%s", gosDir, osAndArchGoVersion))
	if os.IsNotExist(err) {
		if err := downloadGoTarball(osAndArchGoVersion, appDir); err != nil {
			log.Fatalln(err)
		}

		if err := untar(
			fmt.Sprintf("%s/%s.tar.gz", appDir, osAndArchGoVersion),
			gosDir,
		); err != nil {
			log.Fatalln(err)
		}

		if err := os.Rename(
			fmt.Sprintf("%s/%s", gosDir, "go"),
			fmt.Sprintf("%s/%s", gosDir, osAndArchGoVersion),
		); err != nil {
			log.Fatalln(err)
		}

		if err := os.Remove(fmt.Sprintf("%s/%s.tar.gz",
			appDir,
			osAndArchGoVersion,
		)); err != nil {
			log.Fatalln(err)
		}
	}

	if err := os.Setenv("RUNGO_GOBIN", fmt.Sprintf("%s/%s/bin/go",
		gosDir,
		osAndArchGoVersion,
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

func checkLatestGoVersion() (string, string, error) {
	resp, err := http.Get("https://go.dev/VERSION?m=text")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	version := strings.Split(string(body), "\n")[0]
	return fmt.Sprintf("%s.%s-%s",
		version,
		os.Getenv("RUNGO_OS"),
		os.Getenv("RUNGO_ARCH"),
	), version, nil
}

func desktopLayout(
	tabs *container.AppTabs, 
	shortcuts, about, version *widget.Button,
) *fyne.Container {
	return container.NewBorder(
		nil,
		container.NewPadded(
			container.NewGridWithColumns(8,
				shortcuts,
				about,
				layout.NewSpacer(),
				layout.NewSpacer(),
				layout.NewSpacer(),
				layout.NewSpacer(),
				layout.NewSpacer(),
				version,
			),
		),
		nil,
		nil,
		container.NewPadded(tabs),
	)
}

func main() {
	var versionPopUp *widget.PopUp

	myApp := app.New()
	myWindow := myApp.NewWindow("RunGo")
	
	tabs := appTabs(myWindow.Canvas())

	goVersionStr := binding.NewString()
	goVersion := widget.NewButtonWithIcon(
		bareGoVersion, 
		theme.ConfirmIcon(), 
		func() {
			versionPopUp.Resize(fyne.NewSize(440, 540))
			versionPopUp.Show()
		},
	)
	versionPopUp = goVersionPopUp(
		myWindow.Canvas(), 
		goVersion, 
		goVersionStr,
	)

	shortcuts := widget.NewButtonWithIcon("Shortcuts", theme.ContentRedoIcon(), func() {})
	about := widget.NewButtonWithIcon("About RunGo", theme.InfoIcon(), func() {})

	myWindow.Canvas().AddShortcut(altT, tabs.TypedShortcut)
	myWindow.SetContent(desktopLayout(tabs.AppTabs, shortcuts, about, goVersion))
	myWindow.Resize(fyne.NewSize(1280, 720))
	myWindow.ShowAndRun()
}
