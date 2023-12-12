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
	"path"
	"path/filepath"
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
	APP_DIR			= "run-go"
	GOS_DIR			= "gos"
	SNIPPETS_DIR	= "snippets"

	ALT_T		= "CustomDesktop:Alt+T"
	ALT_S		= "CustomDesktop:Alt+S"
	ALT_O		= "CustomDesktop:Alt+O"
	ALT_RETURN	= "CustomDesktop:Alt+Return"

	GO_URL = "https://go.dev/dl"
)

var (
	errRequestFailed	= errors.New("failed to perform http request")
	errUnexpectedStatus = errors.New("received an unexpected http status code")

	altT		= &desktop.CustomShortcut{KeyName: fyne.KeyT, Modifier: fyne.KeyModifierAlt}
	altS		= &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierAlt}
	altO		= &desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierAlt}
	altReturn	= &desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierAlt}

	logger *zap.Logger
)

var aboutMD = `
RunGo is a free cross-platform Go playground that allows users to experiment,
prototype and get instant feedback. It provides support for running Go versions
from 1.16+, and is built on top of [Fyne](https://fyne.io), a cross-platform GUI
toolkit made with Go and inspired by Material Design.

RunGo is mainly built using the following open-source projects:
- [github.com/golang/go](https://github.com/golang/go)
- [github.com/fyne-io/fyne](https://github.com/fyne-io/fyne)
- [github.com/PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)
- [github.com/golang/mod](https://github.com/golang/mod)
- [github.com/fyne-io/fyne-cross](https://github.com/fyne-io/fyne-cross)

If you wish to hack for a bit and work on some open issues, you can do so by checking out
the [Contributing](https://github.com/itsksrof/run-go/tree/master#contributing) section.

[RunGo](https://github.com/itsksrof/run-go) is licensed under the **MIT License**

Copyright (c) 2023 Kevin Suñer
`

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	zapLogger := zap.NewProductionConfig()
	zapLogger.OutputPaths = []string{filepath.Join(homeDir, APP_DIR, "run-go.log")}
	logger, err = zapLogger.Build()
	if err != nil {
		log.Fatalln(err)
	}

	appDirs := []string{
		filepath.Join(homeDir, APP_DIR),
		filepath.Join(homeDir, APP_DIR, GOS_DIR),
		filepath.Join(homeDir, APP_DIR, SNIPPETS_DIR),
	}
	
	for _, appDir := range appDirs {
		_, err = os.ReadDir(appDir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(appDir, 0755)
			if err != nil {
				logger.Fatal("os.MkdirAll()", zap.Error(err))
			}
		}
	}

	getLatestGoVersion := func() (string, error) {
		res, err := http.Get(path.Join(GO_URL, "VERSION?m=text"))
		if err != nil {
			return "", fmt.Errorf("%w: %v", errRequestFailed, err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return "", fmt.Errorf("%w: %s", errUnexpectedStatus, res.Status)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		return strings.Split(string(body), "\n")[0], nil
	}

	version, err := getLatestGoVersion()
	if err != nil {
		logger.Fatal("getLatestGoVersion()", zap.Error(err))
	}

	longVersion := fmt.Sprintf("%s.%s-%s", version, runtime.GOOS, runtime.GOARCH)
	_, err = os.ReadDir(filepath.Join(homeDir, APP_DIR, GOS_DIR, longVersion))
	if os.IsNotExist(err) {
		switch runtime.GOOS {
		case "windows":
			err = getGoSource(fmt.Sprintf("%s.%s", longVersion, "zip"), filepath.Join(homeDir, APP_DIR))
			if err != nil {
				logger.Fatal("getGoSource()", zap.Error(err))
			}
	
			err = uncompressZipFile(
				filepath.Join(homeDir, APP_DIR, fmt.Sprintf("%s.%s", longVersion, "zip")), 
				filepath.Join(homeDir, APP_DIR, GOS_DIR),
			)
			if err != nil {
				logger.Fatal("uncompressZipFile()", zap.Error(err))
			}
		default:
			err = getGoSource(fmt.Sprintf("%s.%s", longVersion, "tar.gz"), filepath.Join(homeDir, APP_DIR))
			if err != nil {
				logger.Fatal("getGoSource()", zap.Error(err))
			}

			err = uncompressTarFile(
				filepath.Join(homeDir, APP_DIR, fmt.Sprintf("%s.%s", longVersion, "tar.gz")),
				filepath.Join(homeDir, APP_DIR, GOS_DIR),
			)
			if err != nil {
				logger.Fatal("uncompressTarFile()", zap.Error(err))
			}
		}

		err = os.Rename(filepath.Join(homeDir, APP_DIR, GOS_DIR, "go"), filepath.Join(homeDir, APP_DIR, GOS_DIR, longVersion))
		if err != nil {
			logger.Fatal("os.Rename()", zap.Error(err))
		}
	}

	setEnvironment := func() error {
		appDirErr := os.Setenv("RUNGO_APP_DIR", filepath.Join(homeDir, APP_DIR))
		goVerErr := os.Setenv("RUNGO_GO_VER", version)
		goBinErr := os.Setenv("RUNGO_GO_BIN", filepath.Join(homeDir, APP_DIR, GOS_DIR, longVersion, "bin", "go"))
		if runtime.GOOS == "windows" {
			goBinErr = os.Setenv("RUNGO_GO_BIN", filepath.Join(homeDir, APP_DIR, GOS_DIR, longVersion, "bin", "go.exe"))
			return errors.Join(appDirErr, goVerErr, goBinErr)
		}

		return errors.Join(appDirErr, goVerErr, goBinErr)
	}

	err = setEnvironment()
	if err != nil {
		logger.Fatal("setEnvironment()", zap.Error(err))
	}
}

func appLayout(tabs *container.AppTabs, shortcutsBtn, aboutBtn, versionBtn *widget.Button) *fyne.Container {
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

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("RunGo")
	
	appTabs := newAppTabs(myWindow)

	var shortcutsModal, aboutModal, versionModal *widget.PopUp
	shortcutsBtn := widget.NewButtonWithIcon("Shortcuts", theme.ContentRedoIcon(), func() {
		shortcutsModal.Resize(fyne.NewSize(440, 540))
		shortcutsModal.Show()
	})
	aboutBtn := widget.NewButtonWithIcon("About RunGo", theme.InfoIcon(), func() {
		aboutModal.Resize(fyne.NewSize(440, 540))
		aboutModal.Show()
	})
	versionBtn := widget.NewButtonWithIcon(os.Getenv("RUNGO_GO_VER"), theme.ConfirmIcon(), func() {
		versionModal.Resize(fyne.NewSize(440, 540))
		versionModal.Show()
	})

	shortcutsModal = newShortcutsModal(myWindow.Canvas(), customShortcuts)
	aboutModal = newAboutModal(myWindow.Canvas(), aboutMD)
	versionModal = newVersionModal(myWindow, versionBtn, binding.NewString())
	
	myWindow.Canvas().AddShortcut(altT, appTabs.TypedShortcut)
	myWindow.SetContent(appLayout(appTabs.AppTabs, shortcutsBtn, aboutBtn, versionBtn))
	myWindow.Resize(fyne.NewSize(1280, 720))
	myWindow.ShowAndRun()
}

