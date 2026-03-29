package integration_test

import (
	"database/sql"
	"testing"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
	"github.com/KristianJBorgwarth/dendrite.daemon/test"
)


func TestCreateNoteHandlerOnSucess(t *testing.T) {

	// Arrange
	// TODO: move to test vars and main_test.go
	db, err := sql.Open("sqlite", test.NewTestVars().DbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	uow := repositories.NewUnitOfWork(db)

	handler := handlers.NewCreateNoteHandler(uow)

	// TODO: add to test vars and move to main_test.go
	ctx := test.Context()

	// Act


	// Assert
}
