package search

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/go-github/github"
)

type byUpdated []github.Issue

func (a byUpdated) Len() int      { return len(a) }
func (a byUpdated) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byUpdated) Less(i, j int) bool {
	return a[i].GetUpdatedAt().Before(a[j].GetUpdatedAt())
}

// Parameter search parameter
type Parameter func() string

func FindStaleIssues(ctx context.Context, client *github.Client, owner string, repositoryName string, parameters ...Parameter) ([]github.Issue, error) {

	var filter string
	for _, param := range parameters {
		if param != nil {
			filter += param()
		}
	}

	query := fmt.Sprintf("repo:%s/%s type:issue state:open %s", owner, repositoryName, filter)
	log.Println(query)

	options := &github.SearchOptions{
		Sort:        "updated",
		Order:       "desc",
		ListOptions: github.ListOptions{PerPage: 25},
	}

	issues, err := findIssues(ctx, client, query, options)
	if err != nil {
		return nil, err
	}
	sort.Sort(byUpdated(issues))

	return issues, nil
}

func findIssues(ctx context.Context, client *github.Client, query string, searchOptions *github.SearchOptions) ([]github.Issue, error) {
	var allIssues []github.Issue
	for {
		issuesSearchResult, resp, err := client.Search.Issues(ctx, query, searchOptions)
		if err != nil {
			return nil, err
		}
		for _, issue := range issuesSearchResult.Issues {
			allIssues = append(allIssues, issue)
		}
		if resp.NextPage == 0 {
			break
		}
		searchOptions.Page = resp.NextPage
	}
	return allIssues, nil
}

func NoOp() string {
	return ""
}

// Cond apply conditionally some parameters
func Cond(apply bool, parameters ...Parameter) Parameter {
	if apply {
		return func() string {
			var filter string
			for _, param := range parameters {
				if param != nil {
					filter += param()
				}
			}
			return filter
		}
	}
	return NoOp
}

// WithLabels add a search filter by labels
func WithLabels(labels ...string) Parameter {
	var labelsFilter string
	for _, lbl := range labels {
		labelsFilter += fmt.Sprintf("label:%s ", lbl)
	}
	return func() string {
		return " " + labelsFilter
	}
}

// WithExcludedLabels add a search filter by unwanted labels
func WithExcludedLabels(labels ...string) Parameter {
	var labelsFilter string
	for _, lbl := range labels {
		labelsFilter += fmt.Sprintf("-label:%s ", lbl)
	}
	return func() string {
		return " " + labelsFilter
	}
}

func CreatedBefore(days int) Parameter {
	if days == 0 {
		return NoOp
	}

	return func() string {
		return fmt.Sprintf(" created:<%s ", daysToDate(days))
	}
}

func CreatedAfter(days int) Parameter {
	if days == 0 {
		return NoOp
	}

	return func() string {
		return fmt.Sprintf(" created:>=%s ", daysToDate(days))
	}
}

func UpdatedBefore(days int) Parameter {
	if days == 0 {
		return NoOp
	}

	return func() string {
		return fmt.Sprintf(" updated:<%s ", daysToDate(days))
	}
}

func UpdatedAfter(days int) Parameter {
	if days == 0 {
		return NoOp
	}

	return func() string {
		return fmt.Sprintf(" updated:>=%s ", daysToDate(days))
	}
}

// updated:>=2013-02-01
// created:<2011-01-01

const GitHubSearchDateLayout = "2006-01-02"

func daysToDate(days int) string {
	date := time.Now().Add(-time.Duration(days) * 24 * time.Hour)
	return date.Format(GitHubSearchDateLayout)
}
