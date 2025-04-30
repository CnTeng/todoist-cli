package sync

import (
	"fmt"
	"time"

	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/cmd/value"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewCmd(f *util.Factory) *cobra.Command {
	params := &daemon.SyncArgs{}
	cmd := &cobra.Command{
		Use: "sync",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := f.Call(cmd.Context(), daemon.Sync, params); err != nil {
				return err
			}

			fmt.Println("Sync success")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&params.Force, "force", "f", false, "force sync")
	cmd.Flags().BoolVarP(&params.All, "all", "a", false, "sync all items")
	cmd.Flags().AddFlag(newSinceFlag(&params.Since))

	return cmd
}

func newSinceFlag(destination *time.Time) *pflag.Flag {
	v := value.NewTimeValue(time.DateOnly, func(v time.Time) error {
		*destination = v
		return nil
	})
	return &pflag.Flag{
		Name:      "since",
		Shorthand: "s",
		Usage:     "completed task since",
		Value:     v,
		DefValue:  time.Now().AddDate(0, -1, 0).String(),
	}
}
