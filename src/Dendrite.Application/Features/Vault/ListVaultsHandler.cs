using Dendrite.Application.Abstractions;
using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed class ListVaultsHandler : IQueryHandler<ListVaults, List<VaultInfoDto>>
{
    private readonly VaultRegistry _registry;
    private readonly IIndexDatabase _indexes;

    public ListVaultsHandler(VaultRegistry registry, IIndexDatabase indexes)
    {
        _registry = registry;
        _indexes = indexes;
    }

    public ValueTask<List<VaultInfoDto>> Handle(ListVaults query, CancellationToken cancellationToken)
    {
        var vaults = _registry.All
            .OrderBy(v => v.VaultName, StringComparer.Ordinal)
            .Select(v => new VaultInfoDto(v.VaultName, v.VaultPath, _indexes.CountNotes(v.VaultName), v.VaultName == _registry.Active))
            .ToList();

        return ValueTask.FromResult(vaults);
    }
}
