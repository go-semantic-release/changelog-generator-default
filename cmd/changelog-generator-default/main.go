package main

import (
	defaultGenerator "github.com/go-semantic-release/changelog-generator-default/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ChangelogGenerator: func() generator.ChangelogGenerator {
			return &defaultGenerator.DefaultChangelogGenerator{}
		},
	})
}
