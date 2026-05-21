package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/rpc"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/store"
	_ "modernc.org/sqlite"
)

type DBFixture struct {
	DB          *sql.DB
	DBPath      string
	TestContext context.Context
	VaultStore  *store.VaultStore
}

func NewDBFixture() *DBFixture {
	vaultPath := os.TempDir()

	err := persistence.InitializeDBContext(vaultPath)
	if err != nil {
		panic(err)
	}

	xdg := os.Getenv("XDG_DATA_HOME")
	if xdg == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		xdg = filepath.Join(home, ".local", "share")
	}
	dbPath := filepath.Join(xdg, "dendrite", vaultPath, "index.db")

	dbContext := persistence.GetDBContext()

	return &DBFixture{
		DB:          dbContext.DB,
		DBPath:      dbPath,
		TestContext: context.Background(),
		VaultStore:  store.NewVaultStore("testVault", vaultPath, path.Join(vaultPath, "templates"), []string{}, false),
	}
}

var Fixture = NewDBFixture()

func TestMain(m *testing.M) {
	code := m.Run()

	os.RemoveAll(Fixture.DBPath)

	os.Exit(code)
}

func CreateTestRequest(method string, ID int, params json.RawMessage) *rpc.Request {
	return &rpc.Request{
		Jsonrpc: "2.0",
		ID:      &ID,
		Method:  method,
		Params:  params,
	}
}
