using System.Text.Json.Serialization;
using Dendrite.Application.Features.Vault;

namespace Dendrite.Infrastructure.FileSystem;

[JsonSourceGenerationOptions(WriteIndented = true, PropertyNameCaseInsensitive = true)]
[JsonSerializable(typeof(VaultConfigDto))]
internal sealed partial class ManifestJsonContext : JsonSerializerContext;
