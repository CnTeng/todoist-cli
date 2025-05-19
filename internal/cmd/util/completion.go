package util

import (
	"fmt"
	"time"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/spf13/cobra"
)

func (f *Factory) NewColorCompletionFunc(n int) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if n > 0 && len(args) >= n {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		colors := sync.ListColors()

		cmps := make([]cobra.Completion, len(colors))

		for i, c := range colors {
			cmps[i] = cobra.CompletionWithDesc(string(c), c.Hex())
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	}
}

func (f *Factory) NewDeadlineCompletionFunc(n int) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if n > 0 && len(args) >= n {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		cmps := []cobra.Completion{
			cobra.CompletionWithDesc(time.Now().Format(time.DateOnly), "Today"),
			cobra.CompletionWithDesc(time.Now().AddDate(0, 0, 1).Format(time.DateOnly), "Tomorrow"),
			cobra.CompletionWithDesc(time.Now().AddDate(0, 0, 7).Format(time.DateOnly), "Next week"),
			cobra.CompletionWithDesc(time.Now().AddDate(0, 1, 0).Format(time.DateOnly), "Next month"),
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	}
}

func (f *Factory) NewFavoriteCompletionFunc(n int) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if n > 0 && len(args) >= n {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		cmps := []cobra.Completion{
			cobra.CompletionWithDesc("true", "Add to favorites"),
			cobra.CompletionWithDesc("false", "Remove from favorites"),
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	}
}

func (f *Factory) NewTaskCompletionFunc(n int) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if n > 0 && len(args) >= n {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.LoadConfig(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.Dial(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		defer f.Close()

		params := &view.TaskViewConfig{}
		result := []*model.Task{}
		if err := f.CallResult(cmd.Context(), daemon.TaskList, params, &result); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		seen := make(map[string]struct{})
		for _, arg := range args {
			seen[arg] = struct{}{}
		}

		cmps := make([]cobra.Completion, len(result))
		for i, task := range result {
			if _, ok := seen[task.ID]; ok {
				continue
			}
			desc := fmt.Sprintf("%s: %s", task.ProjectName, task.Content)
			cmps[i] = cobra.CompletionWithDesc(task.ID, desc)
		}

		return cmps, cobra.ShellCompDirectiveNoFileComp
	}
}

func (f *Factory) NewProjectCompletionFunc(n int) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if n > 0 && len(args) >= n {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.LoadConfig(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.Dial(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		defer f.Close()

		result := []*sync.Project{}
		if err := f.CallResult(cmd.Context(), daemon.ProjectList, nil, &result); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		seen := make(map[string]struct{})
		for _, arg := range args {
			seen[arg] = struct{}{}
		}

		cmps := make([]cobra.Completion, len(result))
		for i, project := range result {
			if _, ok := seen[project.ID]; ok {
				continue
			}

			desc := project.Name
			if project.Description != "" {
				desc = fmt.Sprintf("%s: %s", project.Name, project.Description)
			}
			cmps[i] = cobra.CompletionWithDesc(project.ID, desc)
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	}
}

func (f *Factory) NewLabelCompletionFunc(n int) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if n > 0 && len(args) >= n {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.LoadConfig(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.Dial(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		defer f.Close()

		labels := []*model.Label{}
		if err := f.CallResult(cmd.Context(), daemon.LabelList, nil, &labels); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		seen := make(map[string]struct{})
		for _, arg := range args {
			seen[arg] = struct{}{}
		}

		cmps := make([]cobra.Completion, len(labels))
		for i, label := range labels {
			if _, ok := seen[label.Name]; ok {
				continue
			}

			desc := "personal"
			if label.IsShared {
				desc = "shared"
			}

			cmps[i] = cobra.CompletionWithDesc(label.Name, desc)
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	}
}

func (f *Factory) NewFilterCompletionFunc(n int) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if n > 0 && len(args) >= n {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.LoadConfig(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.Dial(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		defer f.Close()

		filters := []*sync.Filter{}
		if err := f.CallResult(cmd.Context(), daemon.FilterList, nil, &filters); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		seen := make(map[string]struct{})
		for _, arg := range args {
			seen[arg] = struct{}{}
		}

		cmps := make([]cobra.Completion, len(filters))
		for i, filter := range filters {
			if _, ok := seen[filter.ID]; ok {
				continue
			}

			desc := fmt.Sprintf("%s: %s", filter.Name, filter.Query)
			cmps[i] = cobra.CompletionWithDesc(filter.ID, desc)
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	}
}

func (f *Factory) NewSectionCompletionFunc(n int) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if n > 0 && len(args) >= n {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.LoadConfig(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.Dial(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		defer f.Close()

		sections := []*model.Section{}
		if err := f.CallResult(cmd.Context(), daemon.SectionList, nil, &sections); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		seen := make(map[string]struct{})
		for _, arg := range args {
			seen[arg] = struct{}{}
		}

		cmps := make([]cobra.Completion, len(sections))
		for i, section := range sections {
			if _, ok := seen[section.ID]; ok {
				continue
			}

			desc := fmt.Sprintf("%s: %s", section.ProjectName, section.Name)
			cmps[i] = cobra.CompletionWithDesc(section.ID, desc)
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	}
}

func (f *Factory) NewPriorityCompletionFunc(n int) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if n > 0 && len(args) >= n {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		cmps := []cobra.Completion{
			cobra.CompletionWithDesc("1", "P1 Natural"),
			cobra.CompletionWithDesc("2", "P2 Medium"),
			cobra.CompletionWithDesc("3", "P3 High"),
			cobra.CompletionWithDesc("4", "P4 Urgent"),
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	}
}
