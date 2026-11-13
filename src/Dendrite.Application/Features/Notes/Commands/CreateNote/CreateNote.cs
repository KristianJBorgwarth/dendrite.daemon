using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record CreateNote : ICommand<string>, IVaultRequest
{
    public string Title { get; }
    public string TemplateName { get; }
    public string Directory { get; }
    public string? Vault { get; }

    public CreateNote(string title, string templateName, string directory, string? vault = null)
    {
        Title = title;
        TemplateName = templateName;
        Directory = directory;
        Vault = vault;
    }
}
