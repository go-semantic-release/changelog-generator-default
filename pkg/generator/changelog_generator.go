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

var typeToText = map[string]string{
	"%%bc%%":   "Breaking Changes",
	"feat":     "Feature",
	"fix":      "Bug Fixes",
	"revert":   "Reverts",
	"perf":     "Performance Improvements",
	"docs":     "Documentation",
	"test":     "Tests",
	"refactor": "Code Refactoring",
	"style":    "Styles",
	"chore":    "Chores",
	"build":    "Build",
	"ci":       "CI",
}

var typeOrder = map[int]string{
	0:  "%%bc%%",
	1:  "feat",
	2:  "fix",
	3:  "revert",
	4:  "perf",
	5:  "docs",
	6:  "test",
	7:  "refactor",
	8:  "style",
	9:  "chore",
	10: "build",
	11: "ci",
}

func getSortedKeys(m *map[string]string) []string {
	types := make(map[int]string, len(*m))
	keys := make([]string, len(*m))

	i := 0
	for k := range *m {
		types[i] = k
		i++
	}

	i = 0
	for ki := 0; ki < len(typeOrder); ki++ {
		for ti := range types {
			if types[ti] == typeOrder[ki] {
				keys[i] = types[ti]
				delete(types, ti)
				i++
			}
		}
	}

	for ti := range types {
		keys[i] = types[ti]
		i++
	}

	return keys
}

var CGVERSION = "dev"

type DefaultChangelogGenerator struct{}

func (g *DefaultChangelogGenerator) Init(m map[string]string) error {
	return nil
}

func (g *DefaultChangelogGenerator) Name() string {
	return "default"
}

func (g *DefaultChangelogGenerator) Version() string {
	return CGVERSION
}

func (*DefaultChangelogGenerator) Generate(changelogConfig *generator.ChangelogGeneratorConfig) string {
	ret := fmt.Sprintf("## %s (%s)\n\n", changelogConfig.NewVersion, time.Now().UTC().Format("2006-01-02"))
	typeScopeMap := make(map[string]string)
	for _, commit := range changelogConfig.Commits {
		if changelogConfig.LatestRelease.SHA == commit.SHA {
			break
		}
		if commit.Change != nil && commit.Change.Major {
			typeScopeMap["%%bc%%"] += fmt.Sprintf("%s```\n%s\n```\n", formatCommit(commit), strings.Join(commit.Raw[1:], "\n"))
			continue
		}
		if commit.Type == "" {
			continue
		}
		typeScopeMap[commit.Type] += formatCommit(commit)
	}
	for _, t := range getSortedKeys(&typeScopeMap) {
		msg := typeScopeMap[t]
		typeName, found := typeToText[t]
		if !found {
			typeName = t
		}
		ret += fmt.Sprintf("#### %s\n\n%s\n", typeName, msg)
	}
	return ret
}
