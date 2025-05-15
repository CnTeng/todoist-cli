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
		Usage:     "section name",
		Value:     v,
		DefValue:  v.String(),
	}
}
