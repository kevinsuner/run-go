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
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type customShortcut struct {
	keys string
	info string
}

var customShortcuts = []customShortcut{
	{keys: "Alt+T", info: "Open a new tab"},
	{keys: "Alt+S", info: "Open save snippet modal"},
	{keys: "Alt+O", info: "Open load snippet modal"},
	{keys: "Alt+Return", info: "Run code"},
}

func newShortcutsModal(canvas fyne.Canvas, shortcuts []customShortcut) *widget.PopUp {
	var (
		shortcutsModal *widget.PopUp
	)

	shortcutsTable := widget.NewTable(
		func() (int, int) {
			// N rows by 2 cols
			return len(shortcuts), 2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(tid widget.TableCellID, obj fyne.CanvasObject) {
			switch tid.Col {
			case 0:
				obj.(*widget.Label).SetText(shortcuts[tid.Row].keys)
			case 1:
				obj.(*widget.Label).SetText(shortcuts[tid.Row].info)
			}
		},
	)

	// TODO: Need to check if there is a better way to set this up
	shortcutsTable.SetColumnWidth(0, 200)
	shortcutsTable.SetColumnWidth(1, 200)

	shortcutsModal = widget.NewModalPopUp(container.NewBorder(
		container.NewPadded(container.NewGridWithColumns(12,
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				shortcutsModal.Hide()
			}),
		)),
		nil,
		nil,
		nil,
		container.NewPadded(shortcutsTable),
	), canvas)

	return shortcutsModal
}

func newAboutModal(canvas fyne.Canvas, content string) *widget.PopUp {
	var (
		aboutModal *widget.PopUp
	)

	aboutModal = widget.NewModalPopUp(container.NewBorder(
		container.NewPadded(container.NewGridWithColumns(12,
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				aboutModal.Hide()
			}),
		)),
		nil,
		nil,
		nil,
		container.NewPadded(widget.NewRichTextFromMarkdown(content)),
	), canvas)

	return aboutModal
}

func newVersionModal(
	window fyne.Window, 
	versionBtn *widget.Button, 
	versionStr binding.String,
) *widget.PopUp {
	var (
		versionModal *widget.PopUp
		home = os.Getenv("RUNGO_HOME")
		osys = os.Getenv("RUNGO_OS")
		arch = os.Getenv("RUNGO_ARCH")
	)

	versions, err := getGoVersions()
	if err != nil {
		if errors.Unwrap(err) == errRequestFailed ||
			errors.Unwrap(err) == errUnexpectedStatus {
			dialog.NewInformation("An error occurred", err.Error(), window).Show()
			logger.Errorw("getGoVersions()", "error", err.Error())

			// Set versions to an empty array, as this could be due to the
			// service being currently unavailable
			versions = []string{}
		} else {
			logger.Fatalw("getGoVersions()", "error", err.Error())
		}
	}

	versionModal = widget.NewModalPopUp(container.NewBorder(
		container.NewPadded(container.NewGridWithColumns(12,
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				versionModal.Hide()
			}),
		)),
		nil,
		nil,
		nil,
		container.NewPadded(widget.NewList(
			func() int {
				return len(versions)
			},
			func() fyne.CanvasObject {
				return widget.NewButton("template", nil)
			},
			func(lid widget.ListItemID, obj fyne.CanvasObject) {
				button := obj.(*widget.Button)
				button.SetText(versions[lid])
				button.Alignment = widget.ButtonAlignLeading
				button.OnTapped = func() {
					appDir := fmt.Sprintf("%s/%s", home, APP_DIR)
					gosDir := fmt.Sprintf("%s/%s", appDir, GOS_DIR)
					version := fmt.Sprintf("%s.%s-%s", button.Text, osys, arch)

					_, err := os.ReadDir(fmt.Sprintf("%s/%s", gosDir, version))
					if os.IsNotExist(err) {
						// NOTE: Ideally this should represent a real progress bar
						progress := dialog.NewCustomWithoutButtons(
							fmt.Sprintf("Downloading %s", version),
							container.NewPadded(widget.NewProgressBarInfinite()),
							window,
						)
						progress.Show()
						defer progress.Hide()

						if err := getGoTarball(version, appDir); err != nil {
							logger.Fatalw("getGoTarball()", "error", err.Error())
						}

						if err := untarFile(
							fmt.Sprintf("%s/%s.tar.gz", appDir, version),
							gosDir,
						); err != nil {
							logger.Fatalw("untar()", "error", err.Error())
						}

						if err := os.Rename(
							fmt.Sprintf("%s/%s", gosDir, "go"),
							fmt.Sprintf("%s/%s", gosDir, version),
						); err != nil {
							logger.Fatalw("os.Rename()", "error", err.Error())
						}

						if err := os.Remove(
							fmt.Sprintf("%s/%s.tar.gz", appDir, version),
						); err != nil {
							logger.Fatalw("os.Remove()", "error", err.Error())
						}

					}

					if err := os.Setenv("RUNGO_GOBIN",
						fmt.Sprintf("%s/%s/bin/go", gosDir, version),
					); err != nil {
						logger.Fatalw("os.Setenv()", "error", err.Error())
					}

					if err := versionStr.Set(versions[lid]); err != nil {
						logger.Fatalw("versionStr.Set()", "error", err.Error())
					}

					version, err = versionStr.Get()
					if err != nil {
						logger.Fatalw("versionStr.Get()", "error", err.Error())
					}

					versionBtn.SetText(version)
					versionModal.Hide()
				}
			},
		)),
	), window.Canvas())

	return versionModal
}

type customSaveModal struct {
	*widget.PopUp
}

func newSaveModal(
	entry *widget.Entry, 
	appTabs *container.AppTabs,
	snippet binding.String,
	window fyne.Window,
) *customSaveModal {
	var (
		customSaveModal = &customSaveModal{}
		saveModal *widget.PopUp
	)

	input := &widget.Entry{PlaceHolder: "Snippet name"}
	saveModal = widget.NewModalPopUp(container.NewBorder(
		container.NewPadded(container.NewGridWithColumns(12,
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				saveModal.Hide()
			}),
		)),
		nil,
		nil,
		nil,
		container.NewPadded(container.NewVBox(
			input,
			widget.NewButtonWithIcon("Save", theme.ConfirmIcon(), func() {
				if err := newSnippet(input.Text, []byte(entry.Text)); err != nil {
					if errors.Is(err, os.ErrExist) {
						dialog.NewInformation("An error occurred", err.Error(), window)
						logger.Errorw("newSnippet()", "error", err.Error())
					} else {
						logger.Fatalw("newSnippet()", "error", err.Error())
					}
				}

				if err := snippet.Set(input.Text); err != nil {
					logger.Fatalw("snippet.Set()", "error", err.Error())
				}

				appTabs.Selected().Text = input.Text
				appTabs.Refresh()
				saveModal.Hide()
			}),
		)),
	), window.Canvas())

	customSaveModal.PopUp = saveModal
	return customSaveModal
}

func (c *customSaveModal) TypedShortcut(shortcut fyne.Shortcut) {
	customShortcut, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		c.TypedShortcut(shortcut)
		return
	}

	switch customShortcut.ShortcutName() {
	case ALT_S:
		c.PopUp.Resize(fyne.NewSize(440, 200))
		c.PopUp.Show()
	}
}

type customOpenModal struct {
	snippetList binding.StringList
	*widget.PopUp
}

func newOpenModal(
	entry *widget.Entry,
	appTabs *container.AppTabs,
	snippet binding.String,
	snippetList binding.StringList,
	window fyne.Window,
) *customOpenModal {
	var (
		customOpenModal = &customOpenModal{snippetList: snippetList}
		openModal *widget.PopUp
	)

	openModal = widget.NewModalPopUp(container.NewBorder(
		container.NewPadded(container.NewGridWithColumns(12,
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				openModal.Show()
			}),
		)),
		nil,
		nil,
		nil,
		container.NewPadded(widget.NewList(
			func() int {
				return snippetList.Length()
			},
			func() fyne.CanvasObject {
				return widget.NewButton("template", nil)
			},
			func(lid widget.ListItemID, obj fyne.CanvasObject) {
				snippetName, err := snippetList.GetValue(lid)
				if err != nil {
					logger.Fatalw("snippetList.GetValue()", "error", err.Error())
				}

				button := obj.(*widget.Button)
				button.SetText(snippetName)
				button.Alignment = widget.ButtonAlignLeading
				button.OnTapped = func() {
					if len(entry.Text) != 0 {
						dialog.NewInformation("Info", "Tab already in use", window)
						logger.Warnw("user attempted to open snippet in used tab")
						openModal.Hide()
						return
					}

					data, err := openSnippet(button.Text)
					if err != nil {
						if errors.Is(err, os.ErrNotExist) {
							dialog.NewInformation("An error occurred", err.Error(), window)
							logger.Errorw("openSnippet()", "error", err.Error())
						} else {
							logger.Fatalw("openSnippet()", "error", err.Error())
						}
					}

					if err := snippet.Set(button.Text); err != nil {
						logger.Fatalw("snippet.Set()", "error", err.Error())
					}

					entry.SetText(data)
					appTabs.Selected().Text = button.Text
					appTabs.Refresh()
					openModal.Hide()
				}
			},
		)),
	), window.Canvas())

	customOpenModal.PopUp = openModal
	return customOpenModal
}

func (c *customOpenModal) TypedShortcut(shortcut fyne.Shortcut) {
	customShortcut, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		c.TypedShortcut(shortcut)
		return
	}

	switch customShortcut.ShortcutName() {
	case ALT_O:
		snippets, err := listSnippets()
		if err != nil {
			logger.Fatalw("listSnippet()", "error", err.Error())
		}

		if err := c.snippetList.Set(snippets); err != nil {
			logger.Fatalw("c.snippetList.Set()", "error", err.Error())
		}

		c.PopUp.Resize(fyne.NewSize(440, 540))
		c.PopUp.Show()
	}
}
