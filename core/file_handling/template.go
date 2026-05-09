package filehandling

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/store"
)

type Template struct {
	Content     []byte
	Title       string
	Slug        string
	Filename    string
	FrontMatter *FrontMatter
}

func NewTemplate(templateName, title, directory string) (*Template, error) {
	var templatePath string
	if templateName != "" {
		templatePath = store.GetVaultStore().GetTemplatePath(templateName)
	}

	fileNameSlug := Slugify(title)

	pathSlug := filepath.ToSlash(filepath.Join(directory, fileNameSlug))
	filename := fileNameSlug + ".md"

	content, err := renderTemplate(templatePath, title, pathSlug, filename)
	if err != nil {
		return nil, err
	}

	fm, err := ParseFrontMatter(content)
	if err != nil {
		return nil, err
	}

	return &Template{
		Content:     content,
		Title:       title,
		Slug:        pathSlug,
		Filename:    filename,
		FrontMatter: fm,
	}, nil
}

func renderTemplate(templatePath, title, slug, filename string) ([]byte, error) {
	template, err := readTemplate(templatePath, title)
	if err != nil {
		return nil, err
	}

	r := strings.NewReplacer(
		"{{title}}", title,
		"{{date}}", time.Now().Format("2006-01-02"),
		"{{slug}}", slug,
		"{{file}}", filename,
	)

	return []byte(r.Replace(string(template))), nil
}

func readTemplate(templatePath string, title string) ([]byte, error) {
	if templatePath == "" {
		return createDefaultTemplate(title), nil
	}

	template, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func createDefaultTemplate(title string) []byte {
	return []byte("---\ntitle: " + title + "\ndate: " + time.Now().Format("2006-01-02") + "\n---\n")
}
