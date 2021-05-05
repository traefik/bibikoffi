package mjolnir

import (
	"context"

	"github.com/google/go-github/v28/github"
	"github.com/rs/zerolog/log"
	"github.com/traefik/bibikoffi/internal/search"
	"github.com/traefik/bibikoffi/types"
)

// LockIssues lock issues who match criterion.
func LockIssues(ctx context.Context, client *github.Client, owner, repositoryName string, ices []types.Frozen, dryRun bool) error {
	for _, ice := range ices {
		if !ice.Disable {
			err := lockIssues(ctx, client, owner, repositoryName, ice, dryRun)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func lockIssues(ctx context.Context, client *github.Client, owner, repositoryName string, ice types.Frozen, dryRun bool) error {
	oldIssues, err := search.FindIssues(ctx, client, owner, repositoryName,
		search.State("closed"),
		search.Lock(false),
		search.Cond(len(ice.ExcludedLabels) != 0, search.WithExcludedLabels(ice.ExcludedLabels...)),
		search.UpdatedBefore(ice.DaysSinceUpdate),
	)
	if err != nil {
		return err
	}

	log.Info().Msgf("Unlocked: %d", len(oldIssues))

	for i, issue := range oldIssues {
		logger := log.With().Int("issue", issue.GetNumber()).Logger()
		logger.Info().Msgf("%d Lock issue #%d: created %v, updated %v", i+1, issue.GetNumber(), issue.GetCreatedAt(), issue.GetUpdatedAt())

		logger.Debug().Msgf("%s %s", issue.GetTitle(), issue.GetHTMLURL())

		if dryRun {
			logger.Debug().Msgf("lock %s", ice.Label)
			continue
		}

		err = lockIssue(ctx, client, owner, repositoryName, issue.GetNumber(), ice.Label)
		if err != nil {
			logger.Error().Err(err).Msg("unable to lock issue")
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
