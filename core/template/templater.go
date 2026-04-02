package template

import (
	"os"
)


func GenerateTemplate(templatePath string) ([]byte, error) {
	if templatePath == "" {
		return nil, nil
	}
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}
