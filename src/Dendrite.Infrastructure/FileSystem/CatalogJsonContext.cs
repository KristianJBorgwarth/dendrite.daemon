using System.Text.Json.Serialization;
using Dendrite.Application.Vaults;

namespace Dendrite.Infrastructure.FileSystem;

[JsonSourceGenerationOptions(WriteIndented = true, PropertyNameCaseInsensitive = true)]
[JsonSerializable(typeof(VaultCatalogState))]
internal sealed partial class CatalogJsonContext : JsonSerializerContext;
