using System.Text.Json;
using Dendrite.Application;
using Dendrite.Application.Abstractions;
using Dendrite.Application.Vaults;
using Dendrite.Infrastructure.Persistence;

namespace Dendrite.Infrastructure.FileSystem;

/// <summary>Stores the registered vaults and the active one as JSON next to the indexes.</summary>
public sealed class VaultCatalog : IVaultCatalog
{
    private static string FilePath => Path.Combine(Database.DataRoot(), "vaults.json");

    public VaultCatalogState Load()
    {
        if (!File.Exists(FilePath))
            return new VaultCatalogState();

        try
        {
            return JsonSerializer.Deserialize(File.ReadAllText(FilePath), CatalogJsonContext.Default.VaultCatalogState)
                ?? new VaultCatalogState();
        }
        catch (Exception ex) when (ex is JsonException or IOException)
        {
            Log.Warn($"could not read {FilePath}: {ex.Message}");
            return new VaultCatalogState();
        }
    }

    public void Save(VaultCatalogState state)
    {
        Directory.CreateDirectory(Path.GetDirectoryName(FilePath)!);

        var temp = FilePath + ".tmp";
        File.WriteAllText(temp, JsonSerializer.Serialize(state, CatalogJsonContext.Default.VaultCatalogState));
        File.Move(temp, FilePath, overwrite: true);
    }
}
