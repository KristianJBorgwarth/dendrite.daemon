using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record GotoNote : IQuery<GotoNoteResult?>, IVaultRequest
{
    public string Link { get; }
    public string? Vault { get; }

    public GotoNote(string link, string? vault = null)
    {
        Link = link;
        Vault = vault;
    }
}
