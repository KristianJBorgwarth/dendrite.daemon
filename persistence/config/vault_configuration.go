package persistenceconfig

import "path"

type VaultConfiguration struct {
	path              string
	templateDirectory string
	scratches         struct {
		directory    string
		templateName string
	}
	dailyNotes struct {
		directory    string
		templateName string
	}
}

func NewVaultConfiguration(vaultPath, templateDirectory, scratchDirectory, scratchTemplate, dailyDirectory, dailyTemplate string) *VaultConfiguration {
	return &VaultConfiguration{path: vaultPath, templateDirectory: templateDirectory, scratches: struct {
		directory    string
		templateName string
	}{
		directory:    scratchDirectory,
		templateName: scratchTemplate,
	}, dailyNotes: struct {
		directory    string
		templateName string
	}{
		directory:    dailyDirectory,
		templateName: dailyTemplate,
	}}
}

func (vc *VaultConfiguration) VaultPath() string {
	return vc.path
}

func (vc *VaultConfiguration) ScratchTemplatePath() string {
	fileName := vc.fileTypeCheck(vc.scratches.templateName)
	return path.Join(vc.templateDirectory, fileName)
}

func (vc *VaultConfiguration) DailyTemplatePath() string {
	fileName := vc.fileTypeCheck(vc.dailyNotes.templateName)
	return path.Join(vc.templateDirectory, fileName)
}

func (vc *VaultConfiguration) fileTypeCheck(templateName string) string {
	fileType := path.Ext(templateName)
	if fileType != ".md" {
		return templateName + ".md"
	}
	return templateName
}
