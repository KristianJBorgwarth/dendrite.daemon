using YamlDotNet.RepresentationModel;

namespace Dendrite.Domain.FileHandling;

public sealed record FrontMatter
{
    public static readonly FrontMatter Empty = new("", [], []);

    private static readonly HashSet<string> KnownKeys = ["title", "tags", "created", "updated", "date"];

    public string Title { get; }
    public IReadOnlyList<string> Tags { get; }
    public IReadOnlyList<KeyValuePair<string, string>> Custom { get; }

    public FrontMatter(string title, IReadOnlyList<string> tags, IReadOnlyList<KeyValuePair<string, string>> custom)
    {
        Title = title;
        Tags = tags;
        Custom = custom;
    }

    public static FrontMatter? Parse(string file)
    {
        var content = file.Trim();
        if (!content.StartsWith("---", StringComparison.Ordinal))
            return null;

        var yaml = ExtractYaml(content[3..]);

        var stream = new YamlStream();
        stream.Load(new StringReader(yaml));
        if (stream.Documents.Count == 0 || stream.Documents[0].RootNode is not YamlMappingNode root)
            return Empty;

        var title = "";
        var tags = new List<string>();
        var custom = new List<KeyValuePair<string, string>>();

        foreach (var (keyNode, valueNode) in root.Children)
        {
            if (keyNode is not YamlScalarNode keyScalar || keyScalar.Value is not { } key)
                continue;

            switch (valueNode)
            {
                case YamlScalarNode scalar when key == "title":
                    title = scalar.Value ?? "";
                    break;
                case YamlScalarNode scalar when key == "tags":
                    if (!string.IsNullOrWhiteSpace(scalar.Value))
                        tags.Add(scalar.Value);
                    break;
                case YamlSequenceNode seq when key == "tags":
                    tags.AddRange(Scalars(seq));
                    break;
                case YamlScalarNode scalar when !KnownKeys.Contains(key):
                    custom.Add(new(key, scalar.Value ?? ""));
                    break;
                case YamlSequenceNode seq when !KnownKeys.Contains(key):
                    custom.AddRange(Scalars(seq).Select(v => new KeyValuePair<string, string>(key, v)));
                    break;
            }
        }

        return new FrontMatter(title, [.. tags.Distinct()], custom);
    }

    private static IEnumerable<string> Scalars(YamlSequenceNode seq) =>
        seq.Children.OfType<YamlScalarNode>().Select(s => s.Value).OfType<string>();

    private static string ExtractYaml(string afterOpening)
    {
        var lines = afterOpening.Trim().Split('\n');
        var yaml = new List<string>();
        foreach (var line in lines)
        {
            if (line.TrimEnd() == "---")
                break;
            yaml.Add(line);
        }
        return string.Join('\n', yaml);
    }
}
