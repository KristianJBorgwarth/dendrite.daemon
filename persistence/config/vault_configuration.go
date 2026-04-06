package persistenceconfig

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

func (vc *VaultConfiguration) TemplateDirectory() string {
	return vc.templateDirectory
}
