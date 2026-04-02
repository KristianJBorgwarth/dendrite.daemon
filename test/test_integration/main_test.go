package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"testing"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/rpc"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

type TestFixture struct {
	Db *sql.DB
	DbPath string
	TestContext context.Context
}

func NewTestFixture() (*TestFixture, error) {
	dbPath := os.TempDir() + "/dendrite_test_vault"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	return &TestFixture{
		Db: db,
		DbPath: dbPath,
		TestContext: context.Background(),
	}, nil
}

var Fixture, err = NewTestFixture()

func TestMain(m *testing.M) {

	err := persistence.InitializeIndex(Fixture.DbPath)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.RemoveAll(Fixture.DbPath)

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


