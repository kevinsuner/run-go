/*
    SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
    SPDX-License-Identifier: MIT
*/
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

const CTRL_RETURN string = "CustomDesktop:Control+Return"

type Input struct {
	output      binding.String
	projectName binding.String

	widget.Entry
}

func NewInput(out, projectName binding.String) *Input {
	input := &Input{output: out, projectName: projectName}
	input.MultiLine = true
	input.ExtendBaseWidget(input)

	return input
}

func (i *Input) TypedShortcut(shortcut fyne.Shortcut) {
	customShortcut, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		i.Entry.TypedShortcut(shortcut)
		return
	}

	switch customShortcut.ShortcutName() {
	case CTRL_RETURN:
		projectName, err := i.projectName.Get()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if len(projectName) == 0 {
			out, err := events.CreateTempAndRun([]byte(i.Text))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if err := i.output.Set(out); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			return
		}

		out, err := events.RunGoProject(projectName, []byte(i.Text))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := i.output.Set(out); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
