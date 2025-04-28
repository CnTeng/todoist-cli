package project

import (
	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/value"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newNameFlag(destination **string) *pflag.Flag {
	v := value.NewStringPtr(destination)
	return &pflag.Flag{
		Name:      "name",
		Shorthand: "n",
		Usage:     "Project name",
		Value:     v,
		DefValue:  v.String(),
	}
}

func colorCompletion(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	colors := sync.ListColors()

	cmps := make([]cobra.Completion, len(colors))

	for i, c := range colors {
		cmps[i] = cobra.CompletionWithDesc(string(c), c.Hex())
	}
	return cmps, cobra.ShellCompDirectiveNoFileComp
}

func addColorFlag(cmd *cobra.Command, destination **sync.Color) {
	v := value.NewColorPtr(destination)
	flag := &pflag.Flag{
		Name:      "color",
		Shorthand: "c",
		Usage:     "Project color",
		Value:     v,
		DefValue:  v.String(),
	}

	cmd.Flags().AddFlag(flag)
	_ = cmd.RegisterFlagCompletionFunc("color", colorCompletion)
}

func newFavoriteFlag(destination **bool) *pflag.Flag {
	v := value.NewBoolPtr(destination)
	return &pflag.Flag{
		Name:      "favorite",
		Shorthand: "f",
		Usage:     "Add project to favorites",
		Value:     v,
		DefValue:  v.String(),
	}
}
