using System.Text.Json;
using System.Text.Json.Serialization;
using Dendrite.Daemon.Dtos;
using Dendrite.Daemon.Features.Completion;
using Dendrite.Daemon.Features.Diagnostics;
using Dendrite.Daemon.Features.Note;
using Dendrite.Daemon.Features.Vault;

namespace Dendrite.Daemon.Rpc;

public sealed class RpcRequest
{
    [JsonPropertyName("jsonrpc")] public string? JsonRpc { get; init; }
    [JsonPropertyName("id")] public JsonElement? Id { get; init; }
    [JsonPropertyName("method")] public string? Method { get; init; }
    [JsonPropertyName("params")] public JsonElement? Params { get; init; }
}

// Source-generated so the daemon stays reflection-free under Native AOT.
[JsonSourceGenerationOptions(
    PropertyNamingPolicy = JsonKnownNamingPolicy.CamelCase,
    PropertyNameCaseInsensitive = true,
    DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull)]
[JsonSerializable(typeof(RpcRequest))]
[JsonSerializable(typeof(string))]
[JsonSerializable(typeof(InitializeVault))]
[JsonSerializable(typeof(RebuildIndex))]
[JsonSerializable(typeof(CreateNote))]
[JsonSerializable(typeof(SaveNote))]
[JsonSerializable(typeof(DeleteNote))]
[JsonSerializable(typeof(GotoNote))]
[JsonSerializable(typeof(GetBacklinks))]
[JsonSerializable(typeof(GetNotesByTag))]
[JsonSerializable(typeof(GetNotesByCf))]
[JsonSerializable(typeof(CompleteTag))]
[JsonSerializable(typeof(CompleteSlug))]
[JsonSerializable(typeof(LinkDiagnostics))]
[JsonSerializable(typeof(GotoNoteResult))]
[JsonSerializable(typeof(List<string>))]
[JsonSerializable(typeof(List<NoteDto>))]
[JsonSerializable(typeof(List<BacklinkDto>))]
[JsonSerializable(typeof(List<LinkDiagnosticDto>))]
internal sealed partial class RpcJsonContext : JsonSerializerContext;
