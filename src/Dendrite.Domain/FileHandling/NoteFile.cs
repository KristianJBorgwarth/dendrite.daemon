using System.Text;
using System.Text.RegularExpressions;

namespace Dendrite.Domain.FileHandling;

public sealed record ExtractedLink(string Raw, string Display, string TargetSlug, int Line, int Col);

public sealed record NoteFile(
    string Path,
    string Title,
    string Slug,
    FrontMatter FrontMatter,
    IReadOnlyList<ExtractedLink> Links);

public static partial class NoteFileParser
{
    [GeneratedRegex(@"\[\[([^\]|]+)(?:\|([^\]]+))?\]\]")]
    private static partial Regex LinkRegex();

    public static NoteFile Parse(string vaultRoot, string path, string content)
    {
        var body = content.Trim();

        var frontMatter = FrontMatter.Parse(body);
        var title = frontMatter?.Title ?? Path.GetFileNameWithoutExtension(path);

        return new NoteFile(
            path,
            title,
            Slug.Relative(vaultRoot, path),
            frontMatter ?? FrontMatter.Empty,
            ExtractLinks(body));
    }

    public static List<ExtractedLink> ExtractLinks(string body)
    {
        var links = new List<ExtractedLink>();
        var lineNumber = 1;

        foreach (var line in body.Split('\n'))
        {
            foreach (Match match in LinkRegex().Matches(line))
            {
                var display = match.Groups[2].Success ? match.Groups[2].Value : "";
                links.Add(new ExtractedLink(
                    match.Value,
                    display,
                    match.Groups[1].Value,
                    lineNumber,
                    Encoding.UTF8.GetByteCount(line.AsSpan(0, match.Index)) + 1));
            }
            lineNumber++;
        }

        return links;
    }
}
