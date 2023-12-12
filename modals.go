/*
	SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
	SPDX-License-Identifier: MIT
*/
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
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
	var shortcutsModal *widget.PopUp
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
	var aboutModal *widget.PopUp
	mdText := widget.NewRichTextFromMarkdown(content)
	mdText.Wrapping = fyne.TextWrap(fyne.TextWrapBreak)

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
		container.NewPadded(mdText),
	), canvas)

	return aboutModal
}

func newVersionModal(window fyne.Window, versionBtn *widget.Button, versionStr binding.String) *widget.PopUp {
	versions, err := getGoVersions()
	if err != nil {
		if errors.Unwrap(err) == errRequestFailed || errors.Unwrap(err) == errUnexpectedStatus {
			dialog.NewInformation("An error occurred", err.Error(), window).Show()
			logger.Error("getGoVersions()", zap.Error(err))

			// Set versions to an empty array, as this could be due to the
			// service being currently unavailable
			versions = []string{}
		} else {
			logger.Fatal("getGoVersions()", zap.Error(err))
		}
	}

	var versionModal *widget.PopUp
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
					longVersion := fmt.Sprintf("%s.%s-%s", button.Text, runtime.GOOS, runtime.GOARCH)
					_, err = os.ReadDir(filepath.Join(os.Getenv("RUNGO_APP_DIR"), GOS_DIR, longVersion))
					if os.IsNotExist(err) {
						progress := dialog.NewCustomWithoutButtons(fmt.Sprintf("Downloading %s", button.Text),
							container.NewPadded(widget.NewProgressBarInfinite()),
							window,
						)
						progress.Show()
						defer progress.Hide()

						switch runtime.GOOS {
						case "windows":
							err = getGoSource(fmt.Sprintf("%s.%s", longVersion, "zip"), filepath.Join(os.Getenv("RUNGO_APP_DIR")))
							if err != nil {
								logger.Fatal("getGoSource()", zap.Error(err))
							}

							err = uncompressZipFile(
								filepath.Join(os.Getenv("RUNGO_APP_DIR"), fmt.Sprintf("%s.%s", longVersion, "zip")),
								filepath.Join(os.Getenv("RUNGO_APP_DIR"), GOS_DIR),
							)
							if err != nil {
								logger.Fatal("uncompressZipFile()", zap.Error(err))
							}
						default:
							err = getGoSource(fmt.Sprintf("%s.%s", longVersion, "tar.gz"), filepath.Join(os.Getenv("RUNGO_APP_DIR")))
							if err != nil {
								logger.Fatal("getGoSource()", zap.Error(err))
							}

							err = uncompressTarFile(
								filepath.Join(os.Getenv("RUNGO_APP_DIR"), fmt.Sprintf("%s.%s", longVersion, "tar.gz")),
								filepath.Join(os.Getenv("RUNGO_APP_DIR"), GOS_DIR),
							)
							if err != nil {
								logger.Fatal("uncompressTarFile()", zap.Error(err))
							}
						}

						err = os.Rename(filepath.Join(os.Getenv("RUNGO_APP_DIR"), GOS_DIR, "go"), filepath.Join(os.Getenv("RUNGO_APP_DIR"), GOS_DIR, longVersion))
						if err != nil {
							logger.Fatal("os.Rename()", zap.Error(err))
						}
					}


					err = os.Setenv("RUNGO_GO_BIN", filepath.Join(os.Getenv("RUNGO_APP_DIR"), GOS_DIR, longVersion, "bin", "go"))
					if err != nil {
						logger.Fatal("os.Setenv()", zap.Error(err))
					}

					if runtime.GOOS == "windows" {
						err = os.Setenv("RUNGO_GO_BIN", filepath.Join(os.Getenv("RUNGO_APP_DIR"), GOS_DIR, longVersion, "bin", "go.exe"))
						if err != nil {
							logger.Fatal("os.Setenv()", zap.Error(err))
						}
					}

					err = versionStr.Set(versions[lid])
					if err != nil {
						logger.Fatal("versionStr.Set()", zap.Error(err))
					}

					version, err := versionStr.Get()
					if err != nil {
						logger.Fatal("versionStr.Get()", zap.Error(err))
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

func newSaveModal(entry *widget.Entry, appTabs *container.AppTabs, snippet binding.String, window fyne.Window) *customSaveModal {
	customSaveModal := &customSaveModal{}
	
	input := &widget.Entry{PlaceHolder: "Snippet name"}
	var saveModal *widget.PopUp
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
				err := newSnippet(input.Text, []byte(entry.Text))
				if err != nil {
					if errors.Is(err, os.ErrExist) {
						dialog.NewInformation("An error occurred", err.Error(), window).Show()
						logger.Error("newSnippet()", zap.Error(err))
						return
					} else {
						logger.Fatal("newSnippet()", zap.Error(err))
					}
				}

				err = snippet.Set(input.Text)
				if err != nil {
					logger.Fatal("snippet.Set()", zap.Error(err))
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

func newOpenModal(entry *widget.Entry, appTabs *container.AppTabs, snippet binding.String, snippetList binding.StringList, window fyne.Window) *customOpenModal {
	customOpenModal := &customOpenModal{snippetList: snippetList}
	
	var openModal *widget.PopUp
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
				openModal.Hide()
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
					logger.Fatal("snippetList.GetValue()", zap.Error(err))
				}

				button := obj.(*widget.Button)
				button.SetText(snippetName)
				button.Alignment = widget.ButtonAlignLeading
				button.OnTapped = func() {
					if len(entry.Text) != 0 {
						dialog.NewInformation("Info", "Tab already in use", window).Show()
						logger.Warn("user attempted to open snippet in used tab")
						return
					}

					data, err := openSnippet(button.Text)
					if err != nil {
						if errors.Is(err, os.ErrNotExist) {
							dialog.NewInformation("An error occurred", err.Error(), window).Show()
							logger.Error("openSnippet()", zap.Error(err))
							return
						} else {
							logger.Fatal("openSnippet()", zap.Error(err))
						}
					}

					err = snippet.Set(button.Text)
					if err != nil {
						logger.Fatal("snippet.Set()", zap.Error(err))
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
			logger.Fatal("listSnippet()", zap.Error(err))
		}

		err = c.snippetList.Set(snippets)
		if err != nil {
			logger.Fatal("c.snippetList.Set()", zap.Error(err))
		}

		c.PopUp.Resize(fyne.NewSize(440, 540))
		c.PopUp.Show()
	}
}
