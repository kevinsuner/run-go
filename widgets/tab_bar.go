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

const CUSTOM_SHORTCUT_CTRL_T string = "CustomDesktop:Control+T"

var (
	ctrlReturn = &desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierControl}
	ctrlS      = &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
)

type TabBar struct {
	*container.AppTabs
	fyne.Canvas
}

func NewTabBar(canvas fyne.Canvas) *TabBar {
	return &TabBar{
		AppTabs: container.NewAppTabs(
			newTabItem(canvas),
		),
		Canvas: canvas,
	}
}

func (t *TabBar) TypedShortcut(shortcut fyne.Shortcut) {
	cs, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		t.TypedShortcut(shortcut)
		return
	}

	switch cs.ShortcutName() {
	case CUSTOM_SHORTCUT_CTRL_T:
		output := binding.NewString()
		if err := output.Set("Type some code and hit Ctrl+Return to start!"); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		t.AppTabs.Append(newTabItem(t.Canvas))
	}
}

func newTabItem(canvas fyne.Canvas) *container.TabItem {
	output := binding.NewString()
	if err := output.Set("Type some code and hit Ctrl*Return to start!"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	editor := NewEditor(output, binding.NewString())
	console := NewConsole(output)
	saveSnippet := NewSaveSnippet(editor, canvas)

	canvas.AddShortcut(ctrlReturn, editor.Entry.TypedShortcut)
	canvas.AddShortcut(ctrlS, saveSnippet.TypedShortcut)

	return container.NewTabItem("New Snippet", container.New(
		layout.NewGridLayout(2), editor, console,
	))
}
