package project

import (
	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/cmd/value"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func addColorFlag(f *util.Factory, cmd *cobra.Command, destination **sync.Color) {
	v := value.NewStringPtr(func(v string) error {
		color, err := sync.ParseColor(v)
		if err != nil {
			return err
		}
		*destination = &color
		return nil
	})
	flag := &pflag.Flag{
		Name:      "color",
		Shorthand: "c",
		Usage:     "Project color",
		Value:     v,
		DefValue:  v.String(),
	}

	cmd.Flags().AddFlag(flag)
	_ = cmd.RegisterFlagCompletionFunc(flag.Name, f.NewColorCompletionFunc(-1))
}

func addFavoriteFlag(f *util.Factory, cmd *cobra.Command, destination **bool) {
	v := value.NewBoolPtr(func(v bool) error {
		*destination = &v
		return nil
	})
	flag := &pflag.Flag{
		Name:        "favorite",
		Shorthand:   "f",
		Usage:       "Add project to favorites",
		Value:       v,
		NoOptDefVal: "true",
		DefValue:    v.String(),
	}

	cmd.Flags().AddFlag(flag)
	_ = cmd.RegisterFlagCompletionFunc(flag.Name, f.NewFavoriteCompletionFunc(-1))
}

func addNameFlag(f *util.Factory, cmd *cobra.Command, destination **string) {
	v := value.NewStringPtr(func(v string) error {
		*destination = &v
		return nil
	})
	flag := &pflag.Flag{
		Name:      "name",
		Shorthand: "n",
		Usage:     "Project name",
		Value:     v,
		DefValue:  v.String(),
	}

	cmd.Flags().AddFlag(flag)
	_ = cmd.RegisterFlagCompletionFunc(flag.Name, f.NewProjectCompletionFunc(-1))
}
