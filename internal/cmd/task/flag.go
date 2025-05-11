package task

import (
	"time"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/cmd/value"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func addContentFlag(_ *util.Factory, cmd *cobra.Command, destination **string) {
	v := value.NewStringPtr(func(v string) error {
		*destination = &v
		return nil
	})
	flag := &pflag.Flag{
		Name:      "content",
		Shorthand: "c",
		Usage:     "Task content",
		Value:     v,
		DefValue:  v.String(),
	}

	cmd.Flags().AddFlag(flag)
}

func addDeadlineFlag(f *util.Factory, cmd *cobra.Command, destination **sync.Deadline) {
	v := value.NewTimeValue(
		time.DateOnly,
		func(v time.Time) error {
			*destination = &sync.Deadline{Date: &v}
			return nil
		})
	flag := &pflag.Flag{
		Name:     "deadline",
		Usage:    "Deadline date",
		Value:    v,
		DefValue: v.String(),
	}

	cmd.Flags().AddFlag(flag)
	_ = cmd.RegisterFlagCompletionFunc(flag.Name, f.NewDeadlineCompletionFunc(-1))
}

func addDescriptionFlag(_ *util.Factory, cmd *cobra.Command, destination **string) {
	v := value.NewStringPtr(func(v string) error {
		*destination = &v
		return nil
	})
	flag := &pflag.Flag{
		Name:      "description",
		Shorthand: "D",
		Usage:     "Task description",
		Value:     v,
		DefValue:  v.String(),
	}

	cmd.Flags().AddFlag(flag)
}

func addDueFlag(_ *util.Factory, cmd *cobra.Command, destination **sync.Due) {
	v := value.NewStringPtr(func(v string) error {
		*destination = &sync.Due{String: &v}
		return nil
	})
	flag := &pflag.Flag{
		Name:      "due",
		Shorthand: "d",
		Usage:     "Due date",
		Value:     v,
		DefValue:  v.String(),
	}

	cmd.Flags().AddFlag(flag)
}

func addDurationFlag(_ *util.Factory, cmd *cobra.Command, destination **sync.Duration) {
	v := value.NewStringPtr(func(v string) error {
		duration, err := sync.ParseDuration(v)
		if err != nil {
			return err
		}
		*destination = duration
		return nil
	})
	flag := &pflag.Flag{
		Name:     "duration",
		Usage:    "Duration",
		Value:    v,
		DefValue: v.String(),
	}

	cmd.Flags().AddFlag(flag)
	_ = cmd.RegisterFlagCompletionFunc(flag.Name, cobra.NoFileCompletions)
}

func addLabelsFlag(f *util.Factory, cmd *cobra.Command, destination *[]string) {
	cmd.Flags().StringArrayVarP(destination, "labels", "l", nil, "Labels")
}

func addParentFlag(f *util.Factory, cmd *cobra.Command, destination **string) {
	v := value.NewStringPtr(func(v string) error {
		*destination = &v
		return nil
	})
	flag := &pflag.Flag{
		Name:     "parent",
		Usage:    "Parent task ID",
		Value:    v,
		DefValue: v.String(),
	}

	cmd.Flags().AddFlag(flag)
	_ = cmd.RegisterFlagCompletionFunc(flag.Name, f.NewTaskCompletionFunc(-1))
}

func addPriorityFlag(f *util.Factory, cmd *cobra.Command, destination **int) {
	v := value.NewIntValue(func(v int) error {
		*destination = &v
		return nil
	})
	flag := &pflag.Flag{
		Name:      "priority",
		Shorthand: "p",
		Usage:     "Task priority",
		Value:     v,
		DefValue:  "1",
	}

	cmd.Flags().AddFlag(flag)
}

func addProjectFlag(f *util.Factory, cmd *cobra.Command, destination **string) {
	v := value.NewStringPtr(func(v string) error {
		*destination = &v
		return nil
	})
	flag := &pflag.Flag{
		Name:      "project",
		Shorthand: "P",
		Usage:     "Project ID",
		Value:     v,
		DefValue:  v.String(),
	}

	cmd.Flags().AddFlag(flag)
	_ = cmd.RegisterFlagCompletionFunc(flag.Name, f.NewProjectCompletionFunc(-1))
}

func addSectionFlag(f *util.Factory, cmd *cobra.Command, destination **string) {
	v := value.NewStringPtr(func(v string) error {
		*destination = &v
		return nil
	})
	flag := &pflag.Flag{
		Name:      "section",
		Shorthand: "s",
		Usage:     "Section ID",
		Value:     v,
		DefValue:  v.String(),
	}

	cmd.Flags().AddFlag(flag)
}
