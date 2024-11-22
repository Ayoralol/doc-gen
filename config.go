package main

type Config struct {
	Files     []File    `yaml:"files"`
	Output    Output    `yaml:"output"`
	Structure []Section `yaml:"structure"`
}

type File struct {
	Path string   `yaml:"path"`
	Type string   `yaml:"type"`
	Tags []string `yaml:"tags,omitempty"`
}

type Output struct {
	Individual string `yaml:"individual"`
	Aggregated string `yaml:"aggregated"`
}

type Section struct {
	Section string   `yaml:"section"`
	Files   []string `yaml:"files"`
}
