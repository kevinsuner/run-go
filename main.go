/*
	SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
	SPDX-License-Identifier: MIT
*/
package main

import (
	"errors"
	"fmt"
	"io"
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
	"go.uber.org/zap"
)

const (
	APP_DIR			= ".run-go"
	SNIPPETS_DIR	= "snippets"
	GOS_DIR			= ".gos"

	ALT_T		= "CustomDesktop:Alt+T"
	ALT_S		= "CustomDesktop:Alt+S"
	ALT_O		= "CustomDesktop:Alt+O"
	ALT_RETURN	= "CustomDesktop:Alt+Return"

	GO_URL = "https://go.dev/dl"
)

var (
	errUnsupportedOS	= errors.New("unsupported operating system")
	errUnsupportedArch	= errors.New("unsupported processor architecture")
	errRequestFailed	= errors.New("failed to perform http request")
	errUnexpectedStatus = errors.New("received an unexpected http status code")
	errUnknownFileExt	= errors.New("file has an unknown extension")

	altT		= &desktop.CustomShortcut{KeyName: fyne.KeyT, Modifier: fyne.KeyModifierAlt}
	altS		= &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierAlt}
	altO		= &desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierAlt}
	altReturn	= &desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierAlt}

	// TODO: Find another way to do this, I feel like I'm overusing
	// global variables way too much
	osAndArchGoVersion	string
	bareGoVersion		string
	goFileExt			string
	goBinaryExt			string

	logger *zap.SugaredLogger
)

var aboutMD = `
RunGo is a free cross-platform Go playground, that allows users to experiment,
prototype and get instant feedback. It provides support for running Go versions
from 1.16+, and is built on top of [Fyne](https://fyne.io), a cross-platform GUI
toolkit made with Go and inspired by Material Design.

RunGo is mainly built using the following open-source projects:
- [Golang BSD 3-Clause License](https://github.com/golang/go/blob/master/LICENSE)
- [Fyne BSD 3-Clause License](https://github.com/fyne-io/fyne/blob/master/LICENSE)
- [goquery BSD 3-Clause License](https://github.com/PuerkitoBio/goquery/blob/master/LICENSE)

I've only included those who are a direct dependency of the project, but if you
wish to have a complete list of the projects being used, head to RunGo's **go.mod** file.

If you wish to hack for a bit and work on some open issues, you can do so by checking out
the **CONTRIBUTING.md** file.

[RunGo](https://github.com/itsksrof/run-go) is licensed under the **MIT License**

Copyright (c) 2023 Kevin Suñer
`

func init() {
	zapL, _ := zap.NewProduction()
	logger = zapL.Sugar()	

	home, err := os.UserHomeDir()
	if err != nil {
		logger.Fatalw("os.UserHomeDir()", "error", err.Error())
	}

	if err := os.Setenv("RUNGO_HOME", home); err != nil {
		logger.Fatalw("os.Setenv()", "error", err.Error())
	}

	appDir := fmt.Sprintf("%s/%s", home, APP_DIR)
	_, err = os.ReadDir(appDir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(appDir, 0755); err != nil {
			logger.Fatalw("os.Mkdir()", "error", err.Error())
		}
	}

	snippetsDir := fmt.Sprintf("%s/%s", appDir, SNIPPETS_DIR)
	_, err = os.ReadDir(snippetsDir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(snippetsDir, 0755); err != nil {
			logger.Fatalw("os.Mkdir()", "error", err.Error())
		}
	}

	gosDir := fmt.Sprintf("%s/%s", appDir, GOS_DIR)
	_, err = os.ReadDir(gosDir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(gosDir, 0755); err != nil {
			logger.Fatalw("os.Mkdir()", "error", err.Error())	
		}
	}

	if err := setOSAndArch(); err != nil {
		logger.Fatalw("setOSAndArch()", "error", err.Error())
	}

	// TODO: Shouldn't fail on this one, could be a network error
	// and default to the latest Go version that is installed
	osAndArchGoVersion, bareGoVersion, err = getLatestGoVersion()
	if err != nil {
		logger.Fatalw("checkLatestGoVersion()", "error", err.Error())
	}

	_, err = os.ReadDir(fmt.Sprintf("%s/%s", gosDir, osAndArchGoVersion))
	if os.IsNotExist(err) {
		// TODO: As the previous one, this shouldn't fail, and instead
		// default to the latest Go version that is installed
		if err := getGoFile(osAndArchGoVersion, goFileExt, appDir); err != nil {
			logger.Fatalw("getGoFile()", "error", err.Error())
		}

		if err := extractFile(
			fmt.Sprintf("%s/%s%s", appDir, osAndArchGoVersion, goFileExt),
			gosDir,
		); err != nil {
			logger.Fatalw("extractFile()", "error", err.Error())
		}

		if err := os.Rename(
			fmt.Sprintf("%s/%s", gosDir, "go"),
			fmt.Sprintf("%s/%s", gosDir, osAndArchGoVersion),
		); err != nil {
			logger.Fatalw("os.Rename()", "error", err.Error())
		}

		if err := os.Remove(fmt.Sprintf("%s/%s%s",
			appDir,
			osAndArchGoVersion,
			goFileExt,
		)); err != nil {
			logger.Fatalw("os.Remove()", "error", err.Error())
		}
	}

	if err := os.Setenv("RUNGO_GOBIN", fmt.Sprintf("%s/%s/bin/go%s",
		gosDir,
		osAndArchGoVersion,
		goBinaryExt,
	)); err != nil {
		logger.Fatalw("os.Setenv()", "error", err.Error())
	}
}

func setOSAndArch() error {
	var (
		osys, arch string
	)
	
	switch runtime.GOOS {
	case "darwin":
		osys = runtime.GOOS
		goFileExt = ".tar.gz"
		if runtime.GOARCH == "arm64" {
			arch = runtime.GOARCH
		} else if runtime.GOARCH == "amd64" {
			arch = runtime.GOARCH
		} else {
			return errUnsupportedArch
		}
	case "linux":
		osys = runtime.GOOS
		goFileExt = ".tar.gz"
		if runtime.GOARCH == "amd64" {
			arch = runtime.GOARCH
		} else {
			return errUnsupportedArch
		}
	case "windows":
		osys = runtime.GOOS
		goFileExt, goBinaryExt = ".zip", ".exe"
		if runtime.GOARCH == "amd64" {
			arch = runtime.GOARCH
		} else {
			return errUnsupportedArch
		}
	default:
		return errUnsupportedOS
	}

	if err := os.Setenv("RUNGO_OS", osys); err != nil {
		return err
	}

	if err := os.Setenv("RUNGO_ARCH", arch); err != nil {
		return err
	}

	return nil
}

func getLatestGoVersion() (string, string, error) {
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

func main() {
	var (
		shortcutsModal, aboutModal, versionModal *widget.PopUp
	)

	myApp := app.New()
	myWindow := myApp.NewWindow("RunGo")
	
	appTabs := newAppTabs(myWindow)

	shortcutsBtn := widget.NewButtonWithIcon("Shortcuts",
		theme.ContentRedoIcon(),
		func() {
			shortcutsModal.Resize(fyne.NewSize(440, 540))
			shortcutsModal.Show()
		},
	)

	aboutBtn := widget.NewButtonWithIcon("About RunGo",
		theme.InfoIcon(),
		func() {
			aboutModal.Resize(fyne.NewSize(440, 540))
			aboutModal.Show()
		},
	)

	versionBtn := widget.NewButtonWithIcon(bareGoVersion,
		theme.ConfirmIcon(),
		func() {
			versionModal.Resize(fyne.NewSize(440, 540))
			versionModal.Show()
		},
	)

	shortcutsModal = newShortcutsModal(myWindow.Canvas(), customShortcuts)
	aboutModal = newAboutModal(myWindow.Canvas(), aboutMD)
	versionModal = newVersionModal(myWindow, versionBtn, binding.NewString())
	
	myWindow.Canvas().AddShortcut(altT, appTabs.TypedShortcut)
	myWindow.SetContent(appLayout(appTabs.AppTabs, shortcutsBtn, aboutBtn, versionBtn))
	myWindow.Resize(fyne.NewSize(1280, 720))
	myWindow.ShowAndRun()
}

func appLayout(
	tabs *container.AppTabs, 
	shortcutsBtn, aboutBtn, versionBtn *widget.Button,
) *fyne.Container {
	return container.NewBorder(
		nil,
		container.NewPadded(
			container.NewGridWithColumns(8,
				shortcutsBtn,
				aboutBtn,
				layout.NewSpacer(),
				layout.NewSpacer(),
				layout.NewSpacer(),
				layout.NewSpacer(),
				layout.NewSpacer(),
				versionBtn,
			),
		),
		nil,
		nil,
		container.NewPadded(tabs),
	)
}
