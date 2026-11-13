using Dendrite.Application.Services;
using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed class RebuildIndexHandler : ICommandHandler<RebuildIndex, Unit>
{
    private readonly IndexRebuilder _rebuilder;

    public RebuildIndexHandler(IndexRebuilder rebuilder)
    {
        _rebuilder = rebuilder;
    }

    public ValueTask<Unit> Handle(RebuildIndex command, CancellationToken cancellationToken)
    {
        _rebuilder.Rebuild();
        return ValueTask.FromResult(Unit.Value);
    }
}
