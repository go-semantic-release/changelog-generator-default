package generator

type ChangelogType struct {
	Type    string
	Text    string
	Content string
}

type ChangelogTypes []ChangelogType

func NewChangelogTypes() ChangelogTypes {
	ret := make(ChangelogTypes, len(defaultTypes))
	copy(ret, defaultTypes)
	return ret
}

func (ct *ChangelogTypes) AppendContent(cType, content string) {
	for i, cct := range *ct {
		if cct.Type == cType {
			(*ct)[i].Content += content
			return
		}
	}
	*ct = append(*ct, ChangelogType{
		Type:    cType,
		Text:    cType,
		Content: content,
	})
}

var defaultTypes = ChangelogTypes{
	{
		Type: "%%bc%%",
		Text: "Breaking Changes",
	},
	{
		Type: "feat",
		Text: "Feature",
	},
	{
		Type: "fix",
		Text: "Bug Fixes",
	},
	{
		Type: "revert",
		Text: "Reverts",
	},
	{
		Type: "perf",
		Text: "Performance Improvements",
	},
	{
		Type: "docs",
		Text: "Documentation",
	},
	{
		Type: "test",
		Text: "Tests",
	},
	{
		Type: "refactor",
		Text: "Code Refactoring",
	},
	{
		Type: "style",
		Text: "Styles",
	},
	{
		Type: "chore",
		Text: "Chores",
	},
	{
		Type: "build",
		Text: "Build",
	},
	{
		Type: "ci",
		Text: "CI",
	},
}
