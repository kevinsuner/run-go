/*
	SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
	SPDX-License-Identifier: MIT
*/
package main

import (
	"errors"
	"fmt"
	"os"

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

	altT		= &desktop.CustomShortcut{KeyName: fyne.KeyT, Modifier: fyne.KeyModifierAlt}
	altS		= &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierAlt}
	altO		= &desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierAlt}
	altReturn	= &desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierAlt}

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

	osys, arch, err := getOSAndArch()
	if err != nil {
		logger.Fatalw("getOSAndArch()", "error", err.Error())
	}

	// TODO: Shouldn't fail on this one, could be a network error
	// and default to the latest Go version that is installed
	version, err := getLatestGoVersion()
	if err != nil {
		logger.Fatalw("getLatestGoVersion()", "error", err.Error())
	}

	_, err = os.ReadDir(fmt.Sprintf("%s/%s.%s-%s", gosDir, version, osys, arch))
	if os.IsNotExist(err) {
		// TODO: As the previous one, this shouldn't fail, and instead
		// default to the latest Go version that is installed
		if err := getGoSource(
			fmt.Sprintf("%s.%s-%s", version, osys, arch),
			osys,
			appDir,
		); err != nil {
			logger.Fatalw("getGoSource()", "error", err.Error())
		}

		if err := extractGoSource(
			fmt.Sprintf("%s.%s-%s", version, osys, arch),
			osys,
			appDir,
			gosDir,
		); err != nil {
			logger.Fatalw("extractGoSource()", "error", err.Error())
		}

		if err := os.Rename(
			fmt.Sprintf("%s/%s", gosDir, "go"),
			fmt.Sprintf("%s/%s.%s-%s", gosDir, version, osys, arch),
		); err != nil {
			logger.Fatalw("os.Rename()", "error", err.Error())
		}
	}

	if err := setEnvVariables(home, version, osys, arch); err != nil {
		logger.Fatalw("setEnvVariables()", "error", err.Error())
	}
}

func main() {
	var (
		version = os.Getenv("RUNGO_GO_VER")
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

	versionBtn := widget.NewButtonWithIcon(version,
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
