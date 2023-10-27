// SPDX-License-Identifier: MIT
package widgets

import (
	"fyne.io/fyne/v2/widget"
)

func NewForm(text string, entry *widget.Entry, onSubmit func()) *widget.Form {
	return &widget.Form{
		Items: []*widget.FormItem{
			{Text: text, Widget: entry},
		},
		OnSubmit: onSubmit,
	}
}
