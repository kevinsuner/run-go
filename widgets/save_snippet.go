// SPDX-License-Identifier: MIT
package widgets

import (
	"fmt"
	"os"
	"run-go/events"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const CUSTOM_SHORTCUT_CTRL_S string = "CustomDesktop:Control+S"

type SaveSnippet struct {
	*widget.PopUp
}

func NewSaveSnippet(editor *Editor, canvas fyne.Canvas) *SaveSnippet {
	saveSnippet := &SaveSnippet{}

	entry := widget.NewEntry()
	saveSnippet.PopUp = widget.NewModalPopUp(&widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: entry},
		},
		OnSubmit: func() {
			err := events.CreateGoProject(entry.Text, []byte(editor.Entry.Text))
			if err != nil {
				// TODO: Proper error handling
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if err := editor.SnippetName.Set(entry.Text); err != nil {
				// TODO: Proper error handling
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			saveSnippet.PopUp.Hide()
		},
	}, canvas)

	return saveSnippet
}

func (s *SaveSnippet) TypedShortcut(shortcut fyne.Shortcut) {
	cs, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		s.TypedShortcut(shortcut)
		return
	}

	switch cs.ShortcutName() {
	case CUSTOM_SHORTCUT_CTRL_S:
		s.PopUp.Resize(fyne.NewSize(640, 360))
		s.PopUp.Show()
	}
}
