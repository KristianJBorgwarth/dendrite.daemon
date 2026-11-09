namespace Dendrite.Application.Vaults;

/// <summary>A vault request about a file, so the vault can be inferred from the path.</summary>
public interface IVaultPathRequest : IVaultRequest
{
    string Path { get; }
}
