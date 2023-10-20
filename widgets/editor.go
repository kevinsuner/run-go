// SPDX-License-Identifier: MIT
package widgets

import (
	"log"
	"run-go/events"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const CUSTOM_SHORTCUT_CTRL_RETURN string = "CustomDesktop:Control+Return"

type Editor struct {
	widget.Entry
	binding.String
}

func NewEditor(str binding.String) *Editor {
	editor := &Editor{String: str}
	editor.MultiLine = true
	editor.ExtendBaseWidget(editor)

	return editor
}

func (e *Editor) TypedShortcut(shortcut fyne.Shortcut) {
	s, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	switch s.ShortcutName() {
	case CUSTOM_SHORTCUT_CTRL_RETURN:
		// Run code
		out, err := events.CreateTempAndRun([]byte(e.Text))
		if err != nil {
			log.Fatal(err)
		}

		e.String.Set(out)
	}
}
