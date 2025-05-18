package view

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ReorderView interface {
	View
	Reorder() ([]string, error)
}

type ReorderItem struct {
	ID          string
	Description string
}

type reorderView struct {
	items []*ReorderItem
}

func NewReorderView(items []*ReorderItem) ReorderView {
	return &reorderView{
		items: items,
	}
}

func (v *reorderView) Render() string {
	b := &strings.Builder{}
	for _, i := range v.items {
		fmt.Fprintf(b, "%s %s\n", i.ID, i.Description)
	}
	return b.String()
}

func (v *reorderView) Reorder() ([]string, error) {
	data, err := openEditorWithTempFile(v.Render())
	if err != nil {
		return nil, err
	}
	return v.parseEditedItems(data)
}

func (v *reorderView) parseEditedItems(data []byte) ([]string, error) {
	result := make([]string, 0, len(v.items))

	for line := range bytes.Lines(data) {
		if len(line) == 0 {
			continue
		}

		fields := bytes.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}

		result = append(result, string(fields[0]))
	}

	if len(result) != len(v.items) {
		return nil, fmt.Errorf("expected %d items, got %d", len(v.items), len(result))
	}

	return result, nil
}

func openEditorWithTempFile(data string) ([]byte, error) {
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

	editedFile, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read temp file: %w", err)
	}

	return editedFile, nil
}
