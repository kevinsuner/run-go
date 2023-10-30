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

const CTRL_O string = "CustomDesktop:Control+O"

type LoadProjectPopUp struct {
	*widget.PopUp
}

func NewLoadProjectPopUp(
	input *Input,
	projectName binding.String,
	tabs *container.AppTabs,
	canvas fyne.Canvas,
) *LoadProjectPopUp {
	loadProjectPopUp := &LoadProjectPopUp{}

	projects, err := events.ListGoProjects()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	loadProjectPopUp.PopUp = widget.NewModalPopUp(widget.NewList(
		func() int {
			return len(projects)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("template", nil)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			button := obj.(*widget.Button)
			button.SetText(projects[id])
			button.OnTapped = func() {
				if input.Entry.Text == "" {
					data, err := events.LoadGoProject(button.Text)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}

					if err := projectName.Set(button.Text); err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}

					tabs.Selected().Text = button.Text
					tabs.Refresh()
					input.Entry.SetText(data)
				}

				loadProjectPopUp.PopUp.Hide()
			}
		},
	), canvas)

	return loadProjectPopUp
}

func (l *LoadProjectPopUp) TypedShortcut(shortcut fyne.Shortcut) {
	customShortcut, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		l.TypedShortcut(shortcut)
		return
	}

	switch customShortcut.ShortcutName() {
	case CTRL_O:
		l.PopUp.Resize(fyne.NewSize(1024, 640))
		l.PopUp.Show()
	}
}
