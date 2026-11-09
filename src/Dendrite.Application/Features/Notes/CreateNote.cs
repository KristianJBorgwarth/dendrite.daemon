using Dendrite.Application.Abstractions;
using Dendrite.Application.Services;
using Dendrite.Domain.Models;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record CreateNote(string Title, string TemplateName, string Directory) : ICommand<string>;

public sealed class CreateNoteHandler(
    IUnitOfWork uow,
    INoteRepository notes,
    NoteMetadataService metadata,
    TemplateRenderer templates,
    VaultContext vault) : ICommandHandler<CreateNote, string>
{
    public ValueTask<string> Handle(CreateNote command, CancellationToken cancellationToken)
    {
        var template = templates.Render(command.TemplateName, command.Title, command.Directory);
        var notePath = Path.GetFullPath(Path.Join(vault.Config.VaultPath, command.Directory, template.Filename));

        using var tx = uow.Begin();

        if (notes.GetBySlug(template.Slug) is not null)
            return ValueTask.FromResult(notePath);

        var note = Note.Create(notePath, command.Title, template.Slug);
        notes.Insert(note);
        metadata.WriteTags(note.Id, template.FrontMatter.Tags);
        metadata.WriteCustomFrontMatter(note.Id, template.FrontMatter.Custom);

        tx.Stage(notePath, template.Content);
        tx.Commit();

        return ValueTask.FromResult(notePath);
    }
}
