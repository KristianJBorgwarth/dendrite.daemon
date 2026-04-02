package integration_test

import (
	"encoding/json"
	"testing"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

func TestCreateNoteHandlerOnSucess(t *testing.T) {
	// Arrange
	uow := repositories.NewUnitOfWork(Fixture.DB)

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
	response, err := handler.Handle(Fixture.TestContext, request.Params)

	// Assert
	if err != nil {
		t.Fatalf("handler returned an error: %v", err)
	}

	if response != nil {
		t.Fatalf("expected response to be nil, got: %v", response)
	}
}
