package daemon

import (
	"fmt"
	"os"

	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/db"
	"github.com/spf13/cobra"
)

const (
	apiTokenEnv     = "API_TOKEN"
	apiTokenFileEnv = "API_TOKEN_FILE"
)

func loadApiToken() (string, error) {
	apiToken := os.Getenv(apiTokenEnv)
	apiTokenFile := os.Getenv(apiTokenFileEnv)
	if apiToken == "" && apiTokenFile == "" {
		return "", fmt.Errorf("%s or %s is required", apiTokenEnv, apiTokenFileEnv)
	}
	if apiToken == "" && apiTokenFile != "" {
		token, err := os.ReadFile(apiTokenFile)
		if err != nil {
			return "", err
		}
		apiToken = string(token)
	}

	return apiToken, nil
}

func NewCmd(f *util.Factory) *cobra.Command {
	return &cobra.Command{
		Use:          "daemon",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := f.LoadConfig(); err != nil {
				return err
			}

			token, err := loadApiToken()
			if err != nil {
				return err
			}

			f.DeamonConfig.ApiToken = token
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := db.NewDB(f.DataFilePath)
			if err != nil {
				return err
			}

			if err := db.Migrate(); err != nil {
				return err
			}

			daemon := daemon.NewDaemon(db, f.DeamonConfig)
			return daemon.Serve(cmd.Context())
		},
	}
}
