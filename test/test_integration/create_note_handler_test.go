package integration_test

import (
	"encoding/json"
	"testing"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

func TestCreateNoteHandlerOnSucess(t *testing.T) {
	// Arrange
	uow := repositories.NewUnitOfWork(Fixture.Db)

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
	response := handler.Handle(Fixture.TestContext, request.Params)

	// Assert
	if response.Error != nil {
		t.Fatalf("expected no error, got: %v", response.Error)
	}

	if response.Result == nil {
		t.Fatal("expected result, got nil")
	}

	if  *response.ID != 1 {
		t.Fatal("expected non-empty ID in response result")
	}
}
