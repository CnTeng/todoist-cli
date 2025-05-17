package task

import (
	"time"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/value"
	"github.com/spf13/pflag"
)

func newContentFlag(destination **string) *pflag.Flag {
	v := value.NewStringPtr(func(v string) error {
		*destination = &v
		return nil
	})
	return &pflag.Flag{
		Name:      "name",
		Shorthand: "n",
		Usage:     "Set the task name",
		Value:     v,
		DefValue:  v.String(),
	}
}

func newDeadlineFlag(destination **sync.Deadline) *pflag.Flag {
	v := value.NewTimeValue(
		time.DateOnly,
		func(v time.Time) error {
			*destination = &sync.Deadline{Date: &v}
			return nil
		})
	return &pflag.Flag{
		Name:     "deadline",
		Usage:    "Set the task deadline (format: YYYY-MM-DD)",
		Value:    v,
		DefValue: v.String(),
	}
}

func newDescriptionFlag(destination **string) *pflag.Flag {
	v := value.NewStringPtr(func(v string) error {
		*destination = &v
		return nil
	})
	return &pflag.Flag{
		Name:      "description",
		Shorthand: "D",
		Usage:     "Add a description to the task",
		Value:     v,
		DefValue:  v.String(),
	}
}

func newDueFlag(destination **sync.Due) *pflag.Flag {
	v := value.NewStringPtr(func(v string) error {
		*destination = &sync.Due{String: &v}
		return nil
	})
	return &pflag.Flag{
		Name:      "due",
		Shorthand: "d",
		Usage:     "Set the task due date (e.g. 'today', 'tomorrow', 'every day 22:00')",
		Value:     v,
		DefValue:  v.String(),
	}
}

func newDurationFlag(destination **sync.Duration) *pflag.Flag {
	v := value.NewStringPtr(func(v string) error {
		duration, err := sync.ParseDuration(v)
		if err != nil {
			return err
		}
		*destination = duration
		return nil
	})
	return &pflag.Flag{
		Name:     "duration",
		Usage:    "Set the task duration (format: <amount> <minute|day>, e.g. '30 minute', '1 day')",
		Value:    v,
		DefValue: v.String(),
	}
}

func newParentFlag(destination **string) *pflag.Flag {
	v := value.NewStringPtr(func(v string) error {
		*destination = &v
		return nil
	})
	return &pflag.Flag{
		Name:     "parent",
		Usage:    "Set the parent task by <task-id> to create a subtask",
		Value:    v,
		DefValue: v.String(),
	}
}

func newPriorityFlag(destination **int) *pflag.Flag {
	v := value.NewIntValue(func(v int) error {
		*destination = &v
		return nil
	})
	return &pflag.Flag{
		Name:      "priority",
		Shorthand: "P",
		Usage:     "Set the priority level (1-4, where 4 is the highest priority)",
		Value:     v,
		DefValue:  "1",
	}
}

func newProjectFlag(destination **string) *pflag.Flag {
	v := value.NewStringPtr(func(v string) error {
		*destination = &v
		return nil
	})
	return &pflag.Flag{
		Name:      "project",
		Shorthand: "p",
		Usage:     "Assign the task to a specific project by <project-id>",
		Value:     v,
		DefValue:  v.String(),
	}
}

func newSectionFlag(destination **string) *pflag.Flag {
	v := value.NewStringPtr(func(v string) error {
		*destination = &v
		return nil
	})
	return &pflag.Flag{
		Name:      "section",
		Shorthand: "s",
		Usage:     "Assign the task to a specific section by <section-id>",
		Value:     v,
		DefValue:  v.String(),
	}
}
