using Dendrite.Daemon.Persistence;
using Dendrite.Daemon.Services;
using Mediator;

namespace Dendrite.Daemon.Features.Vault;

public sealed record RebuildIndex : ICommand<Unit>;

public sealed class RebuildIndexHandler(IndexRebuilder rebuilder, VaultStore vaultStore)
    : ICommandHandler<RebuildIndex, Unit>
{
    public ValueTask<Unit> Handle(RebuildIndex command, CancellationToken cancellationToken)
    {
        rebuilder.Rebuild(vaultStore.Config.VaultPath);
        return ValueTask.FromResult(Unit.Value);
    }
}
