package sync

import (
	"context"
	"fmt"
	"time"

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
				Destination: &params.Force,
			},
			&cli.BoolFlag{
				Name:        "all",
				Aliases:     []string{"a"},
				Usage:       "sync all items",
				OnlyOnce:    true,
				Value:       false,
				Destination: &params.All,
			},

			&cli.TimestampFlag{
				Name:        "since",
				Aliases:     []string{"s"},
				Usage:       "completed task since",
				OnlyOnce:    true,
				Value:       time.Now().AddDate(0, -1, 0),
				DefaultText: "1 month ago",
				Destination: &params.Since,
				Config: cli.TimestampConfig{
					Layouts: []string{time.DateOnly},
				},
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
