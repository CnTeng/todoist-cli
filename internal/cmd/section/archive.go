package section

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
		Use:               "archive [flags] <section-id>...",
		Short:             "Archive sections",
		Long:              "Archive sections in Todoist.",
		Example:           "  todoist section archive 6X7FxXvX84jHphx2",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: f.NewSectionCompletionFunc(-1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params := []*sync.SectionArchiveArgs{}
			for _, arg := range args {
				params = append(params, &sync.SectionArchiveArgs{ID: arg})
			}
			if _, err := f.Call(cmd.Context(), daemon.SectionArchive, params); err != nil {
				return err
			}

			fmt.Printf("Sections archived: %s\n", strings.Join(args, ", "))
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}

func NewUnarchiveCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "unarchive [flags] <section-id>...",
		Short:             "Unarchive sections",
		Long:              "Unarchive sections in Todoist.",
		Example:           "  todoist section unarchive 6X7FxXvX84jHphx2",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: f.NewSectionCompletionFunc(-1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params := []*sync.SectionUnarchiveArgs{}
			for _, arg := range args {
				params = append(params, &sync.SectionUnarchiveArgs{ID: arg})
			}
			if _, err := f.Call(cmd.Context(), daemon.SectionUnarchive, params); err != nil {
				return err
			}

			fmt.Printf("Sections unarchived: %s\n", strings.Join(args, ", "))
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
