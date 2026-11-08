using Dendrite.Daemon.FileHandling;
using Dendrite.Daemon.Persistence;
using Dendrite.Daemon.Persistence.Repositories;
using Dendrite.Daemon.Services;
using Mediator;

namespace Dendrite.Daemon.Features.Note;

public sealed record CreateNote(string Title, string TemplateName, string Directory) : ICommand<string>;

public sealed class CreateNoteHandler(
    UnitOfWork uow,
    NoteRepository notes,
    NoteMetadataService metadata,
    TemplateRenderer templates,
    VaultStore vaultStore) : ICommandHandler<CreateNote, string>
{
    public ValueTask<string> Handle(CreateNote command, CancellationToken cancellationToken)
    {
        var template = templates.Render(command.TemplateName, command.Title, command.Directory);
        var notePath = Path.GetFullPath(Path.Join(vaultStore.Config.VaultPath, command.Directory, template.Filename));

        using var tx = uow.Begin();

        if (notes.GetBySlug(template.Slug) is not null)
            return ValueTask.FromResult(notePath);

        var note = Models.Note.Create(notePath, command.Title, template.Slug);
        notes.Insert(note);
        metadata.WriteTags(note.Id, template.FrontMatter.Tags);
        metadata.WriteCustomFrontMatter(note.Id, template.FrontMatter.Custom);

        tx.Stage(notePath, template.Content);
        tx.Commit();

        return ValueTask.FromResult(notePath);
    }
}
