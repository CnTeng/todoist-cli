package sync

import (
	"context"
	"fmt"

	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/urfave/cli/v3"
)

func NewCmd(f *util.Factory) *cli.Command {
	params := &daemon.SyncArgs{}
	return &cli.Command{
		Name:                   "sync",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "force",
				Aliases:     []string{"f"},
				Usage:       "force sync",
				OnlyOnce:    true,
				Value:       false,
				Destination: &params.IsForce,
			},
			&cli.BoolFlag{
				Name:        "all",
				Aliases:     []string{"a"},
				Usage:       "sync all items",
				OnlyOnce:    true,
				Value:       false,
				Destination: &params.All,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if _, err := f.Call(ctx, daemon.Sync, params); err != nil {
				return err
			}

			fmt.Println("Sync success")
			return nil
		},
	}
}
