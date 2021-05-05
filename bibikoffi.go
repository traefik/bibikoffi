package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/containous/flaeg"
	"github.com/ogier/pflag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/traefik/bibikoffi/internal/gh"
	"github.com/traefik/bibikoffi/mjolnir"
	"github.com/traefik/bibikoffi/types"
)

func main() {
	options := &types.Options{
		DryRun:         true,
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
	if err != nil && !errors.Is(err, pflag.ErrHelp) {
		log.Fatal().Err(err).Msg("unable to start bibikoffi")
	}
}

func runCmd(options *types.Options) func() error {
	return func() error {
		setupLogger(options.DryRun, options.LogLevel)

		log.Debug().Msgf("Run bibikoffi command with config : %+v", options)

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
			log.Debug().Msg("IMPORTANT: you are using the dry-run mode. Use `--dry-run=false` to disable this mode.")
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

	log.Debug().Msgf("configuration: %+v", metadata)

	ctx := context.Background()
	client := gh.NewGitHubClient(ctx, options.GitHubToken)

	err = mjolnir.CloseIssues(ctx, client, config.Owner, config.RepositoryName, config.Rules, options.DryRun)
	if err != nil {
		return err
	}

	return mjolnir.LockIssues(ctx, client, config.Owner, config.RepositoryName, config.Locks, options.DryRun)
}

func required(field, fieldName string) error {
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
		log.Error().Msgf("Invalid http method: %s", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := launch(s.options)
	if err != nil {
		log.Error().Err(err).Msg("Report error")
		http.Error(w, "Report error.", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Myrmica Bibikoffi: Scheluded.\n")
}

// setupLogger is configuring the logger.
func setupLogger(dryRun bool, level string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Logger = zerolog.New(os.Stderr).With().Caller().Logger()

	logLevel := zerolog.DebugLevel
	if !dryRun {
		var err error
		logLevel, err = zerolog.ParseLevel(strings.ToLower(level))
		if err != nil {
			logLevel = zerolog.InfoLevel
		}
	}

	zerolog.SetGlobalLevel(logLevel)

	log.Trace().Msgf("Log level set to %s.", logLevel)
}
