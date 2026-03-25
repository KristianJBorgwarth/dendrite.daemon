package frontmatter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

type FrontMatter struct {
	Title   string
	Tags    []string
	Created string
	Updated string
	Author  string
}

func ParseFrontMatter(r io.Reader) (map[string]string, error) {
	scanner := bufio.NewScanner(r)
	frontMatter := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			break
		}
		parts := bytes.SplitN([]byte(line), []byte(":"), 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid front matter format")
		}
		key := string(bytes.TrimSpace(parts[0]))
		value := string(bytes.TrimSpace(parts[1]))
		frontMatter[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return frontMatter, nil
}

func ExtractTags(frontMatter map[string]string) ([]string, error) {
	tagsStr, ok := frontMatter["tags"]
	if !ok {
		return nil, errors.New("tags not found in front matter")
	}
	var tags []string
	err := json.Unmarshal([]byte(tagsStr), &tags)
	if err != nil {
		return nil, err
	}
	return tags, nil
}
