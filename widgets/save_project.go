// SPDX-License-Identifier: MIT
package widgets

import (
	"fmt"
	"os"
	"run-go/events"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const CTRL_S string = "CustomDesktop:Control+S"

type SaveProjectPopUp struct {
	*widget.PopUp
}

func NewSaveProjectPopUp(
	input *Input,
	projectName binding.String,
	tabs *container.AppTabs,
	canvas fyne.Canvas,
) *SaveProjectPopUp {
	saveProjectPopUp := &SaveProjectPopUp{}

	entry := widget.NewEntry()
	saveProjectPopUp.PopUp = widget.NewModalPopUp(NewForm("Name", entry, func() {
		err := events.CreateGoProject(entry.Text, []byte(input.Entry.Text))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := projectName.Set(entry.Text); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		tabs.Selected().Text = entry.Text
		tabs.Refresh()
		saveProjectPopUp.PopUp.Hide()
	}), canvas)

	return saveProjectPopUp
}

func (s *SaveProjectPopUp) TypedShortcut(shortcut fyne.Shortcut) {
	customShortcut, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		s.TypedShortcut(shortcut)
		return
	}

	switch customShortcut.ShortcutName() {
	case CTRL_S:
		s.PopUp.Resize(fyne.NewSize(640, 360))
		s.PopUp.Show()
	}
}
