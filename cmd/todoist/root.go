package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

type config struct {
	Token   string `json:"token"`
	WSToken string `json:"ws_token"`
	Lang    string `json:"lang"`
}

var cfg *config

var configFile string

var rootCmd = &cli.Command{
	Name:  "todoist",
	Usage: "A CLI for Todoist",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Usage:       "config file (default is config.json)",
			Destination: &configFile,
		},
	},
	Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		file, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}

		cfg = &config{}
		if err := json.Unmarshal(file, cfg); err != nil {
			return nil, fmt.Errorf("error decoding config file: %v", err)
		}

		return ctx, nil
	},

	Commands: []*cli.Command{
		taskListCmd(),
		taskAddCmd(),
		taskQuickAddCmd(),
		taskModifyCmd(),
		taskCloseCmd(),
		taskRemoveCmd(),
		taskMoveCmd(),
		projectCmd(),
		syncCmd(),
		daemonCmd(),
	},
}

func Execute() error {
	return rootCmd.Run(context.Background(), os.Args)
}
