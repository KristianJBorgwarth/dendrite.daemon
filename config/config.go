// Package config provides configuration management for the application, allowing for easy loading and parsing of configuration files, environment variables, and command-line arguments. It supports various formats such as JSON, YAML, and TOML, and provides a unified interface for accessing configuration values throughout the application.
package config

type ScratchConfig struct {
	Dir string
	TemplateName string
}

type DailyConfig struct {
	Dir string
	TemplateName string
	FilenameFormat string 
}

type Config struct {
	VaultPath string
	TemplateDir string
	ScratchNote ScratchConfig
	DailyNote DailyConfig
}

func (c *Config) SetDefaults() {
	if c.TemplateDir == "" {
		c.TemplateDir = "templates"
	}

	if c.ScratchNote.Dir == "" {
		c.ScratchNote.Dir = "scratch"
	}

	if c.ScratchNote.TemplateName == "" {
		c.ScratchNote.TemplateName = "scratch.md"
	}

	if c.DailyNote.Dir == "daily" {
		c.DailyNote.Dir = "daily"
	}

	if c.DailyNote.TemplateName == "" {
		c.DailyNote.TemplateName = "daily.md"
	}

	if c.DailyNote.FilenameFormat == "" {
		c.DailyNote.FilenameFormat = "2006-01-02.md"
	}
}
