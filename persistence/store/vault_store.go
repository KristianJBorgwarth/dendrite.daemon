package store

import (
	"path"

	persistenceconfig "github.com/KristianJBorgwarth/dendrite.daemon/persistence/config"
)

type VaultStore struct {
	Config *persistenceconfig.VaultConfig
}

var vaultStore *VaultStore

func NewVaultStore(
	vaultName,
	vaultPath,
	templateDir string,
	exludeIndexFiles []string,
	overrideDefaultIgnores bool,
	dailyDir string,
	dailyFilenameFormat string,
	dailyTemplateName string,
) *VaultStore {
	if vaultStore != nil {
		return vaultStore
	}
	vaultStore = &VaultStore{Config: persistenceconfig.NewVaultConfiguration(
		vaultName,
		vaultPath,
		templateDir,
		exludeIndexFiles,
		overrideDefaultIgnores,
		persistenceconfig.NewDailyNotes(dailyDir, dailyFilenameFormat, dailyTemplateName))}
	return vaultStore
}

func GetVaultStore() *VaultStore {
	if vaultStore == nil {
		panic("VaultStore not initialized. Call NewVaultStore first.")
	}
	return vaultStore
}

func (vs *VaultStore) GetTemplatePath(templateName string) string {
	templateName = vs.fileTypeCheck(templateName)
	path := path.Join(vs.Config.VaultPath(), vs.Config.TemplateDirectory(), templateName)
	return path
}

func (vs *VaultStore) fileTypeCheck(templateName string) string {
	fileType := path.Ext(templateName)
	if fileType != ".md" {
		return templateName + ".md"
	}
	return templateName
}

func (vs *VaultStore) GetExcludeIndexFiles() []string {
	return vs.Config.ExcludeIndexFiles()
}
