package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type console struct {
	widget.Label
}

func playgroundConsole(output binding.String) *console {
	console := &console{}
	console.Label.Wrapping = fyne.TextWrapBreak
	console.Bind(output)
	console.ExtendBaseWidget(console)
	return console
}
