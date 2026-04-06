package store

import (
	"path"

	persistenceconfig "github.com/KristianJBorgwarth/dendrite.daemon/persistence/config"
)

type VaultStore struct {
	Config *persistenceconfig.VaultConfiguration
}

var vaultStore *VaultStore

func NewVaultStore(vault, templateDir string) *VaultStore {
	if vaultStore != nil {
		return vaultStore
	}
	vaultStore = &VaultStore{Config: persistenceconfig.NewVaultConfiguration(vault, templateDir)}
	return vaultStore
}

func GetVaultStore() *VaultStore {
	if vaultStore == nil {
		panic("VaultStore not initialized. Call NewVaultStore first.")
	}
	return vaultStore
}

func (vs *VaultStore) SetConfig(config persistenceconfig.VaultConfiguration) {
	vs.Config = &config
}

func (vs *VaultStore) GetTemplatePath(templateName string) string {
	templateName = vs.fileTypeCheck(templateName)
	return path.Join(vs.Config.TemplateDirectory(), templateName)
}

func (vs *VaultStore) fileTypeCheck(templateName string) string {
	fileType := path.Ext(templateName)
	if fileType != ".md" {
		return templateName + ".md"
	}
	return templateName
}

