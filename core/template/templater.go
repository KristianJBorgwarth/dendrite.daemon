package templatepackage template

import (
	"os"
	"strings"
	"time"
)

package template

import (
	"os"
	"strings"
	"time"
)

func RenderTemplate(templatePath string, title string, slug string) ([]byte, error) {
	tmpl, err := readTemplate(templatePath, title)
	if err != nil {
		return nil, err
	}

	r := strings.NewReplacer(
		"{{title}}", title,
		"{{date}}", time.Now().Format("2006-01-02"),
		"{{slug}}", slug,
	)

	return []byte(r.Replace(string(tmpl))), nil
}

func readTemplate(templatePath string, title string) ([]byte, error) {
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return createDefaultTemplate(title), nil
	}
	return os.ReadFile(templatePath)
}

func createDefaultTemplate(title string) []byte {
	return []byte("---\ntitle: " + title + "\ndate: " + time.Now().Format("2006-01-02") + "\n---\n")
}

