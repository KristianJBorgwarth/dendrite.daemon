using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record GetBacklinks : IQuery<List<BacklinkDto>>, IVaultPathRequest
{
    public string Path { get; }
    public string? Vault { get; }

    public GetBacklinks(string path, string? vault = null)
    {
        Path = path;
        Vault = vault;
    }
}
