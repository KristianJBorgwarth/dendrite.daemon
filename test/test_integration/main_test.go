package integration_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/rpc"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

var DbPath string = os.TempDir() + "/dendrite_test_vault"
var TestContext = context.Background()

func TestMain(m *testing.M) {

	err := persistence.InitializeIndex(DbPath)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.RemoveAll(DbPath)

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
