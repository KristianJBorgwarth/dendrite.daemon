using Dendrite.Application.Services;
using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed record RebuildIndex : ICommand<Unit>;

public sealed class RebuildIndexHandler(IndexRebuilder rebuilder) : ICommandHandler<RebuildIndex, Unit>
{
    public ValueTask<Unit> Handle(RebuildIndex command, CancellationToken cancellationToken)
    {
        rebuilder.Rebuild();
        return ValueTask.FromResult(Unit.Value);
    }
}
