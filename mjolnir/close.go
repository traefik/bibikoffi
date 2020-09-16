package mjolnir

import (
	"context"
	"log"

	"github.com/google/go-github/v28/github"
	"github.com/traefik/bibikoffi/internal/search"
	"github.com/traefik/bibikoffi/types"
)

// CloseIssues close issues who match criterion.
func CloseIssues(ctx context.Context, client *github.Client, owner, repositoryName string, rules []types.Rule, dryRun, debug bool) error {
	for _, rule := range rules {
		if !rule.Disable {
			err := closeIssuesByRule(ctx, client, owner, repositoryName, rule, dryRun, debug)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func closeIssuesByRule(ctx context.Context, client *github.Client, owner, repositoryName string, rule types.Rule, dryRun, debug bool) error {
	staleIssues, err := search.FindIssues(ctx, client, owner, repositoryName,
		search.State("open"),
		search.Cond(len(rule.IncludedLabels) != 0, search.WithLabels(rule.IncludedLabels...)),
		search.Cond(len(rule.ExcludedLabels) != 0, search.WithExcludedLabels(rule.ExcludedLabels...)),
		search.CreatedBefore(rule.DaysSinceCreation),
		search.UpdatedBefore(rule.DaysSinceUpdate),
	)
	if err != nil {
		return err
	}

	log.Printf("%v: %d\n", rule.IncludedLabels, len(staleIssues))

	for _, issue := range staleIssues {
		log.Printf("Close issue #%d: created %v, updated %v", issue.GetNumber(), issue.GetCreatedAt(), issue.GetUpdatedAt())

		if debug {
			log.Println(issue.GetTitle())
		}

		if dryRun {
			log.Println(rule.Message)
		} else {
			err = closeIssue(ctx, client, owner, repositoryName, issue.GetNumber(), rule.Message)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func closeIssue(ctx context.Context, client *github.Client, owner, repositoryName string, issueNumber int, comment string) error {
	issueRequest := &github.IssueRequest{
		State: github.String("closed"),
	}
	_, _, err := client.Issues.Edit(ctx, owner, repositoryName, issueNumber, issueRequest)
	if err != nil {
		return err
	}

	issueComment := &github.IssueComment{
		Body: github.String(comment),
	}
	_, _, err = client.Issues.CreateComment(ctx, owner, repositoryName, issueNumber, issueComment)
	return err
}
