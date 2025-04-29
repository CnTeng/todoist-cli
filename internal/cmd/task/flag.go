package task

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/cmd/value"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func taskCompletion(f *util.Factory) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if err := f.ReadConfig(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.Dial(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		params := &view.TaskViewConfig{}
		result := []*model.Task{}
		if err := f.CallResult(cmd.Context(), daemon.TaskList, params, &result); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		cmps := make([]cobra.Completion, len(result))
		for i, task := range result {
			desc := fmt.Sprintf("%s: %s", task.Project.Name, task.Content)
			cmps[i] = cobra.CompletionWithDesc(task.ID, desc)
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	}
}

func newContentFlag(destination **string) *pflag.Flag {
	v := value.NewStringPtr(destination)
	return &pflag.Flag{
		Name:      "content",
		Shorthand: "c",
		Usage:     "Task content",
		Value:     v,
		DefValue:  v.String(),
	}
}

func newDescriptionFlag(destination **string) *pflag.Flag {
	v := value.NewStringPtr(destination)
	return &pflag.Flag{
		Name:      "description",
		Shorthand: "D",
		Usage:     "Task description",
		Value:     v,
		DefValue:  v.String(),
	}
}

func newProjectFlag(destination **string) *pflag.Flag {
	v := value.NewStringPtr(destination)
	return &pflag.Flag{
		Name:      "project",
		Shorthand: "P",
		Usage:     "Project ID",
		Value:     v,
		DefValue:  v.String(),
	}
}

func newDueFlag(destination **sync.Due) *pflag.Flag {
	v := value.NewDuePtr(destination)
	return &pflag.Flag{
		Name:      "due",
		Shorthand: "d",
		Usage:     "Due date",
		Value:     v,
		DefValue:  v.String(),
	}
}

func newDeadlineFlag(destination **sync.Deadline) *pflag.Flag {
	v := value.NewDeadlinePtr(destination)
	return &pflag.Flag{
		Name:     "deadline",
		Usage:    "Deadline date",
		Value:    v,
		DefValue: v.String(),
	}
}

func newPriorityFlag(destination **int) *pflag.Flag {
	v := value.NewIntPtr(destination)
	return &pflag.Flag{
		Name:      "priority",
		Shorthand: "p",
		Usage:     "Task priority",
		Value:     v,
		DefValue:  "1",
	}
}

func newParentFlag(destination **string) *pflag.Flag {
	v := value.NewStringPtr(destination)
	return &pflag.Flag{
		Name:     "parent",
		Usage:    "Parent task ID",
		Value:    v,
		DefValue: v.String(),
	}
}

func newSectionFlag(destination **string) *pflag.Flag {
	v := value.NewStringPtr(destination)
	return &pflag.Flag{
		Name:      "section",
		Shorthand: "s",
		Usage:     "Section ID",
		Value:     v,
		DefValue:  v.String(),
	}
}

func addLabelsFlag(cmd *cobra.Command, destination *[]string) {
	cmd.Flags().StringArrayVarP(destination, "labels", "l", nil, "Labels")
}

func newDurationFlag(destination **sync.Duration) *pflag.Flag {
	v := value.NewDurationPtr(destination)
	return &pflag.Flag{
		Name:     "duration",
		Usage:    "Duration",
		Value:    v,
		DefValue: v.String(),
	}
}
