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
	return &cobra.Command{
		Use:   "archive [flags] <section>...",
		Short: "Archive sections",
		Long:  "Archive sections in Todoist.",
		Example: `  todoist archive 6Xm5HVVRcX00MCjv
  todoist section archive 6Xm5HVVRcX00MCjv 6XxxpwJ00459cJWg`,
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
}

func NewUnarchiveCmd(f *util.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "unarchive [flags] <section>...",
		Short: "Unarchive sections",
		Long:  "Unarchive sections in Todoist.",
		Example: `  todoist unarchive 6Xm5HVVRcX00MCjv
  todoist section unarchive 6Xm5HVVRcX00MCjv 6XxxpwJ00459cJWg`,
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
}
