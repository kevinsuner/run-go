package main

import (
	"log"
	"run-go/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
)

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
