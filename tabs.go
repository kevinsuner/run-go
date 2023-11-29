package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
)

type tabs struct {
	canvas fyne.Canvas
	*container.AppTabs
}

func appTabs(canvas fyne.Canvas) *tabs {
	tabs := &tabs{canvas: canvas}

	output := binding.NewString()
	snippet := binding.NewString()

	editor := playgroundEditor(output, snippet)
	console := playgroundConsole(output)

	tabs.AppTabs = container.NewAppTabs(
		container.NewTabItem("New snippet", container.NewGridWithColumns(2,
			editor,
			console,
		)),
	)

	saveSnippetPopUp := saveSnippetPopUp(
		&editor.Entry, 
		tabs.AppTabs, 
		snippet, 
		canvas,
	)

	loadSnippetPopUp := loadSnippetPopUp(
		&editor.Entry,
		tabs.AppTabs,
		snippet,
		snippetList,
		canvas,
	)
	
	canvas.AddShortcut(altReturn, editor.Entry.TypedShortcut)
	canvas.AddShortcut(altS, saveSnippetPopUp.TypedShortcut)
	canvas.AddShortcut(altO, loadSnippetPopUp.TypedShortcut)

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
		t.Append(newTab(t.AppTabs, t.canvas))
	}
}

func newTab(appTabs *container.AppTabs, canvas fyne.Canvas) *container.TabItem {
	output := binding.NewString()
	snippet := binding.NewString()

	editor := playgroundEditor(output, snippet)
	console := playgroundConsole(output)

	saveSnippetPopUp := saveSnippetPopUp(
		&editor.Entry, 
		appTabs, 
		snippet, 
		canvas,
	)

	loadSnippetPopUp := loadSnippetPopUp(
		&editor.Entry,
		appTabs,
		snippet,
		snippetList,
		canvas,
	)
	
	canvas.AddShortcut(altReturn, editor.Entry.TypedShortcut)
	canvas.AddShortcut(altS, saveSnippetPopUp.TypedShortcut)
	canvas.AddShortcut(altO, loadSnippetPopUp.TypedShortcut)

	return container.NewTabItem("New snippet", container.NewGridWithColumns(2,
		editor,
		console,
	))
}
