/*
	SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
	SPDX-License-Identifier: MIT
*/
package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type editor struct {
	output binding.String
	snippet binding.String
	widget.Entry
}

func playgroundEditor(output, snippet binding.String) *editor {
	editor := &editor{output: output, snippet: snippet}
	editor.MultiLine = true
	editor.ExtendBaseWidget(editor)
	return editor
}

func (e *editor) TypedShortcut(shortcut fyne.Shortcut) {
	customShortcut, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
	}

	switch customShortcut.ShortcutName() {
	case ALT_RETURN:
		snippet, err := e.snippet.Get()
		if err != nil { log.Fatalln(err) }

		if len(snippet) != 0 {
			output, err := runFromSnippet(snippet, []byte(e.Text))
			if err != nil { log.Fatalln(err) }

			e.output.Set(output)
		}

		output, err := runFromEditor([]byte(e.Text))
		if err != nil { log.Fatalln(err) }

		e.output.Set(output)
	}
}
