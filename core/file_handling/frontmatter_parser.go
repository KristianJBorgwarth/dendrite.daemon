package filehandling

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	Title   string   `yaml:"title"`
	Tags    []string `yaml:"tags"`
	Created string   `yaml:"created"`
	Updated string   `yaml:"updated"`
	Date		string   `yaml:"date"`
	Author  string   `yaml:"author"`
}

func ParseFrontMatter(file []byte) (*FrontMatter, error) {
	content := bytes.TrimSpace(file)
	if bytes.HasPrefix(content, []byte("---")) {
		content = bytes.TrimSpace(content[3:])
	} else {
		return nil, nil
	}

	if idx := bytes.Index(content, []byte("---")); idx != -1 {
		content = bytes.TrimSpace(content[:idx])
	}

	var fm FrontMatter
	if err := yaml.Unmarshal(content, &fm); err != nil {
		return nil, err
	}
	return &fm, nil
}
