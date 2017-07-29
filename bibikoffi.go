package main

import (
	"context"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/containous/flaeg"
	"github.com/containous/myrmica-bibikoffi/mjolnir"
	"github.com/containous/myrmica-bibikoffi/types"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {

	options := &types.Options{
		DryRun:         true,
		Debug:          false,
		ConfigFilePath: "./bibikoffi.toml",
	}

	defaultPointersOptions := &types.Options{}
	rootCmd := &flaeg.Command{
		Name:                  "bibikoffi",
		Description:           `Myrmica Bibikoffi: Closes stale issues.`,
		Config:                options,
		DefaultPointersConfig: defaultPointersOptions,
		Run: func() error {
			if options.Debug {
				log.Printf("Run bibikoffi command with config : %+v\n", options)
			}

			if options.DryRun {
				log.Print("IMPORTANT: you are using the dry-run mode. Use `--dry-run=false` to disable this mode.")
			}

			err := process(options)
			if err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}

	flag := flaeg.New(rootCmd, os.Args[1:])
	flag.Run()
}

func process(options *types.Options) error {

	config := &types.Configuration{}
	meta, err := toml.DecodeFile(options.ConfigFilePath, config)

	if err != nil {
		return err
	}

	if options.Debug {
		log.Printf("configuration: %+v\n", meta)
	}

	ctx := context.Background()
	client := NewGitHubClient(ctx, config.GitHubToken)

	return mjolnir.CloseIssues(client, ctx, config.Owner, config.RepositoryName, config.Rules, options.DryRun, options.Debug)
}

// NewGitHubClient create a new GitHub client
func NewGitHubClient(ctx context.Context, token string) *github.Client {
	var client *github.Client
	if len(token) == 0 {
		client = github.NewClient(nil)
	} else {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}
	return client
}
