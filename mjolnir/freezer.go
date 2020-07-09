package mjolnir

import (
	"context"
	"log"

	"github.com/containous/bibikoffi/internal/search"
	"github.com/containous/bibikoffi/types"
	"github.com/google/go-github/v28/github"
)

// LockIssues lock issues who match criterion.
func LockIssues(ctx context.Context, client *github.Client, owner, repositoryName string, ices []types.Frozen, dryRun, debug bool) error {
	for _, ice := range ices {
		if !ice.Disable {
			err := lockIssues(ctx, client, owner, repositoryName, ice, dryRun, debug)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func lockIssues(ctx context.Context, client *github.Client, owner, repositoryName string, ice types.Frozen, dryRun, debug bool) error {
	oldIssues, err := search.FindIssues(ctx, client, owner, repositoryName,
		search.State("closed"),
		search.Lock(false),
		search.Cond(len(ice.ExcludedLabels) != 0, search.WithExcludedLabels(ice.ExcludedLabels...)),
		search.UpdatedBefore(ice.DaysSinceUpdate),
	)
	if err != nil {
		return err
	}

	log.Printf("Unlocked: %d\n", len(oldIssues))

	for i, issue := range oldIssues {
		log.Printf("%d Lock issue #%d: created %v, updated %v", i+1, issue.GetNumber(), issue.GetCreatedAt(), issue.GetUpdatedAt())

		if debug {
			log.Println(issue.GetTitle(), issue.GetHTMLURL())
		}

		if dryRun {
			log.Println("lock ", ice.Label)
			continue
		}

		err = lockIssue(ctx, client, owner, repositoryName, issue.GetNumber(), ice.Label)
		if err != nil {
			return err
		}
	}

	return nil
}

func lockIssue(ctx context.Context, client *github.Client, owner, repositoryName string, issueNumber int, label string) error {
	_, err := client.Issues.Lock(ctx, owner, repositoryName, issueNumber, nil)
	if err != nil {
		return err
	}

	if label == "" {
		return nil
	}

	_, _, err = client.Issues.AddLabelsToIssue(ctx, owner, repositoryName, issueNumber, []string{label})
	return err
}
