package types

// Configuration global configuration
type Configuration struct {
	Owner          string
	RepositoryName string
	GitHubToken    string
	Rules          []Rule
}

// Rule to closes issues
type Rule struct {
	IncludedLabels    []string
	ExcludedLabels    []string
	Message           string
	DaysSinceCreation int
	DaysSinceUpdate   int
	Disable           bool
}

// Options CLI options
type Options struct {
	DryRun         bool   `long:"dry-run" description:"Dry run mode."`
	Debug          bool   `long:"debug" description:"Debug mode."`
	ConfigFilePath string `long:"config-path" description:"Configuration file path."`
}
