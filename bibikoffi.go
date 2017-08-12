package main

import (
	"context"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/containous/bibikoffi/internal/gh"
	"github.com/containous/bibikoffi/mjolnir"
	"github.com/containous/bibikoffi/types"
	"github.com/containous/flaeg"
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
			required(options.ConfigFilePath, "config-path")
			required(options.GitHubToken, "token")

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
	client := gh.NewGitHubClient(ctx, options.GitHubToken)

	return mjolnir.CloseIssues(client, ctx, config.Owner, config.RepositoryName, config.Rules, options.DryRun, options.Debug)
}

func required(field string, fieldName string) error {
	if len(field) == 0 {
		log.Fatalf("%s is mandatory.", fieldName)
	}
	return nil
}
