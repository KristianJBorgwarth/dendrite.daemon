package store

import persistenceconfig "github.com/KristianJBorgwarth/dendrite.daemon/persistence/config"

type VaultStore struct {
	Config *persistenceconfig.VaultConfiguration
}

var vaultStore *VaultStore

func NewVaultStore(vaultPath string) *VaultStore {
	if vaultStore != nil {
		return vaultStore
	}
	vaultStore = &VaultStore{}
	return vaultStore
}

func GetVaultStore() *VaultStore {
	if vaultStore == nil {
		panic("VaultStore not initialized. Call NewVaultStore first.")
	}
	return vaultStore
}
