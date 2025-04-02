package main

import (
	"fmt"

	"github.com/CnTeng/todoist-cli/cmd/todoist"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Print(err)
	}
}
