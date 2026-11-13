using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record GetNotesByCf : IQuery<List<NoteDto>>, IVaultRequest
{
    public string Key { get; }
    public string Value { get; }
    public string? Vault { get; }

    public GetNotesByCf(string key, string value, string? vault = null)
    {
        Key = key;
        Value = value;
        Vault = vault;
    }
}
