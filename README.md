# :memo: changelog-generator-default
[![CI](https://github.com/go-semantic-release/changelog-generator-default/workflows/CI/badge.svg?branch=master)](https://github.com/go-semantic-release/changelog-generator-default/actions?query=workflow%3ACI+branch%3Amaster)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-semantic-release/changelog-generator-default)](https://goreportcard.com/report/github.com/go-semantic-release/changelog-generator-default)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-semantic-release/changelog-generator-default)](https://pkg.go.dev/github.com/go-semantic-release/changelog-generator-default)

The default changelog generator for [go-semantic-release](https://github.com/go-semantic-release/semantic-release).

## Output of the changelog

The changelog generator will order the types of commits in the changelog in the following order:
- Breaking Changes
- Feature
- Bug Fixes
- Reverts
- Performance Improvements
- Documentation
- Tests
- Code Refactoring
- Styles
- Chores
- Build
- CI

## Emoji Changelogs

In order to use emoji changelogs including a prefixed emoji, you need to provide the following config when calling semantic-relase: `--changelog-generator-opt "emojis=true"`. Or add the config within your `.semrelrc` file.

[Example Change Log](./examples/GENERATED_CHANGELOG.md)


## Format Commit Template

The plugin allows to specify the template which is used to render commits with the `--changelog-generator-opt` CLI flag, e.g., `--changelog-generator-opt "format_commit_template={{.Message}}"`. Or by adding the `format_commit_template` option within your `.semrelrc` file.


The following variables are available:

| Variable     | Description                                                                             |
|--------------|-----------------------------------------------------------------------------------------|
| .SHA         | The commit SHA (e.g., c862ecee7682be648289579b515dbc03a5357c89).                        |
| .Type        | The type of the commit (e.g., feat, fix, chore, etc).                                   |
| .Scope       | The scope of the commit.                                                                |
| .Message     | The first line of the commit message.                                                   |
| .Raw         | The raw commit message as a string array representing each line of the commit.          |
| .Annotations | A map containing different commit annotations like the `author_name` or `author_email`. |

Additionally, the following functions are available:

| Function | Description                                     |
|----------|-------------------------------------------------|
| trimSHA  | Trims the SHA to the its first eight characters |

### Examples:

| Template                                                                                                                                       | Example Output                                             |
|------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------|
| `* {{with .Scope -}} **{{.}}:** {{end}} {{- .Message}} ({{trimSHA .SHA}})`                                                                     | `* **app:** commit message (12345678)`                     |
| `* {{with .Scope -}} **{{.}}:** {{end}} {{- .Message}} ({{trimSHA .SHA}}) {{- with index .Annotations."author_login" }} - by @{{.}} {{- end}}` | `* **app:** commit message (12345678) - by @commit-author` |


## Licence

The [MIT License (MIT)](http://opensource.org/licenses/MIT)

Copyright © 2020 [Christoph Witzko](https://twitter.com/christophwitzko)
