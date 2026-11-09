using System.Text.Json.Serialization;

namespace Dendrite.Application.Features.Vault;

public sealed class DailyNotesConfig
{
    [JsonPropertyName("dir")] public string Dir { get; init; } = "";
    [JsonPropertyName("filename_format")] public string FilenameFormat { get; init; } = "";
    [JsonPropertyName("template_name")] public string TemplateName { get; init; } = "";
}
