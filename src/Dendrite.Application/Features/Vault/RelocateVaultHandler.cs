using Dendrite.Application.Services;
using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed class RelocateVaultHandler : ICommandHandler<RelocateVault, Unit>
{
    private readonly VaultRelocator _relocator;

    public RelocateVaultHandler(VaultRelocator relocator)
    {
        _relocator = relocator;
    }

    public ValueTask<Unit> Handle(RelocateVault command, CancellationToken cancellationToken)
    {
        _relocator.Relocate(command.Vault, command.Path);
        return ValueTask.FromResult(Unit.Value);
    }
}
