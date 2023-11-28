package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type tabs struct {
	*container.AppTabs
}

func appTabs() *tabs {
	tabs := &tabs{}
	tabs.AppTabs = container.NewAppTabs(
		container.NewTabItem("Tab 1", container.NewGridWithColumns(2,
			widget.NewEntry(),
			widget.NewLabel("Output"),
		)),
	)

	return tabs
}

func (t *tabs) TypedShortcut(shortcut fyne.Shortcut) {
	customShortcut, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		t.TypedShortcut(shortcut)
		return
	}

	switch customShortcut.ShortcutName() {
	case ALT_T:
		t.Append(container.NewTabItem("Tab 2", container.NewGridWithColumns(2,
			widget.NewEntry(),
			widget.NewLabel("Output"),
		)))
	}
}
