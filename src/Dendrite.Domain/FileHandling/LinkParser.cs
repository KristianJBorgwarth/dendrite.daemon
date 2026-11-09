namespace Dendrite.Domain.FileHandling;

public enum LinkKind
{
    Unknown,
    Note,
    Url,
}

public readonly record struct ParsedLink(LinkKind Kind, string Target);

public static class LinkParser
{
    public static ParsedLink Parse(string s)
    {
        if (s.Length >= 5 && s.StartsWith("[[", StringComparison.Ordinal) && s.EndsWith("]]", StringComparison.Ordinal))
        {
            var content = s[2..^2];
            var bar = content.IndexOf('|');
            return new ParsedLink(LinkKind.Note, bar >= 0 ? content[..bar] : content);
        }

        if (s.StartsWith('['))
        {
            var close = s.IndexOf("](", StringComparison.Ordinal);
            if (close >= 0 && s.EndsWith(')'))
                return new ParsedLink(LinkKind.Url, s[(close + 2)..^1]);
        }

        return new ParsedLink(LinkKind.Note, s);
    }
}
