using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed record OpenVault : ICommand<Unit>
{
    public string Path { get; }

    public OpenVault(string path)
    {
        Path = path;
    }
}
