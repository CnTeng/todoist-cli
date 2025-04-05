package view

type Icons struct {
	none   string
	done   string
	undone string

	inbox    string
	favorite string

	indent     string
	lastIndent string
}

var NerdIcons = Icons{
	none:   "  ",
	done:   " ",
	undone: " ",

	inbox:    " ",
	favorite: " ",

	indent:     "│ ",
	lastIndent: "└ ",
}

var TextIcons = Icons{
	none:   "  ",
	done:   "[x]",
	undone: "[ ]",

	inbox:    "IN",
	favorite: "* ",

	indent:     "  ",
	lastIndent: "  ",
}

type iconType int

const (
	Nerd iconType = iota
	Text
)
