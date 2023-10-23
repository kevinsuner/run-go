// SPDX-License-Identifier: MIT
package widgets

import (
	"fmt"
	"os"
	"run-go/events"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const CUSTOM_SHORTCUT_CTRL_RETURN string = "CustomDesktop:Control+Return"

type Editor struct {
	widget.Entry
	Output      binding.String
	SnippetName binding.String
}

func NewEditor(output, snippetName binding.String) *Editor {
	editor := &Editor{Output: output, SnippetName: snippetName}
	editor.MultiLine = true
	editor.ExtendBaseWidget(editor)

	return editor
}

func (e *Editor) TypedShortcut(shortcut fyne.Shortcut) {
	s, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	switch s.ShortcutName() {
	case CUSTOM_SHORTCUT_CTRL_RETURN:
		snippetName, err := e.SnippetName.Get()
		if err != nil {
			// TODO: Proper error handling
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if len(snippetName) == 0 {
			out, err := events.CreateTempAndRun([]byte(e.Text))
			if err != nil {
				// TODO: Proper error handling
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if err := e.Output.Set(out); err != nil {
				// TODO: Proper error handling
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			return
		}

		out, err := events.RunGoProject(snippetName, []byte(e.Text))
		if err != nil {
			// TODO: Proper error handling
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := e.Output.Set(out); err != nil {
			// TODO: Proper error handling
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
