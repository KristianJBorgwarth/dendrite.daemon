namespace Dendrite.Application.Vaults;

/// <summary>A request that targets a vault. When <see cref="Vault"/> is null the vault is inferred.</summary>
public interface IVaultRequest
{
    string? Vault { get; }
}
