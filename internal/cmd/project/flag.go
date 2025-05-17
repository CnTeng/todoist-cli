package project

import (
	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/value"
	"github.com/spf13/pflag"
)

func newColorFlag(destination **sync.Color) *pflag.Flag {
	v := value.NewStringPtr(func(v string) error {
		color, err := sync.ParseColor(v)
		if err != nil {
			return err
		}
		*destination = &color
		return nil
	})
	return &pflag.Flag{
		Name:      "color",
		Shorthand: "c",
		Usage:     "Set the project color",
		Value:     v,
		DefValue:  v.String(),
	}
}

func newFavoriteFlag(destination **bool) *pflag.Flag {
	v := value.NewBoolPtr(func(v bool) error {
		*destination = &v
		return nil
	})
	return &pflag.Flag{
		Name:        "favorite",
		Shorthand:   "f",
		Usage:       "Mark project to favorites",
		Value:       v,
		NoOptDefVal: "true",
		DefValue:    v.String(),
	}
}

func newNameFlag(destination **string) *pflag.Flag {
	v := value.NewStringPtr(func(v string) error {
		*destination = &v
		return nil
	})
	return &pflag.Flag{
		Name:      "name",
		Shorthand: "n",
		Usage:     "Set the project name",
		Value:     v,
		DefValue:  v.String(),
	}
}
