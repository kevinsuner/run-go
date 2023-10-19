// SPDX-License-Identifier: MIT
package main

import (
	"fmt"
	"log"
	"os"
	"run-go/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
)

const APP_DIR string = ".rungo"

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dir := fmt.Sprintf("%s/%s", home, APP_DIR)
	_, err = os.ReadDir(dir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("RunGo")

	ctrlReturn := &desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierControl}

	str := binding.NewString()
	if err := str.Set("Type some code and hit Ctrl+Return to start!"); err != nil {
		log.Fatal(err)
	}

	editor := widgets.NewEditor(str)
	console := widgets.NewConsole(str)

	playground := container.New(layout.NewGridLayout(2), editor, console)

	myWindow.Canvas().AddShortcut(ctrlReturn, editor.Entry.TypedShortcut)
	myWindow.Canvas().SetContent(playground)
	myWindow.Resize(fyne.NewSize(1024, 640))
	myWindow.ShowAndRun()
}
