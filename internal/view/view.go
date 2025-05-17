package view

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/CnTeng/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type View interface {
	Render() string
}

type InteractiveView interface {
	View
	Interact() error
}

var headerStyle = &table.CellStyle{
	CellAttrs: text.Colors{text.FgGreen, text.Underline},
}

func boolPtr(v bool) *bool {
	return &v
}

func tempFile(data string) ([]byte, error) {
	tmpFile, err := os.CreateTemp("", "todoist-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(data); err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to open editor: %w", err)
	}

	fmt.Println("Saving changes...")

	updatedFile, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read temp file: %w", err)
	}

	return updatedFile, nil
}
