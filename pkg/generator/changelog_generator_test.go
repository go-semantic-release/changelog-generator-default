package generator

import (
	"encoding/json"
	"os"
	"testing"
	"text/template"

	"github.com/go-semantic-release/semantic-release/v2/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
	"github.com/stretchr/testify/require"
)

var testCommits = []*semrel.Commit{
	{},
	{
		SHA: "123456789", Type: "feat", Scope: "app", Message: "commit message",
		Annotations: map[string]string{"author_login": "test"},
	},
	{
		SHA: "deadbeef", Type: "fix", Scope: "", Message: "commit message",
		Annotations: map[string]string{"author_login": "test"},
	},
	{
		SHA: "87654321", Type: "ci", Scope: "", Message: "commit message",
		Annotations: map[string]string{"author_login": "test"},
	},
	{
		SHA: "43218765", Type: "build", Scope: "", Message: "commit message",
		Annotations: map[string]string{"author_login": "test"},
	},
	{
		SHA: "12345678", Type: "yolo", Scope: "swag", Message: "commit message",
	},
	{
		SHA: "12345678", Type: "chore", Scope: "", Message: "commit message",
		Raw:         []string{"", "BREAKING CHANGE: test"},
		Change:      &semrel.Change{Major: true},
		Annotations: map[string]string{"author_login": "test"},
	},
	{
		SHA: "12345679", Type: "chore!", Scope: "user", Message: "another commit message",
		Raw:    []string{"another commit message", "changed ID int into UUID"},
		Change: &semrel.Change{Major: true},
	},
	{
		SHA: "stop", Type: "chore", Scope: "", Message: "not included",
	},
}

var testChangelogConfig = &generator.ChangelogGeneratorConfig{
	Commits:       testCommits,
	LatestRelease: &semrel.Release{SHA: "stop"},
	NewVersion:    "2.0.0",
}

func TestDefaultGenerator(t *testing.T) {
	clGen := &DefaultChangelogGenerator{}
	require.NoError(t, clGen.Init(map[string]string{}))
	changelog := clGen.Generate(testChangelogConfig)

	require.Contains(t, changelog, "* **app:** commit message (12345678)")
	require.Contains(t, changelog, "* commit message (deadbeef)")
	require.Contains(t, changelog, "#### yolo")
	require.Contains(t, changelog, "#### Build")
	require.Contains(t, changelog, "#### CI")
	require.Contains(t, changelog, "```\nBREAKING CHANGE: test\n```")
	require.NotContains(t, changelog, "not included")
}

func TestEmojiGenerator(t *testing.T) {
	clGen := &DefaultChangelogGenerator{}
	require.NoError(t, clGen.Init(map[string]string{"emojis": "true"}))
	changelog := clGen.Generate(testChangelogConfig)

	require.Contains(t, changelog, "* **app:** commit message (12345678)")
	require.Contains(t, changelog, "* commit message (deadbeef)")
	require.Contains(t, changelog, "#### 🎁 Feature")
	require.Contains(t, changelog, "#### 🐞 Bug Fixes")
	require.Contains(t, changelog, "#### 🔁 CI")
	require.Contains(t, changelog, "#### 📦 Build")
	require.Contains(t, changelog, "#### 📣 Breaking Changes")
	require.Contains(t, changelog, "#### yolo")
	require.Contains(t, changelog, "```\nBREAKING CHANGE: test\n```")
	require.NotContains(t, changelog, "not included")
}

func TestFormatCommit(t *testing.T) {
	testCases := []struct {
		tpl            string
		commit         *semrel.Commit
		expectedOutput string
	}{
		{
			tpl:            defaultFormatCommitTemplateStr,
			commit:         &semrel.Commit{SHA: "123456789", Type: "feat", Scope: "", Message: "commit message"},
			expectedOutput: "* commit message (12345678)",
		},
		{
			tpl:            defaultFormatCommitTemplateStr,
			commit:         &semrel.Commit{SHA: "123", Type: "feat", Scope: "app", Message: "commit message"},
			expectedOutput: "* **app:** commit message (123)",
		},
		{
			tpl:            `* {{.SHA}} - {{.Message}} {{- with index .Annotations "author_login" }} [by @{{.}}] {{- end}}`,
			commit:         &semrel.Commit{SHA: "deadbeef", Type: "fix", Message: "custom template", Annotations: map[string]string{"author_login": "test"}},
			expectedOutput: "* deadbeef - custom template [by @test]",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.expectedOutput, func(t *testing.T) {
			tpl := template.Must(template.New("test").Funcs(templateFuncMap).Parse(tc.tpl))
			output := formatCommit(tpl, tc.commit)
			require.Equal(t, tc.expectedOutput, output)
		})
	}
}

func TestFormatCommitWithCustomTemplate(t *testing.T) {
	clGen := &DefaultChangelogGenerator{}
	require.NoError(t, clGen.Init(map[string]string{
		"format_commit_template": "* `{{ trimSHA .SHA}}` - {{.Message}} {{- with index .Annotations \"author_login\" }} [by @{{.}}] {{- end}}",
	}))
	changelog := clGen.Generate(testChangelogConfig)
	require.Contains(t, changelog, "* `12345678` - commit message [by @test]")
	require.NotContains(t, changelog, "* `deadbeef` - commit message (deadbeef) [by @test]")
}

func TestFormatCommitWithCustomTypes(t *testing.T) {
	myTypes := make(ChangelogTypes, len(defaultTypes))
	copy(myTypes, defaultTypes)
	for ix, clTy := range myTypes {
		if clTy.Type == "feat" {
			myTypes[ix].Text = "New Features"
			break
		}
	}
	typesFile, err := os.CreateTemp("", "changelog_types_*.json")
	require.NoError(t, err)
	defer os.Remove(typesFile.Name())
	err = json.NewEncoder(typesFile).Encode(myTypes)
	require.NoError(t, err)

	clGen := &DefaultChangelogGenerator{}
	require.NoError(t, clGen.Init(map[string]string{"types_path": typesFile.Name()}))
	changelog := clGen.Generate(testChangelogConfig)
	require.Contains(t, changelog, "#### New Features")
	require.NotContains(t, changelog, "#### Feature")
}
