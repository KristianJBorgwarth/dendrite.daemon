package persistenceconfig

type VaultConfiguration struct {
	name              string
	path              string
	templateDirectory string
}

func NewVaultConfiguration(vaultName, vaultPath, templateDirectory string) *VaultConfiguration {
	return &VaultConfiguration{name: vaultName, path: vaultPath, templateDirectory: templateDirectory}
}

func (vc *VaultConfiguration) VaultName() string {
	return vc.name
}

func (vc *VaultConfiguration) VaultPath() string {
	return vc.path
}

func (vc *VaultConfiguration) TemplateDirectory() string {
	return vc.templateDirectory
}
