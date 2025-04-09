package project

import (
	"context"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/urfave/cli/v3"
)

func newNameFlag(destination **string) cli.Flag {
	return &cli.StringFlag{
		Name:     "name",
		Aliases:  []string{"n"},
		Usage:    "Project name",
		OnlyOnce: true,
		Action: func(ctx context.Context, cmd *cli.Command, v string) error {
			*destination = &v
			return nil
		},
	}
}

func newColorFlag(destination **sync.Color) cli.Flag {
	return &cli.StringFlag{
		Name:     "color",
		Aliases:  []string{"c"},
		Usage:    "Project color",
		OnlyOnce: true,
		Action: func(ctx context.Context, cmd *cli.Command, v string) error {
			color, err := sync.ParseColor(v)
			if err != nil {
				return err
			}
			*destination = &color
			return nil
		},
	}
}

func newFavoriteFlag(destination **bool) cli.Flag {
	return &cli.BoolFlag{
		Name:     "favorite",
		Aliases:  []string{"f"},
		Usage:    "Add project to favorites",
		OnlyOnce: true,
		Action: func(ctx context.Context, cmd *cli.Command, v bool) error {
			*destination = &v
			return nil
		},
	}
}
