package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/containous/bibikoffi/internal/gh"
	"github.com/containous/bibikoffi/mjolnir"
	"github.com/containous/bibikoffi/types"
	"github.com/containous/flaeg"
	"github.com/ogier/pflag"
)

func main() {
	options := &types.Options{
		DryRun:         true,
		Debug:          false,
		ConfigFilePath: "./bibikoffi.toml",
		ServerPort:     80,
	}

	rootCmd := &flaeg.Command{
		Name:                  "bibikoffi",
		Description:           `Myrmica Bibikoffi: Closes stale issues.`,
		DefaultPointersConfig: &types.Options{},
		Config:                options,
		Run:                   runCmd(options),
	}

	flag := flaeg.New(rootCmd, os.Args[1:])

	// version
	versionCmd := &flaeg.Command{
		Name:                  "version",
		Description:           "Display the version.",
		Config:                &types.NoOption{},
		DefaultPointersConfig: &types.NoOption{},
		Run: func() error {
			displayVersion()
			return nil
		},
	}
	flag.AddCommand(versionCmd)

	// Run command
	err := flag.Run()
	if err != nil && err != pflag.ErrHelp {
		log.Fatalf("Error: %v\n", err)
	}
}

func runCmd(options *types.Options) func() error {
	return func() error {
		if options.Debug {
			log.Printf("Run bibikoffi command with config : %+v\n", options)
		}

		if len(options.GitHubToken) == 0 {
			options.GitHubToken = os.Getenv("GITHUB_TOKEN")
		}

		err := required(options.GitHubToken, "token")
		if err != nil {
			return err
		}
		err = required(options.ConfigFilePath, "config-path")
		if err != nil {
			return err
		}

		if options.DryRun {
			log.Print("IMPORTANT: you are using the dry-run mode. Use `--dry-run=false` to disable this mode.")
		}

		return process(options)
	}
}

func process(options *types.Options) error {
	if options.ServerMode {
		server := &server{options: options}
		return server.ListenAndServe()
	}
	return launch(options)
}

func launch(options *types.Options) error {
	config := &types.Configuration{}
	metadata, err := toml.DecodeFile(options.ConfigFilePath, config)
	if err != nil {
		return err
	}

	if options.Debug {
		log.Printf("configuration: %+v\n", metadata)
	}

	ctx := context.Background()
	client := gh.NewGitHubClient(ctx, options.GitHubToken)

	err = mjolnir.CloseIssues(ctx, client, config.Owner, config.RepositoryName, config.Rules, options.DryRun, options.Debug)
	if err != nil {
		return err
	}

	return mjolnir.LockIssues(ctx, client, config.Owner, config.RepositoryName, config.Locks, options.DryRun, options.Debug)
}

func required(field string, fieldName string) error {
	if len(field) == 0 {
		return fmt.Errorf("%s is mandatory", fieldName)
	}
	return nil
}

type server struct {
	options *types.Options
}

func (s *server) ListenAndServe() error {
	return http.ListenAndServe(":"+strconv.Itoa(s.options.ServerPort), s)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Invalid http method: %s", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := launch(s.options)
	if err != nil {
		log.Printf("Report error: %v", err)
		http.Error(w, "Report error.", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Myrmica Bibikoffi: Scheluded.\n")
}
