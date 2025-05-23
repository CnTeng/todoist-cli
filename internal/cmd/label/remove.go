package label

import (
	"fmt"
	"strings"

	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/spf13/cobra"
)

func NewRemoveCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "remove [flags] <label-name>...",
		Aliases:           []string{"rm"},
		Short:             "Remove labels",
		Long:              "Remove labels in Todoist, similar to the 'rm' command in shell.",
		Example:           "  todoist label remove Food",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: f.NewLabelCompletionFunc(-1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params := []*model.LabelDeleteArgs{}
			for _, arg := range args {
				params = append(params, &model.LabelDeleteArgs{Name: arg})
			}
			if _, err := f.Call(cmd.Context(), daemon.LabelRemove, params); err != nil {
				return err
			}

			fmt.Printf("Labels deleted: %s\n", strings.Join(args, ", "))
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
