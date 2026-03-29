package integration_test

import (
	"os"
	"testing"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/test"
)

func TestMain(m *testing.M) {

	dbPath := test.NewTestVars().DbPath

	err := persistence.InitializeIndex(dbPath)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.RemoveAll(dbPath)

	os.Exit(code)
}
