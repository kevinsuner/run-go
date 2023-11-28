package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/PuerkitoBio/goquery"
)

func goVersionsPopUp(canvas fyne.Canvas) *widget.PopUp {
	versions, err := getGoVersions()
	if err != nil {
		// Should display a modal with the error and in (debug) mode
		// log to the console
		log.Fatalln(err)
	}

	var goVersionsPopUp *widget.PopUp
	goVersionsPopUp = widget.NewModalPopUp(container.NewBorder(
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
				goVersionsPopUp.Hide()
			}),
		)),
		nil,
		nil,
		nil,
		container.NewPadded(
			widget.NewList(
				func() int {
					return len(versions)
				},
				func() fyne.CanvasObject {
					return widget.NewButton("template", nil)
				},
				func(id widget.ListItemID, obj fyne.CanvasObject) {
					button := obj.(*widget.Button)
					button.SetText(versions[id])
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

						goVersionsPopUp.Hide()
					}
				},
			),
		),
	), canvas)

	return goVersionsPopUp
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

	versions := make([]string, 0)
	doc.Find(".toggleButton").Each(func(i int, s *goquery.Selection) {
		version := s.Find("span").Text()
		
		// Match versions from go1.16 ahead
		r := regexp.MustCompile(`^go(\d+)\.(1[6-9]|[2-9]\d+)(?:\.(\d+))?$`)
		if r.MatchString(version) {
			versions = append(versions, version)
		}
	})

	versions = slices.Compact(versions)
	slices.Sort(versions)
	slices.Reverse(versions)

	return versions, nil
}
