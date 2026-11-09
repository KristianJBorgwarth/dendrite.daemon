using Dendrite.Application.Abstractions;
using Dendrite.Application.Services;
using Dendrite.Application.Vaults;
using Dendrite.Domain.Vault;
using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed class InitializeVaultHandler : ICommandHandler<InitializeVault, Unit>
{
    private readonly VaultRegistry _registry;
    private readonly VaultContext _vault;
    private readonly IIndexDatabase _indexDatabase;
    private readonly INoteQueries _queries;
    private readonly IndexRebuilder _rebuilder;

    public InitializeVaultHandler(
        VaultRegistry registry,
        VaultContext vault,
        IIndexDatabase indexDatabase,
        INoteQueries queries,
        IndexRebuilder rebuilder)
    {
        _registry = registry;
        _vault = vault;
        _indexDatabase = indexDatabase;
        _queries = queries;
        _rebuilder = rebuilder;
    }

    public ValueTask<Unit> Handle(InitializeVault command, CancellationToken cancellationToken)
    {
        var c = command.Config;
        if (string.IsNullOrEmpty(c.VaultPath))
            throw new InvalidOperationException("vault_path is empty — check your dendrite.nvim config");

        var daily = c.DailyNotes ?? new DailyNotesConfig();
        var vaultPath = c.VaultPath.Length > 1 ? c.VaultPath.TrimEnd('/', '\\') : c.VaultPath;

        var config = new VaultConfig(
            c.VaultName,
            vaultPath,
            c.TemplatesDir,
            c.ExcludeIndexFiles ?? [],
            c.OverrideDefaultIgnores,
            new DailyNotes(daily.Dir, daily.FilenameFormat, daily.TemplateName));

        var added = _registry.Register(config);
        try
        {
            _vault.Select(config);
            _indexDatabase.Open(c.VaultName);
        }
        catch
        {
            if (added)
                _registry.Remove(c.VaultName);
            throw;
        }

        if (_queries.CountNotes() > 0)
        {
            Log.Debug("Vault already initialized, skipping index rebuild");
            return ValueTask.FromResult(Unit.Value);
        }

        _rebuilder.Rebuild();
        return ValueTask.FromResult(Unit.Value);
    }
}
