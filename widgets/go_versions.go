package widgets

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/PuerkitoBio/goquery"
)

const APP_DIR string = ".run-go"
const CTRL_G string = "CustomDesktop:Control+G"

type GoVersionsPopUp struct {
	*widget.PopUp
}

func NewGoVersionsPopUp(canvas fyne.Canvas) *GoVersionsPopUp {
	goVersionPopUp := &GoVersionsPopUp{}

	// Request all available Go versions
	// This should be cached
	goReleasesURL := "https://go.dev/dl"
	resp, err := http.Get(goReleasesURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, resp.Status)
		os.Exit(1)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	goVersions := make([]string, 0)
	doc.Find(".toggleButton").Each(func(i int, s *goquery.Selection) {
		version := s.Find("span").Text()
		r := regexp.MustCompile(`^go(\d+)\.(1[6-9]|[2-9]\d+)(?:\.(\d+))?$`)
		if r.MatchString(version) {
			goVersions = append(goVersions, version)
		}
	})

	goVersions = slices.Compact(goVersions)
	slices.Sort(goVersions)
	slices.Reverse(goVersions)

	fmt.Println(goVersions)

	// Create the list with all the Go versions and on-click download the binary
	// and update the GOPATH
	goVersionPopUp.PopUp = widget.NewModalPopUp(container.NewGridWithRows(2,
		widget.NewList(
			func() int {
				return len(goVersions)
			},
			func() fyne.CanvasObject {
				return widget.NewButton("template", nil)
			},
			func(id widget.ListItemID, obj fyne.CanvasObject) {
				goVersionButton := obj.(*widget.Button)
				goVersionButton.SetText(goVersions[id])
				goVersionButton.OnTapped = func() {
					home, err := os.UserHomeDir()
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}

					appDir := fmt.Sprintf("%s/%s", home, APP_DIR)
					gosDir := fmt.Sprintf("%s/%s", appDir, "gos")
					
					goDir := fmt.Sprintf("%s/%s.%s-%s", 
						gosDir, 
						goVersionButton.Text, 
						os.Getenv("RUNGO_OS"), 
						os.Getenv("RUNGO_ARCH"),
					)

					_, err = os.ReadDir(goDir)
					if os.IsNotExist(err) {
						// Download Go version tarball in .run-go
						goDownloadURL := fmt.Sprintf("https://go.dev/dl/%s.%s-%s.tar.gz",
							goVersionButton.Text,
							os.Getenv("RUNGO_OS"),
							os.Getenv("RUNGO_ARCH"),
						)

						fmt.Println(goDownloadURL)

						resp, err := http.Get(goDownloadURL)
						if err != nil {
							fmt.Fprintln(os.Stderr, err)
							os.Exit(1)
						}
						defer resp.Body.Close()

						if resp.StatusCode != http.StatusOK {
							fmt.Fprintln(os.Stderr, resp.Status)
							os.Exit(1)
						}

						out, err := os.Create(fmt.Sprintf("%s/%s.%s-%s.tar.gz",
							appDir,
							goVersionButton.Text,
							os.Getenv("RUNGO_OS"),
							os.Getenv("RUNGO_ARCH"),
						))
						if err != nil {
							fmt.Fprintln(os.Stderr, err)
							os.Exit(1)
						}
						defer out.Close()

						_, err = io.Copy(out, resp.Body)
						if err != nil {
							fmt.Fprintln(os.Stderr, err)
							os.Exit(1)
						}

						// Untar the Go version tarball in .run-go/gos
						reader, err := os.Open(fmt.Sprintf("%s/%s.%s-%s.tar.gz",
							appDir,
							goVersionButton.Text,
							os.Getenv("RUNGO_OS"),
							os.Getenv("RUNGO_ARCH"),
						))
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

						// Rename the folder using the Go version, os and architecture of the system
						if err := os.Rename(
							fmt.Sprintf("%s/%s", gosDir, "go"),
							fmt.Sprintf("%s/%s", gosDir,
								fmt.Sprintf("%s.%s-%s", 
									goVersionButton.Text, 
									os.Getenv("RUNGO_OS"),
									os.Getenv("RUNGO_ARCH"),	
								),
							),
						); err != nil {
							fmt.Fprintln(os.Stderr, err)
							os.Exit(1)
						}

						// Remove the downloaded .tar.gz file
						if err := os.Remove(fmt.Sprintf("%s/%s.%s-%s.tar.gz",
							appDir,
							goVersionButton.Text,
							os.Getenv("RUNGO_OS"),
							os.Getenv("RUNGO_ARCH"),
						)); err != nil {
							fmt.Fprintln(os.Stderr, err)
							os.Exit(1)
						}
					}

					// Use the selected Go version
					if err := os.Setenv("GOPATH", fmt.Sprintf("%s/%s.%s-%s/bin/go",
						gosDir,
						goVersionButton.Text,
						os.Getenv("RUNGO_OS"),
						os.Getenv("RUNGO_ARCH"),
					)); err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}

					fmt.Println("go version successfully downloaded and setted")

					goVersionPopUp.PopUp.Hide()
				}
			},
		),
		widget.NewButton("Close", func() {
			fmt.Fprintln(os.Stdout, "closeModalButton clicked")
			goVersionPopUp.PopUp.Hide()
		}),
	), canvas)
	
	return goVersionPopUp
}

func (g *GoVersionsPopUp) TypedShortcut(shortcut fyne.Shortcut) {
	customShortcut, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		g.TypedShortcut(shortcut)
		return
	}

	switch customShortcut.ShortcutName() {
	case CTRL_G:
		g.PopUp.Resize(fyne.NewSize(1024, 640))
		g.PopUp.Show()
	}
}
