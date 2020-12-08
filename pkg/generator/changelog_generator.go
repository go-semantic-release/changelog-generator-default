package generator

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-semantic-release/semantic-release/v2/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
)

func trimSHA(sha string) string {
	if len(sha) < 9 {
		return sha
	}
	return sha[:8]
}

func formatCommit(c *semrel.Commit) string {
	ret := "* "
	if c.Scope != "" {
		ret += fmt.Sprintf("**%s:** ", c.Scope)
	}
	ret += fmt.Sprintf("%s (%s)\n", c.Message, trimSHA(c.SHA))
	return ret
}

var CGVERSION = "dev"

type DefaultChangelogGenerator struct {
	prettifyOutput bool
}

func (g *DefaultChangelogGenerator) Init(m map[string]string) error {
	prettifyOutput := false

	prettifyConfig := m["prettify"]

	if prettifyConfig == "true" {
		prettifyOutput = true
	}

	g.prettifyOutput = prettifyOutput

	return nil
}

func (g *DefaultChangelogGenerator) Name() string {
	return "default"
}

func (g *DefaultChangelogGenerator) Version() string {
	return CGVERSION
}

func (g *DefaultChangelogGenerator) Generate(changelogConfig *generator.ChangelogGeneratorConfig) string {
	ret := fmt.Sprintf("## %s (%s)\n\n", changelogConfig.NewVersion, time.Now().UTC().Format("2006-01-02"))
	clTypes := NewChangelogTypes()
	for _, commit := range changelogConfig.Commits {
		if changelogConfig.LatestRelease.SHA == commit.SHA {
			break
		}
		if commit.Change != nil && commit.Change.Major {
			bc := fmt.Sprintf("%s```\n%s\n```\n", formatCommit(commit), strings.Join(commit.Raw[1:], "\n"))
			clTypes.AppendContent("%%bc%%", bc)
			continue
		}
		if commit.Type == "" {
			continue
		}
		clTypes.AppendContent(commit.Type, formatCommit(commit))
	}
	for _, ct := range clTypes {
		if ct.Content == "" {
			continue
		}
		prettifyPrefix := ""
		if g.prettifyOutput {
			prettifyPrefix = ct.Emoji
		}
		ret += fmt.Sprintf("#### %s%s\n\n%s\n", prettifyPrefix, ct.Text, ct.Content)
	}
	return ret
}
