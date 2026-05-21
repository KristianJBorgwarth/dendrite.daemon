package persistenceconfig

type VaultConfiguration struct {
	name                   string
	path                   string
	templateDirectory      string
	exludeIndexFiles       []string
	defaultExcludes        []string
	overrideDefaultIgnores bool
}

func NewVaultConfiguration(
	vaultName,
	vaultPath,
	templateDirectory string,
	excludeIndexFiles []string,
	overrideDefaultIgnores bool,
) *VaultConfiguration {
	return &VaultConfiguration{
		name:              vaultName,
		path:              vaultPath,
		templateDirectory: templateDirectory,
		exludeIndexFiles:  excludeIndexFiles,
		defaultExcludes:   []string{".git", ".templates"},
	}
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

func (vc *VaultConfiguration) ExcludeIndexFiles() []string {
	return vc.exludeIndexFiles
}

func (vc *VaultConfiguration) GetExcludeIndexFiles() []string {
	defaultIgnores := []string{".git", ".templates"}
	if vc.overrideDefaultIgnores {
		return vc.exludeIndexFiles
	}
	return append(defaultIgnores, vc.exludeIndexFiles...)
}
