package generator

type ChangelogType struct {
	Type    string
	Text    string
	Content string
	Emoji   string
}

type ChangelogTypes []ChangelogType

func NewChangelogTypes(overrideTypes *ChangelogTypes) ChangelogTypes {
	src := &defaultTypes
	if overrideTypes != nil {
		src = overrideTypes
	}
	ret := make(ChangelogTypes, len(*src))
	copy(ret, *src)
	return ret
}

func (ct *ChangelogTypes) AppendContent(cType, content string) {
	for i, cct := range *ct {
		if cct.Type == cType {
			(*ct)[i].Content += content + "\n"
			return
		}
	}
	*ct = append(*ct, ChangelogType{
		Type:    cType,
		Text:    cType,
		Content: content + "\n",
	})
}

var defaultTypes = ChangelogTypes{
	{
		Type:  "%%bc%%",
		Text:  "Breaking Changes",
		Emoji: "📣",
	},
	{
		Type:  "feat",
		Text:  "Feature",
		Emoji: "🎁",
	},
	{
		Type:  "fix",
		Text:  "Bug Fixes",
		Emoji: "🐞",
	},
	{
		Type:  "revert",
		Text:  "Reverts",
		Emoji: "🔙",
	},
	{
		Type:  "perf",
		Text:  "Performance Improvements",
		Emoji: "📈",
	},
	{
		Type:  "docs",
		Text:  "Documentation",
		Emoji: "📄",
	},
	{
		Type:  "test",
		Text:  "Tests",
		Emoji: "🔎",
	},
	{
		Type:  "refactor",
		Text:  "Code Refactoring",
		Emoji: "🔀",
	},
	{
		Type:  "style",
		Text:  "Styles",
		Emoji: "🎨",
	},
	{
		Type:  "chore",
		Text:  "Chores",
		Emoji: "🚧",
	},
	{
		Type:  "build",
		Text:  "Build",
		Emoji: "📦",
	},
	{
		Type:  "ci",
		Text:  "CI",
		Emoji: "🔁",
	},
}
