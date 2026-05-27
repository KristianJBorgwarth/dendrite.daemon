package persistenceconfig

type DailyNotes struct {
	dir            string
	filenameFormat string
	templateName   string
}

func NewDailyNotes(dir, filenameFormat, templateName string) DailyNotes {
	return DailyNotes{
		dir:            dir,
		filenameFormat: filenameFormat,
		templateName:   templateName,
	}
}

type VaultConfig struct {
	vaultName              string
	vaultPath              string
	templatesDir           string
	excludeIndexFiles      []string
	overrideDefaultIgnores bool
	dailyNotes             DailyNotes
}

func NewVaultConfiguration(
	vaultName,
	vaultPath,
	templatesDir string,
	excludeIndexFiles []string,
	overrideDefaultIgnores bool,
	dailyNotes DailyNotes,
) *VaultConfig {
	return &VaultConfig{
		vaultName:    vaultName,
		vaultPath:    vaultPath,
		templatesDir: templatesDir,
		dailyNotes:   dailyNotes,
	}
}

func (vc *VaultConfig) VaultName() string {
	return vc.vaultName
}

func (vc *VaultConfig) VaultPath() string {
	return vc.vaultPath
}

func (vc *VaultConfig) TemplateDirectory() string {
	return vc.templatesDir
}

func (vc *VaultConfig) ExcludeIndexFiles() []string {
	defaultIgnores := []string{"index.md", "index"}
	if vc.overrideDefaultIgnores {
		return vc.excludeIndexFiles
	}
	return append(defaultIgnores, vc.excludeIndexFiles...)
}

