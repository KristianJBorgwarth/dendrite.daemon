package persistenceconfig

import "path"

type VaultConfiguration struct {
	path              string
	templateDirectory string
}

func NewVaultConfiguration(vaultPath, templateDirectory string) *VaultConfiguration {
	return &VaultConfiguration{path: vaultPath, templateDirectory: templateDirectory}
}

func (vc *VaultConfiguration) VaultPath() string {
	return vc.path
}

func (vc *VaultConfiguration) fileTypeCheck(templateName string) string {
	fileType := path.Ext(templateName)
	if fileType != ".md" {
		return templateName + ".md"
	}
	return templateName
}
