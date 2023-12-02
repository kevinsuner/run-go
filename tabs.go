/*
	SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
	SPDX-License-Identifier: MIT
*/
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
)

type customAppTabs struct {
	window fyne.Window
	*container.AppTabs
}

func newAppTabs(window fyne.Window) *customAppTabs {
	var (
		output = binding.NewString()
		snippet = binding.NewString()
		snippetList = binding.NewStringList()
	)

	appTabs := &customAppTabs{window: window}

	editor := playgroundEditor(output, snippet)
	console := playgroundConsole(output)

	appTabs.AppTabs = container.NewAppTabs(
		container.NewTabItem("New snippet", container.NewGridWithColumns(2,
			editor,
			console,
		)),
	)

	saveModal := newSaveModal(&editor.Entry, appTabs.AppTabs, snippet, window)
	openModal := newOpenModal(&editor.Entry, appTabs.AppTabs, snippet, snippetList, window)
	
	window.Canvas().AddShortcut(altReturn, editor.Entry.TypedShortcut)
	window.Canvas().AddShortcut(altS, saveModal.TypedShortcut)
	window.Canvas().AddShortcut(altO, openModal.TypedShortcut)

	return appTabs
}

func (c *customAppTabs) TypedShortcut(shortcut fyne.Shortcut) {
	customShortcut, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		c.TypedShortcut(shortcut)
		return
	}

	switch customShortcut.ShortcutName() {
	case ALT_T:
		c.Append(newTab(c.AppTabs, c.window))
	}
}

func newTab(appTabs *container.AppTabs, window fyne.Window) *container.TabItem {
	var (
		output = binding.NewString()
		snippet = binding.NewString()
		snippetList = binding.NewStringList()
	)

	editor := playgroundEditor(output, snippet)
	console := playgroundConsole(output)

	saveModal := newSaveModal(&editor.Entry, appTabs, snippet, window)
	openModal := newOpenModal(&editor.Entry, appTabs, snippet, snippetList, window)

	window.Canvas().AddShortcut(altReturn, editor.Entry.TypedShortcut)
	window.Canvas().AddShortcut(altS, saveModal.TypedShortcut)
	window.Canvas().AddShortcut(altO, openModal.TypedShortcut)

	return container.NewTabItem("New snippet", container.NewGridWithColumns(2,
		editor,
		console,
	))
}
