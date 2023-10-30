// SPDX-License-Identifier: MIT
package widgets

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
)

const CTRL_T string = "CustomDesktop:Control+T"

var (
	ctrlReturn = &desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierControl}
	ctrlS      = &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
	ctrlO      = &desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierControl}
)

type TabBar struct {
	*container.AppTabs
	fyne.Canvas
}

func NewTabBar(canvas fyne.Canvas) *TabBar {
	out, projectName := binding.NewString(), binding.NewString()
	if err := out.Set("Type some code and hit Ctrl+Return to start!"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	input := NewInput(out, projectName)
	output := NewOutput(out)

	tabs := container.NewAppTabs(container.NewTabItem(
		"New Snippet",
		container.New(
			layout.NewGridLayout(2), input, output,
		),
	))

	saveProjectPopUp := NewSaveProjectPopUp(input, projectName, tabs, canvas)
	loadProjectPopUp := NewLoadProjectPopUp(input, projectName, tabs, canvas)

	canvas.AddShortcut(ctrlReturn, input.Entry.TypedShortcut)
	canvas.AddShortcut(ctrlS, saveProjectPopUp.TypedShortcut)
	canvas.AddShortcut(ctrlO, loadProjectPopUp.TypedShortcut)

	return &TabBar{
		AppTabs: tabs,
		Canvas:  canvas,
	}
}

func (t *TabBar) TypedShortcut(shortcut fyne.Shortcut) {
	customShortcut, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		t.TypedShortcut(shortcut)
		return
	}

	switch customShortcut.ShortcutName() {
	case CTRL_T:
		t.AppTabs.Append(newTabItem(t.AppTabs, t.Canvas))
	}
}

func newTabItem(tabs *container.AppTabs, canvas fyne.Canvas) *container.TabItem {
	out, projectName := binding.NewString(), binding.NewString()
	if err := out.Set("Type some code and hit Ctrl+Return to start!"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	input := NewInput(out, projectName)
	output := NewOutput(out)
	saveProjectPopUp := NewSaveProjectPopUp(input, projectName, tabs, canvas)
	loadProjectPopUp := NewLoadProjectPopUp(input, projectName, tabs, canvas)

	canvas.AddShortcut(ctrlReturn, input.Entry.TypedShortcut)
	canvas.AddShortcut(ctrlS, saveProjectPopUp.TypedShortcut)
	canvas.AddShortcut(ctrlO, loadProjectPopUp.TypedShortcut)

	return container.NewTabItem("New Snippet", container.New(
		layout.NewGridLayout(2), input, output,
	))
}
