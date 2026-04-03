package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/rpc"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	_ "modernc.org/sqlite"
)

type DBFixture struct {
	DB          *sql.DB
	DBPath      string
	TestContext context.Context
}

func NewDBFixture() *DBFixture {
	vaultPath := os.TempDir()

	db, err := persistence.InitializeIndex(vaultPath)
	if err != nil {
		panic(err)
	}

	dbPath := filepath.Join(os.TempDir(), ".index", "index.db")

	return &DBFixture{
		DB:          db,
		DBPath:      dbPath,
		TestContext: context.Background(),
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
