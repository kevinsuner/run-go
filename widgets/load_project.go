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
	input *Input
	projectName binding.String
	tabs *container.AppTabs
	canvas fyne.Canvas
}

func NewLoadProjectPopUp(
	input *Input,
	projectName binding.String,
	tabs *container.AppTabs,
	canvas fyne.Canvas,
) *LoadProjectPopUp {
	loadProjectPopUp := &LoadProjectPopUp{
		input: input,
		projectName: projectName,
		tabs: tabs,
		canvas: canvas,
	}
	loadProjectPopUp.PopUp = widget.NewModalPopUp(nil, canvas)

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
		projects, err := events.ListGoProjects()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		l.PopUp = widget.NewModalPopUp(widget.NewList(
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
					if l.input.Entry.Text == "" {
						data, err := events.LoadGoProject(button.Text)
						if err != nil {
							fmt.Fprintln(os.Stderr, err)
							os.Exit(1)
						}

						if err := l.projectName.Set(button.Text); err != nil {
							fmt.Fprintln(os.Stderr, err)
							os.Exit(1)
						}

						l.tabs.Selected().Text = button.Text
						l.tabs.Refresh()
						l.input.Entry.SetText(data)
					}

					l.PopUp.Hide()
				}
			},
		), l.canvas)

		l.PopUp.Resize(fyne.NewSize(1024, 640))
		l.PopUp.Show()
	}
}
