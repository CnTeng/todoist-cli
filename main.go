package main

import (
	"log"

	"github.com/CnTeng/todoist-cli/cmd/todoist"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
