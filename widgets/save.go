// SPDX-License-Identifier: MIT
package widgets

import (
	"log"
	"run-go/events"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// Maybe rename to SaveSnippet???

const CUSTOM_SHORTCUT_CTRL_S string = "CustomDesktop:Control+S"

type SavePopUp struct {
	*widget.PopUp
}

func NewSavePopUp(editor *Editor, canvas fyne.Canvas) *SavePopUp {
	savePopUp := &SavePopUp{}

	entry := widget.NewEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: entry},
		},
		OnSubmit: func() {
			err := events.CreateGoProject(entry.Text, []byte(editor.Entry.Text))
			if err != nil {
				log.Fatal(err)
			}

			savePopUp.PopUp.Hide()
		},
	}

	savePopUp.PopUp = widget.NewModalPopUp(form, canvas)
	return savePopUp
}

func (p *SavePopUp) TypedShortcut(shortcut fyne.Shortcut) {
	s, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		p.TypedShortcut(shortcut)
	}

	switch s.ShortcutName() {
	case CUSTOM_SHORTCUT_CTRL_S:
		p.PopUp.Resize(fyne.NewSize(640, 360))
		p.PopUp.Show()
	}
}
