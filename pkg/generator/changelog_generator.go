package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
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

var templateFuncMap = template.FuncMap{
	"trimSHA": trimSHA,
}

var defaultFormatCommitTemplateStr = `* {{with .Scope -}} **{{.}}:** {{end}} {{- .Message}} ({{trimSHA .SHA}})`

func formatCommit(tpl *template.Template, c *semrel.Commit) string {
	ret := &bytes.Buffer{}
	err := tpl.Execute(ret, c)
	if err != nil {
		panic(err)
	}
	return ret.String()
}

var CGVERSION = "dev"

type DefaultChangelogGenerator struct {
	emojis          bool
	changelogTypes  *ChangelogTypes
	formatCommitTpl *template.Template
}

func (g *DefaultChangelogGenerator) Init(m map[string]string) error {
	emojis := false
	emojiConfig := m["emojis"]
	if emojiConfig == "true" {
		emojis = true
	}
	g.emojis = emojis

	templateStr := defaultFormatCommitTemplateStr
	if tplStr := m["format_commit_template"]; tplStr != "" {
		templateStr = tplStr
	}

	parsedTemplate, err := template.New("commit-template").
		Funcs(templateFuncMap).
		Parse(templateStr)
	if err != nil {
		return fmt.Errorf("failed to parse commit template: %w", err)
	}
	g.formatCommitTpl = parsedTemplate

	if typesPath, ok := m["types_path"]; ok {
		f, err := os.Open(typesPath)
		if err != nil {
			return fmt.Errorf("failed to open types_path: %v", err)
		}
		defer f.Close()
		var changelogTypes ChangelogTypes
		if err := json.NewDecoder(f).Decode(&changelogTypes); err != nil {
			return fmt.Errorf("error parsing %s: %v", typesPath, err)
		}
		g.changelogTypes = &changelogTypes
	}

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
	clTypes := NewChangelogTypes(g.changelogTypes)
	for _, commit := range changelogConfig.Commits {
		if changelogConfig.LatestRelease.SHA == commit.SHA {
			break
		}
		if commit.Change != nil && commit.Change.Major {
			bc := fmt.Sprintf("%s\n```\n%s\n```", formatCommit(g.formatCommitTpl, commit), strings.Join(commit.Raw[1:], "\n"))
			clTypes.AppendContent("%%bc%%", bc)
			continue
		}
		if commit.Type == "" {
			continue
		}
		clTypes.AppendContent(commit.Type, formatCommit(g.formatCommitTpl, commit))
	}
	for _, ct := range clTypes {
		if ct.Content == "" {
			continue
		}
		emojiPrefix := ""
		if g.emojis && ct.Emoji != "" {
			emojiPrefix = fmt.Sprintf("%s ", ct.Emoji)
		}
		ret += fmt.Sprintf("#### %s%s\n\n%s\n", emojiPrefix, ct.Text, ct.Content)
	}
	return ret
}
