// SPDX-License-Identifier: MIT
package widgets

import (
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Console struct {
	widget.Label
}

func NewConsole(str binding.String) *Console {
	console := &Console{}
	console.Bind(str)
	console.ExtendBaseWidget(console)

	return console
}
