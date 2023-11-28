package widgets

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const CTRL_U string = "CustomDesktop:Control+U"

type ShortcutsPopUp struct {
	*widget.PopUp
}

var data = []string{
	"Open saved snippets modal CTRL+O",
	"Open save snippet modal CTRL+S",
	"Open a new tab CTRL+T",
	"Open Go version modal CTRL+G",
	"Open about modal CTRL+A",
	"Open shortcuts modal CTRL+U",
	"Run code from editor CTRL+ENTER",
}

func NewShortcutsPopUp(canvas fyne.Canvas) *ShortcutsPopUp {
	shortcutsPopUp := &ShortcutsPopUp{}
	shortcutsPopUp.PopUp = widget.NewModalPopUp(container.NewGridWithRows(2,
		widget.NewList(
			func() (int) {
				return len(data)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("template")
			},
			func(id widget.ListItemID, obj fyne.CanvasObject) {
				obj.(*widget.Label).SetText(data[id])
			},
		),
		widget.NewButton("Close", func() {
			fmt.Fprintln(os.Stdout, "closeModalButton clicked")
			shortcutsPopUp.PopUp.Hide()
		}),
	), canvas)

	return shortcutsPopUp
}

func (s *ShortcutsPopUp) TypedShortcut(shortcut fyne.Shortcut) {
	customShortcut, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		s.TypedShortcut(shortcut)
		return
	}

	switch customShortcut.ShortcutName() {
	case CTRL_U:
		s.PopUp.Resize(fyne.NewSize(1024, 640))
		s.PopUp.Show()
	}
}
