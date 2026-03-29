package integration_test

import (
	"os"
	"testing"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

func TestMain(m *testing.M) {

	currentDir := os.TempDir() + "/dendrite_test_vault"

	err := persistence.InitializeIndex(currentDir)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.RemoveAll(currentDir)

	os.Exit(code)
}
