package template

import (
	"os"
	"strings"
	"time"
)

func RenderTemplate(templatePath string, title string, slug string) ([]byte, error) {
	template, err := readTemplate(templatePath, title)
	if err != nil {
		return nil, err
	}

	r := strings.NewReplacer(
		"{{title}}", title,
		"{{date}}", time.Now().Format("2006-01-02"),
		"{{slug}}", slug,
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
