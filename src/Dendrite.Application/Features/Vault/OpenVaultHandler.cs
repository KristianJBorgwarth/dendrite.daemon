using Dendrite.Application.Services;
using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed class OpenVaultHandler : ICommandHandler<OpenVault, Unit>
{
    private readonly VaultOpener _opener;

    public OpenVaultHandler(VaultOpener opener)
    {
        _opener = opener;
    }

    public ValueTask<Unit> Handle(OpenVault command, CancellationToken cancellationToken)
    {
        _opener.Open(command.Path, activate: true);
        return ValueTask.FromResult(Unit.Value);
    }
}
