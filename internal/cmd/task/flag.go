package task

import (
	"context"
	"time"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/urfave/cli/v3"
)

func newContentFlag(destination **string) cli.Flag {
	return &cli.StringFlag{
		Name:     "content",
		Aliases:  []string{"c"},
		Usage:    "Task content",
		OnlyOnce: true,
		Action: func(ctx context.Context, cmd *cli.Command, v string) error {
			*destination = &v
			return nil
		},
	}
}

func newDescriptionFlag(destination **string) cli.Flag {
	return &cli.StringFlag{
		Name:     "description",
		Aliases:  []string{"D"},
		Usage:    "Task description",
		OnlyOnce: true,
		Action: func(ctx context.Context, cmd *cli.Command, v string) error {
			*destination = &v
			return nil
		},
	}
}

func newProjectFlag(destination **string) cli.Flag {
	return &cli.StringFlag{
		Name:     "project",
		Aliases:  []string{"P"},
		Usage:    "Project ID",
		OnlyOnce: true,
		Action: func(ctx context.Context, cmd *cli.Command, v string) error {
			*destination = &v
			return nil
		},
	}
}

func newDueFlag(destination **sync.Due) cli.Flag {
	return &cli.StringFlag{
		Name:     "due",
		Aliases:  []string{"d"},
		Usage:    "Due date",
		OnlyOnce: true,
		Action: func(ctx context.Context, cmd *cli.Command, v string) error {
			*destination = &sync.Due{String: &v}
			return nil
		},
	}
}

func newDeadlineFlag(destination **sync.Deadline) cli.Flag {
	return &cli.TimestampFlag{
		Name:     "deadline",
		Usage:    "Deadline date",
		OnlyOnce: true,
		Config: cli.TimestampConfig{
			Layouts: []string{time.DateOnly},
		},
		Action: func(ctx context.Context, cmd *cli.Command, v time.Time) error {
			*destination = &sync.Deadline{Date: &v}

			return nil
		},
	}
}

func newPriorityFlag(destination **int) cli.Flag {
	return &cli.IntFlag{
		Name:     "priority",
		Aliases:  []string{"p"},
		Usage:    "Task priority",
		Value:    1,
		OnlyOnce: true,
		Action: func(ctx context.Context, cmd *cli.Command, v int64) error {
			priority := int(v)
			*destination = &priority
			return nil
		},
	}
}

func newParentFlag(destination **string) cli.Flag {
	return &cli.StringFlag{
		Name:     "parent",
		Usage:    "Parent task ID",
		OnlyOnce: true,
		Action: func(ctx context.Context, cmd *cli.Command, v string) error {
			*destination = &v
			return nil
		},
	}
}

func newSectionFlag(destination **string) cli.Flag {
	return &cli.StringFlag{
		Name:     "section",
		Aliases:  []string{"s"},
		Usage:    "Section ID",
		OnlyOnce: true,
		Action: func(ctx context.Context, cmd *cli.Command, v string) error {
			*destination = &v
			return nil
		},
	}
}

func newLabelsFlag(destination *[]string) cli.Flag {
	return &cli.StringSliceFlag{
		Name:     "labels",
		Aliases:  []string{"l"},
		Usage:    "Labels",
		OnlyOnce: true,
		Action: func(ctx context.Context, cmd *cli.Command, v []string) error {
			destination = &v
			return nil
		},
	}
}

func newDurationFlag(destination **sync.Duration) cli.Flag {
	return &cli.StringFlag{
		Name:     "duration",
		Usage:    "Duration",
		OnlyOnce: true,
		Action: func(ctx context.Context, cmd *cli.Command, v string) error {
			duration, err := sync.ParseDuration(v)
			if err != nil {
				return err
			}
			*destination = duration
			return nil
		},
	}
}
