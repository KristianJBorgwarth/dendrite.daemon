package integration_test

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

func TestCreateNoteHandlerOnSucess(t *testing.T) {
	// Arrange
	// TODO: move to test vars and main_test.go
	db, err := sql.Open("sqlite", DbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	uow := repositories.NewUnitOfWork(db)

	handler := handlers.NewCreateNoteHandler(uow)

	requestParams := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}{
		Title:   "Test Note",
		Content: "This is a test note.",
	}

	requestParamsBytes, err := json.Marshal(requestParams)
	if err != nil {
		t.Fatalf("failed to marshal request params: %v", err)
	}

	request := CreateTestRequest("createNote", 1, requestParamsBytes)

	// Act
	response, rpcError = handler.Handle(TestContext, request.Params)
}
