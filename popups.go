package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/mod/semver"
)

type shortcuts struct {
	shortcut string
	info string
}

func shortcutsPopUp(canvas fyne.Canvas) *widget.PopUp {
	shortcuts := []shortcuts{
		{shortcut: "Alt+T", info: "Open a new tab"},
		{shortcut: "Alt+S", info: "Open save snippet modal"},
		{shortcut: "Alt+O", info: "Open load snippet modal"},
		{shortcut: "Alt+Return", info: "Run code from editor"},
	}

	table := widget.NewTable(
		func() (rows int, cols int) {
			return len(shortcuts), 2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			switch id.Col {
			case 0:
				obj.(*widget.Label).SetText(shortcuts[id.Row].shortcut)
			case 1:
				obj.(*widget.Label).SetText(shortcuts[id.Row].info)
			}
		},
	)
	table.SetColumnWidth(0, 200)
	table.SetColumnWidth(1, 200)


	var shortcutsPopUp *widget.PopUp
	shortcutsPopUp = widget.NewModalPopUp(container.NewBorder(
		// Insanity again
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
				shortcutsPopUp.Hide()
			}),
		)),
		nil,
		nil,
		nil,
		container.NewPadded(table),
	), canvas)

	return shortcutsPopUp
}

func aboutPopUp(canvas fyne.Canvas) *widget.PopUp {
	text := `
RunGo is a free cross-platform Go playground, that allows the user to experiment, prototype and get instant feedback. It provides support for running Go versions from 1.16+, and is built on top of Fyne, a cross-platform GUI toolkit made with Go and insipred by Material Design.

MIT License
Copyright (c) 2023 Kevin Suñer
	`

	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrap(fyne.TextWrapBreak)

	fyneURL, _ := url.Parse("https://github.com/fyne-io/fyne/blob/master/LICENSE")
	goQueryURL, _ := url.Parse("https://github.com/PuerkitoBio/goquery/blob/master/LICENSE")
	runGoURL, _ := url.Parse("https://github.com/itsksrof/run-go")

	var aboutPopUp *widget.PopUp
	aboutPopUp = widget.NewModalPopUp(container.NewBorder(
		// Insanity again
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
				aboutPopUp.Hide()
			}),
		)),
		nil,
		nil,
		nil,
		container.NewPadded(
			container.NewVBox(
				label,
				widget.NewHyperlink("Fyne.io BSD 3-Clause License", fyneURL),
				widget.NewHyperlink("goquery BSD 3-Clause License", goQueryURL),
				widget.NewHyperlink("RunGo GitHub Repository", runGoURL),
			),
		),
	), canvas)

	return aboutPopUp
}

func goVersionPopUp(
	canvas fyne.Canvas,  
	goVersionBtn *widget.Button, 
	goVersionStr binding.String,
) *widget.PopUp {
	goVersions, err := getGoVersions()
	if err != nil {
		// Should display a modal with the error and in (debug) mode
		// log to the console
		log.Fatalln(err)
	}

	var goVersionPopUp *widget.PopUp
	goVersionPopUp = widget.NewModalPopUp(container.NewBorder(
		// I know, this looks insane, but it is the only way I found for it to work
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
				goVersionPopUp.Hide()
			}),
		)),
		nil,
		nil,
		nil,
		container.NewPadded(
			widget.NewList(
				func() int {
					return len(goVersions)
				},
				func() fyne.CanvasObject {
					return widget.NewButton("template", nil)
				},
				func(id widget.ListItemID, obj fyne.CanvasObject) {
					button := obj.(*widget.Button)
					button.SetText(goVersions[id])
					button.Alignment = widget.ButtonAlignLeading
					button.OnTapped = func() {
						appDir := fmt.Sprintf("%s/%s", os.Getenv("RUNGO_HOME"), APP_DIR)
						gosDir := fmt.Sprintf("%s/%s", appDir, GOS_DIR)
						goVersion := fmt.Sprintf("%s.%s-%s", 
							button.Text, 
							os.Getenv("RUNGO_OS"),
							os.Getenv("RUNGO_ARCH"),
						)

						_, err := os.ReadDir(fmt.Sprintf("%s/%s", gosDir, goVersion))
						if os.IsNotExist(err) {
							if err := downloadGoTarball(goVersion, appDir); err != nil {
								log.Fatalln("is this")
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
						
						goVersionStr.Set(goVersions[id])
						version, err := goVersionStr.Get()
						if err != nil { log.Fatalln(err) }

						goVersionBtn.SetText(version)
						goVersionPopUp.Hide()
					}
				},
			),
		),
	), canvas)

	return goVersionPopUp
}

type saveSnippet struct {
	*widget.PopUp
}

func saveSnippetPopUp(
	entry *widget.Entry,
	appTabs *container.AppTabs,
	snippet binding.String,
	canvas fyne.Canvas,
) *saveSnippet {
	saveSnippet := &saveSnippet{}

	var saveSnippetPopUp *widget.PopUp
	input := &widget.Entry{PlaceHolder: "Snippet name"}
	saveSnippetPopUp = widget.NewModalPopUp(container.NewBorder(
		// Insanity again
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
				saveSnippetPopUp.Hide()
			}),
		)),
		nil,
		nil,
		nil,
		container.NewPadded(
			container.NewVBox(
				input,
				widget.NewButtonWithIcon("Save", theme.ConfirmIcon(), func() {
					err := newSnippet(input.Text, []byte(entry.Text))
					if err != nil { log.Fatalln(err) }
					
					snippet.Set(input.Text)
					appTabs.Selected().Text = input.Text
					appTabs.Refresh()
					saveSnippetPopUp.Hide()
				}),
			),
		),
	), canvas)

	saveSnippet.PopUp = saveSnippetPopUp
	return saveSnippet
}

func (s *saveSnippet) TypedShortcut(shortcut fyne.Shortcut) {
	customShortcut, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		s.TypedShortcut(shortcut)
		return
	}

	switch customShortcut.ShortcutName() {
	case ALT_S:
		s.PopUp.Resize(fyne.NewSize(440, 200))
		s.PopUp.Show()
	}
}

type loadSnippet struct {
	snippetList binding.StringList
	*widget.PopUp
}

func loadSnippetPopUp(
	entry *widget.Entry,
	appTabs *container.AppTabs,
	snippet binding.String,
	snippetList binding.StringList,
	canvas fyne.Canvas,
) *loadSnippet {
	loadSnippet := &loadSnippet{snippetList: snippetList}

	var loadSnippetPopUp *widget.PopUp
	loadSnippetPopUp = widget.NewModalPopUp(container.NewBorder(
		// Insanity again
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
				loadSnippetPopUp.Hide()
			}),
		)),
		nil,
		nil,
		nil,
		container.NewPadded(
			widget.NewList(
				func() int {
					return snippetList.Length()
				},
				func() fyne.CanvasObject {
					return widget.NewButton("template", nil)
				},
				func(id widget.ListItemID, obj fyne.CanvasObject) {
					snippetName, err := snippetList.GetValue(id)
					if err != nil { log.Fatalln(err) }

					button := obj.(*widget.Button)
					button.SetText(snippetName)
					button.Alignment = widget.ButtonAlignLeading
					button.OnTapped = func() {
						if len(entry.Text) != 0 {
							loadSnippetPopUp.Hide()
							return
						}

						data, err := openSnippet(button.Text)
						if err != nil { log.Fatalln(err) }
						entry.SetText(data)

						snippet.Set(button.Text)
						appTabs.Selected().Text = button.Text
						appTabs.Refresh()
						loadSnippetPopUp.Hide()
					}
				},
			),

		),
	), canvas)

	loadSnippet.PopUp = loadSnippetPopUp
	return loadSnippet
}

func (l *loadSnippet) TypedShortcut(shortcut fyne.Shortcut) {
	customShortcut, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		l.TypedShortcut(shortcut)
		return
	}

	switch customShortcut.ShortcutName() {
	case ALT_O:
		snippets, err := listSnippets()
		if err != nil { log.Fatalln(err) }
		l.snippetList.Set(snippets)

		l.PopUp.Resize(fyne.NewSize(440, 540))
		l.PopUp.Show()
	}
}

func getGoVersions() ([]string, error) {
	// This should be cached, I'm doing some webscraping, yes,
	// but I'm not a complete dick
	resp, err := http.Get(GO_DOWNLOADS_URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	goVersionsRaw := make([]string, 0)
	doc.Find(".toggleButton").Each(func(i int, s *goquery.Selection) {
		goVersion := s.Find("span").Text()
		
		// Match versions from go1.16 ahead and replace the leading "go" prefix
		// for a "v" prefix, allowing us to sort it using the semver package
		r := regexp.MustCompile(`^go(\d+)\.(1[6-9]|[2-9]\d+)(?:\.(\d+))?$`)
		if r.MatchString(goVersion) {
			goVersionsRaw = append(
				goVersionsRaw, 
				strings.Replace(goVersion, "go", "v", 1),
			)
		}
	})

	goVersionsRaw = slices.Compact(goVersionsRaw)
	semver.Sort(goVersionsRaw)
	slices.Reverse(goVersionsRaw)

	goVersions := make([]string, 0)
	for _, goVersionRaw := range goVersionsRaw {
		goVersions = append(goVersions, strings.Replace(goVersionRaw, "v", "go", 1))
	}

	return goVersions, nil
}
