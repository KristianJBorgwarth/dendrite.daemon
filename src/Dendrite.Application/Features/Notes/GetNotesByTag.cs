using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record GetNotesByTag : IQuery<List<NoteDto>>, IVaultRequest
{
    public string Tag { get; }
    public string? Vault { get; }

    public GetNotesByTag(string tag, string? vault = null)
    {
        Tag = tag;
        Vault = vault;
    }
}
