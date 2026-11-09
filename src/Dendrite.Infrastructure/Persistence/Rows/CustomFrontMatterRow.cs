namespace Dendrite.Infrastructure.Persistence.Rows;

internal sealed record CustomFrontMatterRow
{
    public string Key { get; }
    public string Value { get; }

    public CustomFrontMatterRow(string key, string value)
    {
        Key = key;
        Value = value;
    }
}
