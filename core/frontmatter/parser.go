package frontmatter

import (
	"bytes"
	"errors"
	"log/slog"

	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	Title   string   `yaml:"title"`
	Tags    []string `yaml:"tags"`
	Created string   `yaml:"created"`
	Updated string   `yaml:"updated"`
	Author  string   `yaml:"author"`
}

func ParseTags(file []byte) ([]string, error) {
	fm, err := parseFrontMatter(file)
	if err != nil {
		return nil, err
	}
	slog.Debug("parsed front matter", "fm", fm)
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
	slog.Debug("yaml content to parse", "content", string(content))

	var fm FrontMatter
	if err := yaml.Unmarshal(content, &fm); err != nil {
		return nil, err
	}
	return &fm, nil
}
