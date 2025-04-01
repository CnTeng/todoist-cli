package cli

import "github.com/fatih/color"

type progressIcons struct {
	none      string
	done      string
	separator string
	undone    string
}

type icons struct {
	done   string
	undone string

	none   string
	add    string
	change string
	delete string

	inbox    string
	favorite string

	indent     string
	lastIndent string

	progress progressIcons
}

var nerdIcons = icons{
	done:   " ",
	undone: " ",

	none:   "  ",
	add:    "+",
	change: "~",
	delete: "-",

	inbox:    "",
	favorite: "",

	indent:     "│ ",
	lastIndent: "└ ",

	progress: progressIcons{
		none:      "─",
		done:      "━",
		separator: "╺",
		undone:    "━",
	},
}

var textIcons = icons{
	done:   "[x]",
	undone: "[ ]",

	none:   " ",
	add:    "+",
	change: "~",
	delete: "-",

	inbox:    "IN",
	favorite: "*",

	progress: progressIcons{
		none:      "─",
		done:      "━",
		separator: "╺",
		undone:    "━",
	},
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

	i.add = color.GreenString(i.add)
	i.change = color.YellowString(i.change)
	i.delete = color.RedString(i.delete)

	i.progress.done = color.GreenString(i.progress.done)
	// i.progress.undone = color..Sprint(i.progress.undone)
	// i.progress.separator = text.Faint.Sprint(i.progress.separator)
	// i.progress.none = text.Faint.Sprint(i.progress.none)

	return i
}
