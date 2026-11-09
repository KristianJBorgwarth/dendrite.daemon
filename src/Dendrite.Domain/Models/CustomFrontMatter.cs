namespace Dendrite.Domain.Models;

/// <summary>A custom frontmatter key/value pair on a note. Owned by the <see cref="Note"/> aggregate.</summary>
public sealed record CustomFrontMatter
{
    public string Key { get; }
    public string Value { get; }

    public CustomFrontMatter(string key, string value)
    {
        Key = key;
        Value = value;
    }
}
