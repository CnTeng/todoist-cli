package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

type config struct {
	Token string `json:"token"`
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
			fmt.Printf("Error opening config file: %v\n", err)
			os.Exit(1)
		}

		cfg = &config{}
		if err := json.Unmarshal(file, cfg); err != nil {
			fmt.Printf("Error decoding config file: %v\n", err)
			os.Exit(1)
		}

		return ctx, nil
	},

	Commands: []*cli.Command{taskListCmd, taskAddCmd, daemonCmd, syncCmd, ProjectCmd},
}

func Execute() error {
	return rootCmd.Run(context.Background(), os.Args)
}
