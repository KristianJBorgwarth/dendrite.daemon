namespace Dendrite.Domain.FileHandling;

public readonly record struct ParsedLink
{
    public LinkKind Kind { get; }
    public string Target { get; }

    public ParsedLink(LinkKind kind, string target)
    {
        Kind = kind;
        Target = target;
    }
}
