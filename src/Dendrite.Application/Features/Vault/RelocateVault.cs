using Mediator;

namespace Dendrite.Application.Features.Vault;

/// <summary>Points a registered vault at a new folder, keeping its index.</summary>
public sealed record RelocateVault : ICommand<Unit>
{
    public string Vault { get; }
    public string Path { get; }

    public RelocateVault(string vault, string path)
    {
        Vault = vault;
        Path = path;
    }
}
