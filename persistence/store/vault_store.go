package store

import persistenceconfig "github.com/KristianJBorgwarth/dendrite.daemon/persistence/config"

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
