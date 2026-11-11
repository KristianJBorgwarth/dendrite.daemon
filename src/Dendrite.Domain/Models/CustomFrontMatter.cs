namespace Dendrite.Domain.Models;

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
