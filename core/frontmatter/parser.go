package frontmatter

import (
	"bytes"
	"errors"
	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	Title   string   `yaml:"Title"`
	Tags    []string `yaml:"Tags"`
	Created string   `yaml:"Created"`
	Updated string   `yaml:"Updated"`
	Author  string   `yaml:"Author"`
}

func ParseTags(file []byte) ([]string, error) {
	fm, err := parseFrontMatter(file)
	if err != nil {
		return nil, err
	}
	return fm.Tags, nil
}

func parseFrontMatter(file []byte) (*FrontMatter, error) {
	content := bytes.TrimSpace(file)
	if bytes.HasPrefix(content, []byte("---")) {
		content = bytes.TrimSpace(content[3:])
	} else {
		return nil, errors.New("missing front matter delimiter")
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
