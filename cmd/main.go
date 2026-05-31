package main

import (
	"log/slog"
	"os"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers/completion"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers/diagnostics"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers/note"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers/vault"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/logging"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/server"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/services"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
	_ "modernc.org/sqlite"
)

func main() {
	logging.Init()
	server := server.NewServer()

	uow := repositories.NewUnitOfWork();

	indexRepo := repositories.NewIndexRepository(persistence.NewReadContext())
	linkRepo := repositories.NewLinkRepository(persistence.NewReadContext())
	tagRepo := repositories.NewTagRepository(persistence.NewReadContext())
	noteRepo := repositories.NewNoteRepository(persistence.NewReadContext())
	cfeRepo := repositories.NewCfeRepository(persistence.NewReadContext())

	tagService := services.NewTagService(tagRepo)
	linkService := services.NewLinkService(linkRepo)
	noteService := services.NewNoteService(tagRepo, linkRepo, noteRepo, cfeRepo)
	cfeSvc := services.NewCfeService(cfeRepo)
	idxr := services.NewIndexRebuilder(uow, noteRepo, linkRepo, tagRepo, indexRepo, cfeRepo)

	server.RegisterHandler("vault/init", vault.NewInitializeHandler(idxr, noteService))
	server.RegisterHandler("vault/rebuild", vault.NewRebuildIndexHandler(idxr))

	server.RegisterHandler("note/create", note.NewCreateNoteHandler(uow, tagService, noteRepo, cfeSvc))
	server.RegisterHandler("note/save", note.NewSaveNoteHandler(uow, noteRepo, tagService, noteService, linkService, cfeSvc))
	server.RegisterHandler("note/delete", note.NewDeleteNoteHandler(uow, tagService, noteService))
	server.RegisterHandler("note/goto", note.NewGotoNoteHandler(noteRepo))
	server.RegisterHandler("note/backlinks", note.NewGetBackLinksHandler(linkRepo, noteRepo))
	server.RegisterHandler("note/search_by_tag", note.NewGetNotesByTagHandler(noteRepo))
	server.RegisterHandler("note/search_by_cf", note.NewGetNotesByCfeHandler(cfeRepo))

	server.RegisterHandler("completion/tag", completion.NewCompleteTagHandler(tagRepo))
	server.RegisterHandler("completion/slug", completion.NewCompleteSlugHandler(noteRepo))

	server.RegisterHandler("diagnostics/links", diagnostics.NewLinkDiagnosticHandler(noteRepo, linkRepo))

	if err := server.Run(os.Stdin, os.Stdout); err != nil {
		slog.Error("server error", "error", err)
	}
}
