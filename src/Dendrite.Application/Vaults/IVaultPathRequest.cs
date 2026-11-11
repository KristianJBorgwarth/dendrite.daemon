namespace Dendrite.Application.Vaults;

public interface IVaultPathRequest : IVaultRequest
{
    string Path { get; }
}
