// SPDX-License-Identifier: MIT
package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Console struct {
	widget.Label
}

func NewConsole(output binding.String) *Console {
	console := &Console{}
	console.Label.Wrapping = fyne.TextWrapBreak
	console.Bind(output)
	console.ExtendBaseWidget(console)

	return console
}
