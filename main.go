package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type KeyableEntry struct {
	binding.String
	widget.Entry
}

func NewKeyableEntry(str binding.String) *KeyableEntry {
	entry := &KeyableEntry{String: str}
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
		timestamp := time.Now().Unix()

		if err := os.WriteFile(
			fmt.Sprintf("%d.go", timestamp),
			[]byte(e.Text),
			0644,
		); err != nil {
			log.Fatal(err)
		}

		out, err := exec.Command("go", "run", fmt.Sprintf("%d.go", timestamp)).CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}

		if err := os.Remove(fmt.Sprintf("%d.go", timestamp)); err != nil {
			log.Fatal(err)
		}

		fmt.Println("[DEBUG]", out)
		fmt.Println("[DEBUG]", string(out))
		e.String.Set(string(out))
	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("RunGo")

	entry := NewKeyableEntry(binding.NewString())
	entry.String.Set("Type some code to start!")
	label := widget.NewLabelWithData(entry.String)

	ctrlReturn := &desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierControl}
	myWindow.Canvas().AddShortcut(ctrlReturn, entry.Entry.TypedShortcut)

	grid := container.New(layout.NewGridLayout(2), entry, label)
	myWindow.Canvas().SetContent(grid)

	myWindow.Resize(fyne.NewSize(1024, 640))
	myWindow.ShowAndRun()
}
