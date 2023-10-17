package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type KeyableEntry struct {
	widget.Entry
}

func NewKeyableEntry() *KeyableEntry {
	entry := &KeyableEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *KeyableEntry) TypedShortcut(shortcut fyne.Shortcut) {
	if _, ok := shortcut.(*desktop.CustomShortcut); !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	log.Println("Text:", e.Text)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("RunGo")

	entry := &KeyableEntry{}
	entry.ExtendBaseWidget(entry)

	ctrlReturn := &desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierControl}
	myWindow.Canvas().AddShortcut(ctrlReturn, entry.TypedShortcut)

	myWindow.Canvas().SetContent(entry)

	myWindow.Resize(fyne.NewSize(1024, 640))
	myWindow.ShowAndRun()
}
