// SPDX-License-Identifier: MIT
package main

import (
	"fmt"
	"log"
	"os"
	"run-go/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

const APP_DIR string = ".run-go/snippets"

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dir := fmt.Sprintf("%s/%s", home, APP_DIR)
	_, err = os.ReadDir(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal(err)
		}
	}
}

var (
	ctrlT = &desktop.CustomShortcut{KeyName: fyne.KeyT, Modifier: fyne.KeyModifierControl}
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("RunGo")

	tabbar := widgets.NewTabBar(myWindow.Canvas())

	myWindow.Canvas().AddShortcut(ctrlT, tabbar.TypedShortcut)
	myWindow.Canvas().SetContent(tabbar.AppTabs)
	myWindow.Resize(fyne.NewSize(1024, 640))
	myWindow.ShowAndRun()
}
