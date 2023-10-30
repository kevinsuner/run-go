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

type PopUp struct {
	*widget.PopUp
}

func NewPopUpWithData(
	input *Input,
	projectName binding.String,
	tabs *container.AppTabs,
	canvas fyne.Canvas,
) *PopUp {
	popup := &PopUp{}

	entry := widget.NewEntry()
	popup.PopUp = widget.NewModalPopUp(NewForm("Name", entry, func() {
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
		popup.PopUp.Hide()
	}), canvas)

	return popup
}

func (p *PopUp) TypedShortcut(shortcut fyne.Shortcut) {
	customShortcut, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		p.TypedShortcut(shortcut)
		return
	}

	switch customShortcut.ShortcutName() {
	case CTRL_S:
		p.PopUp.Resize(fyne.NewSize(640, 360))
		p.PopUp.Show()
	}
}
