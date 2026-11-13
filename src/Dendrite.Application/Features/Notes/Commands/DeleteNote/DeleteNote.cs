using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record DeleteNote : ICommand<Unit>, IVaultPathRequest
{
    public string Path { get; }
    public string? Vault { get; }

    public DeleteNote(string path, string? vault = null)
    {
        Path = path;
        Vault = vault;
    }
}
