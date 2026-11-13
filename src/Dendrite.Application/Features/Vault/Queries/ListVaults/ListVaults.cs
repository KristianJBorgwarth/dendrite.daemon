using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed record ListVaults : IQuery<List<VaultInfoDto>>;
