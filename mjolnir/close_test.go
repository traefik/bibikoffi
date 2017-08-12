package mjolnir

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/containous/bibikoffi/types"
	"github.com/google/go-github/github"
)

const fixturesDir = "./test-fixtures"

func TestIntegrationBibikoffi(t *testing.T) {
	client := github.NewClient(nil)
	ctx := context.Background()

	config := &types.Configuration{}
	_, err := toml.DecodeFile(fixturePath("test01.toml"), config)
	if err != nil {
		t.Fatal(err)
	}

	err = CloseIssues(client, ctx, config.Owner, config.RepositoryName, config.Rules, true, true)
	if err != nil {
		t.Fatal(err)
	}
}

func fixturePath(filename string) string {
	return filepath.Join(fixturesDir, filename)
}
