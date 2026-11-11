using Dendrite.Application;
using Dendrite.Application.Abstractions;
using Dendrite.Application.Services;
using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Daemon;

/// <summary>Selects the vault a request targets, then sends it through Mediator.</summary>
public sealed class VaultDispatcher
{
    private readonly ISender _sender;
    private readonly VaultRegistry _registry;
    private readonly VaultContext _context;
    private readonly IVaultManifest _manifest;
    private readonly VaultOpener _opener;

    public VaultDispatcher(
        ISender sender,
        VaultRegistry registry,
        VaultContext context,
        IVaultManifest manifest,
        VaultOpener opener)
    {
        _sender = sender;
        _registry = registry;
        _context = context;
        _manifest = manifest;
        _opener = opener;
    }

    public ValueTask<TResponse> Send<TResponse>(ICommand<TResponse> command, CancellationToken ct = default)
    {
        Select(command);
        return _sender.Send(command, ct);
    }

    public ValueTask<TResponse> Send<TResponse>(IQuery<TResponse> query, CancellationToken ct = default)
    {
        Select(query);
        return _sender.Send(query, ct);
    }

    private void Select(object request)
    {
        if (request is not IVaultRequest vaultRequest)
            return;

        var path = (request as IVaultPathRequest)?.Path;
        if (string.IsNullOrEmpty(vaultRequest.Vault))
            AdoptVaultAround(path);

        _context.Select(_registry.Resolve(vaultRequest.Vault, path));
    }

    // A file inside a vault the daemon does not know at that location (moved or new) brings its own manifest.
    private void AdoptVaultAround(string? path)
    {
        if (string.IsNullOrEmpty(path) || _registry.Owns(path))
            return;

        try
        {
            if (_manifest.FindRoot(path) is { } root)
                _opener.Open(root, activate: false);
        }
        catch (Exception ex)
        {
            Log.Warn($"could not open the vault around '{path}': {ex.Message}");
        }
    }
}
