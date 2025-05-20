package util

import (
	"context"
	"fmt"
	"time"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/spf13/cobra"
)

func newCompletionFunc(n int, cmpFunc func(ctx context.Context, args []string) ([]cobra.Completion, cobra.ShellCompDirective)) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if n > 0 && len(args) >= n {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return cmpFunc(cmd.Context(), args)
	}
}

func (f *Factory) NewColorCompletionFunc(n int) cobra.CompletionFunc {
	return newCompletionFunc(n, func(context.Context, []string) ([]cobra.Completion, cobra.ShellCompDirective) {
		colors := sync.ListColors()
		cmps := make([]cobra.Completion, 0, len(colors))
		for _, c := range colors {
			cmps = append(cmps, cobra.CompletionWithDesc(string(c), c.Hex()))
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	})
}

func (f *Factory) NewDeadlineCompletionFunc(n int) cobra.CompletionFunc {
	return newCompletionFunc(n, func(context.Context, []string) ([]cobra.Completion, cobra.ShellCompDirective) {
		cmps := []cobra.Completion{
			cobra.CompletionWithDesc(time.Now().Format(time.DateOnly), "Today"),
			cobra.CompletionWithDesc(time.Now().AddDate(0, 0, 1).Format(time.DateOnly), "Tomorrow"),
			cobra.CompletionWithDesc(time.Now().AddDate(0, 0, 7).Format(time.DateOnly), "Next week"),
			cobra.CompletionWithDesc(time.Now().AddDate(0, 1, 0).Format(time.DateOnly), "Next month"),
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	})
}

func (f *Factory) NewFavoriteCompletionFunc(n int) cobra.CompletionFunc {
	return newCompletionFunc(n, func(context.Context, []string) ([]cobra.Completion, cobra.ShellCompDirective) {
		cmps := []cobra.Completion{
			cobra.CompletionWithDesc("true", "Mark as favorite"),
			cobra.CompletionWithDesc("false", "Unmark as favorite"),
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	})
}

func (f *Factory) NewFilterCompletionFunc(n int) cobra.CompletionFunc {
	return newCompletionFunc(n, func(ctx context.Context, args []string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if err := f.LoadConfig(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.Dial(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		defer f.Close()

		filters := []*sync.Filter{}
		if err := f.CallResult(ctx, daemon.FilterList, nil, &filters); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		seen := make(map[string]struct{})
		for _, arg := range args {
			seen[arg] = struct{}{}
		}

		cmps := make([]cobra.Completion, 0, len(filters))
		for _, filter := range filters {
			if _, ok := seen[filter.ID]; ok {
				continue
			}
			desc := fmt.Sprintf("%s: %s", filter.Name, filter.Query)
			cmps = append(cmps, cobra.CompletionWithDesc(filter.ID, desc))
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	})
}

func (f *Factory) NewLabelCompletionFunc(n int) cobra.CompletionFunc {
	return newCompletionFunc(n, func(ctx context.Context, args []string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if err := f.LoadConfig(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.Dial(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		defer f.Close()

		labels := []*model.Label{}
		if err := f.CallResult(ctx, daemon.LabelList, nil, &labels); err != nil {
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
	})
}

func (f *Factory) NewPriorityCompletionFunc(n int) cobra.CompletionFunc {
	return newCompletionFunc(n, func(context.Context, []string) ([]cobra.Completion, cobra.ShellCompDirective) {
		cmps := []cobra.Completion{
			cobra.CompletionWithDesc("1", "P1 Natural"),
			cobra.CompletionWithDesc("2", "P2 Medium"),
			cobra.CompletionWithDesc("3", "P3 High"),
			cobra.CompletionWithDesc("4", "P4 Urgent"),
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	})
}

func (f *Factory) NewProjectCompletionFunc(n int, params *model.ProjectListArgs) cobra.CompletionFunc {
	return newCompletionFunc(n, func(ctx context.Context, args []string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if err := f.LoadConfig(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.Dial(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		defer f.Close()

		projects := []*sync.Project{}
		if err := f.CallResult(ctx, daemon.ProjectList, params, &projects); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		seen := make(map[string]struct{})
		for _, arg := range args {
			seen[arg] = struct{}{}
		}

		cmps := make([]cobra.Completion, 0, len(projects))
		for _, p := range projects {
			if _, ok := seen[p.ID]; ok {
				continue
			}
			desc := p.Name
			if p.Description != "" {
				desc = fmt.Sprintf("%s: %s", p.Name, p.Description)
			}
			cmps = append(cmps, cobra.CompletionWithDesc(p.ID, desc))
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	})
}

func (f *Factory) NewSectionCompletionFunc(n int, params *model.SectionListArgs) cobra.CompletionFunc {
	return newCompletionFunc(n, func(ctx context.Context, args []string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if err := f.LoadConfig(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.Dial(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		defer f.Close()

		sections := []*model.Section{}
		if err := f.CallResult(ctx, daemon.SectionList, params, &sections); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		seen := make(map[string]struct{})
		for _, arg := range args {
			seen[arg] = struct{}{}
		}

		cmps := make([]cobra.Completion, 0, len(sections))
		for _, section := range sections {
			if _, ok := seen[section.ID]; ok {
				continue
			}
			desc := fmt.Sprintf("%s: %s", section.ProjectName, section.Name)
			cmps = append(cmps, cobra.CompletionWithDesc(section.ID, desc))
		}
		return cmps, cobra.ShellCompDirectiveNoFileComp
	})
}

func (f *Factory) NewTaskCompletionFunc(n int, params *model.TaskListArgs) cobra.CompletionFunc {
	return newCompletionFunc(n, func(ctx context.Context, args []string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if err := f.LoadConfig(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if err := f.Dial(); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		defer f.Close()

		tasks := []*model.Task{}
		if err := f.CallResult(ctx, daemon.TaskList, params, &tasks); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		seen := make(map[string]struct{})
		for _, arg := range args {
			seen[arg] = struct{}{}
		}

		cmps := make([]cobra.Completion, 0, len(tasks))
		for _, task := range tasks {
			if _, ok := seen[task.ID]; ok {
				continue
			}
			desc := fmt.Sprintf("%s: %s", task.ProjectName, task.Content)
			cmps = append(cmps, cobra.CompletionWithDesc(task.ID, desc))
		}

		return cmps, cobra.ShellCompDirectiveNoFileComp
	})
}
