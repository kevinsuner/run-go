package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type KeyableEntry struct {
	widget.Entry
}

func NewKeyableEntry() *KeyableEntry {
	entry := &KeyableEntry{}
	entry.MultiLine = true
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *KeyableEntry) TypedShortcut(shortcut fyne.Shortcut) {
	cs, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	if cs.Key() == fyne.KeyReturn && cs.Mod() == fyne.KeyModifierControl {
		log.Println(e.Text)
	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("RunGo")

	entry := NewKeyableEntry()
	label := widget.NewLabel("Type some code to start!")

	ctrlReturn := &desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierControl}
	myWindow.Canvas().AddShortcut(ctrlReturn, entry.Entry.TypedShortcut)

	grid := container.New(layout.NewGridLayout(2), entry, label)
	myWindow.Canvas().SetContent(grid)

	myWindow.Resize(fyne.NewSize(1024, 640))
	myWindow.ShowAndRun()
}
