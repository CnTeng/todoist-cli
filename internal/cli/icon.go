package cli

import "github.com/fatih/color"

type icons struct {
	done   string
	undone string

	none string

	inbox    string
	favorite string

	indent     string
	lastIndent string
}

var nerdIcons = icons{
	done:   " ",
	undone: " ",

	none: "  ",

	inbox:    " ",
	favorite: " ",

	indent:     "│ ",
	lastIndent: "└ ",
}

var textIcons = icons{
	done:   "[x]",
	undone: "[ ]",

	none: "  ",

	inbox:    "IN",
	favorite: "* ",
}

type iconType int

const (
	Nerd iconType = iota
	Text
)

func newIcons(t iconType) *icons {
	var icons icons

	if t == Nerd {
		icons = nerdIcons
	} else {
		icons = textIcons
	}

	return &icons
}

func (i *icons) withColor() *icons {
	i.done = color.GreenString(i.done)

	return i
}
