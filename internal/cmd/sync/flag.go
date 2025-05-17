package sync

import (
	"fmt"
	"time"

	"github.com/CnTeng/todoist-cli/internal/cmd/value"
	"github.com/spf13/pflag"
)

func newSinceFlag(destination *time.Time) *pflag.Flag {
	v := value.NewTimeValue(time.DateOnly,
		func(v time.Time) error {
			*destination = v
			return nil
		})
	return &pflag.Flag{
		Name:      "since",
		Shorthand: "s",
		Usage:     fmt.Sprintf("Sync completed tasks since <YYYY-MM-DD> (default: %s)", destination.Format(time.DateOnly)),
		Value:     v,
		DefValue:  v.String(),
	}
}
