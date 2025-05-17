package project

import (
	"fmt"
	"strings"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewArchiveCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "archive [flags] <project-id>...",
		Short:             "Archive projects",
		Long:              "Archive projects in Todoist.",
		Example:           "  todoist project archive 6X7fphhgwcXVGccJ",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: f.NewProjectCompletionFunc(-1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params := []*sync.ProjectArchiveArgs{}
			for _, arg := range args {
				params = append(params, &sync.ProjectArchiveArgs{ID: arg})
			}
			if _, err := f.Call(cmd.Context(), daemon.ProjectArchive, params); err != nil {
				return err
			}

			fmt.Printf("Projects archived: %s\n", strings.Join(args, ", "))
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}

func NewUnarchiveCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "unarchive [flags] <project-id>...",
		Short:             "Unarchive projects",
		Long:              "Unarchive projects in Todoist.",
		Example:           "  todoist project unarchive 6X7fphhgwcXVGccJ",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: f.NewProjectCompletionFunc(-1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params := []*sync.ProjectUnarchiveArgs{}
			for _, arg := range args {
				params = append(params, &sync.ProjectUnarchiveArgs{ID: arg})
			}
			if _, err := f.Call(cmd.Context(), daemon.ProjectUnarchive, params); err != nil {
				return err
			}

			fmt.Printf("Projects unarchived: %s\n", strings.Join(args, ", "))
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
