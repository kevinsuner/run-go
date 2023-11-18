/*
    SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
    SPDX-License-Identifier: MIT
*/
package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Output struct {
	widget.Label
}

func NewOutput(out binding.String) *Output {
	output := &Output{}
	output.Label.Wrapping = fyne.TextWrapBreak
	output.Bind(out)
	output.ExtendBaseWidget(output)

	return output
}
