package view

import "fmt"

type iconType string

const (
	Nerd iconType = "nerd"
	Text iconType = "text"
)

type IconConfig struct {
	IconType iconType `toml:"default"`
	*Icons
}

func (c *IconConfig) UnmarshalTOML(p any) error {
	data, ok := p.(map[string]any)
	if !ok {
		return fmt.Errorf("expected map[string]any, got %T", p)
	}

	if _, ok := data["default"]; ok {
		c.IconType = iconType(data["default"].(string))
		switch c.IconType {
		case Nerd:
			c.Icons = &NerdIcons
		case Text:
			c.Icons = &TextIcons
		}
	}

	if v, ok := data["none"].(string); ok {
		c.None = v
	}
	if v, ok := data["done"].(string); ok {
		c.Done = v
	}
	if v, ok := data["undone"].(string); ok {
		c.Undone = v
	}
	if v, ok := data["inbox"].(string); ok {
		c.Inbox = v
	}
	if v, ok := data["favorite"].(string); ok {
		c.Favorite = v
	}
	if v, ok := data["indent"].(string); ok {
		c.Indent = v
	}
	if v, ok := data["last_indent"].(string); ok {
		c.LastIndent = v
	}

	return nil
}

var DefaultIconConfig = &IconConfig{
	IconType: Nerd,
	Icons:    &NerdIcons,
}

type Icons struct {
	None   string `toml:"none"`
	Done   string `toml:"done"`
	Undone string `toml:"undone"`

	Inbox    string `toml:"inbox"`
	Favorite string `toml:"favorite"`

	Indent     string `toml:"indent"`
	LastIndent string `toml:"last_indent"`
}

var NerdIcons = Icons{
	None:   "  ",
	Done:   " ",
	Undone: " ",

	Inbox:    " ",
	Favorite: " ",

	Indent:     "│ ",
	LastIndent: "└ ",
}

var TextIcons = Icons{
	None:   "  ",
	Done:   "[x]",
	Undone: "[ ]",

	Inbox:    "IN",
	Favorite: "* ",

	Indent:     "  ",
	LastIndent: "  ",
}
