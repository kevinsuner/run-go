package main

import (
	"context"
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

type ctxKey struct{}

type KeyableEntry struct {
	context.Context
	widget.Entry
}

func NewKeyableEntry(ctx context.Context) *KeyableEntry {
	entry := &KeyableEntry{Context: ctx}
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
		e.Context = context.WithValue(e.Context, ctxKey{}, string(out))
	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("RunGo")

	entry := NewKeyableEntry(context.Background())
	str := binding.NewString()
	str.Set("Type some code to start!")
	label := widget.NewLabelWithData(str)

	ctrlReturn := &desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierControl}
	myWindow.Canvas().AddShortcut(ctrlReturn, entry.Entry.TypedShortcut)

	grid := container.New(layout.NewGridLayout(2), entry, label)
	myWindow.Canvas().SetContent(grid)

	go func() {
		for range time.Tick(time.Millisecond * 100) {
			if entry.Context.Value(ctxKey{}) == nil {
				continue
			} else {
				str.Set(entry.Context.Value(ctxKey{}).(string))
			}
		}
	}()

	myWindow.Resize(fyne.NewSize(1024, 640))
	myWindow.ShowAndRun()
}
