using Dendrite.Application.Abstractions;
using Dendrite.Application.Services;
using Dendrite.Application.Vaults;
using Dendrite.Domain.FileHandling;
using Dendrite.Domain.Models;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed class CreateNoteHandler(
    IUnitOfWork uow, 
    INoteRepository notes, 
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

        notes.Add(Note.Create(new NoteFile(notePath, command.Title, template.Slug, template.FrontMatter, [])));

        tx.Stage(notePath, template.Content);
        tx.Commit();

        return ValueTask.FromResult(notePath);
    }
}
