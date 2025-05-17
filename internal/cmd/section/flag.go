package section

import (
	"github.com/CnTeng/todoist-cli/internal/cmd/value"
	"github.com/spf13/pflag"
)

func newNameFlag(destination **string) *pflag.Flag {
	v := value.NewStringPtr(func(v string) error {
		*destination = &v
		return nil
	})
	return &pflag.Flag{
		Name:      "name",
		Shorthand: "n",
		Usage:     "Set the section name",
		Value:     v,
		DefValue:  v.String(),
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
		Usage:     "Assign the section to a specific project by <project-id>",
		Value:     v,
		DefValue:  v.String(),
	}
}
