package types

// Configuration global configuration
type Configuration struct {
	Owner          string   `toml:"owner,omitempty"`
	RepositoryName string   `toml:"repositoryName,omitempty"`
	Rules          []Rule   `toml:"rules,omitempty"`
	Locks          []Frozen `toml:"locks,omitempty"`
}

// Frozen to lock issues
type Frozen struct {
	Label           string   `toml:"label,omitempty"`
	ExcludedLabels  []string `toml:"excludedLabels,omitempty"`
	DaysSinceUpdate int      `toml:"daysSinceUpdate,omitempty"`
	Disable         bool     `toml:"disable,omitempty"`
}

// Rule to closes issues
type Rule struct {
	IncludedLabels    []string `toml:"includedLabels,omitempty"`
	ExcludedLabels    []string `toml:"excludedLabels,omitempty"`
	Message           string   `toml:"message,omitempty"`
	DaysSinceCreation int      `toml:"daysSinceCreation,omitempty"`
	DaysSinceUpdate   int      `toml:"daysSinceUpdate,omitempty"`
	Disable           bool     `toml:"disable,omitempty"`
}

// Options CLI options
type Options struct {
	GitHubToken    string `long:"token" short:"t" description:"GitHub Token. [required]"`
	DryRun         bool   `long:"dry-run" description:"Dry run mode."`
	Debug          bool   `long:"debug" description:"Debug mode."`
	ConfigFilePath string `long:"config-path" description:"Configuration file path."`
	ServerMode     bool   `long:"server" description:"Server mode."`
	ServerPort     int    `long:"port" description:"Server port."`
}

// NoOption empty struct.
type NoOption struct{}
