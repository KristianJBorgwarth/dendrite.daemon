using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Diagnostics;

public sealed record LinkDiagnostics : IQuery<List<LinkDiagnosticDto>>, IVaultPathRequest
{
    public string Path { get; }
    public string? Vault { get; }

    public LinkDiagnostics(string path, string? vault = null)
    {
        Path = path;
        Vault = vault;
    }
}
