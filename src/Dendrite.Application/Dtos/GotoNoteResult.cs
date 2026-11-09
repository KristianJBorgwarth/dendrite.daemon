namespace Dendrite.Application;

public sealed record GotoNoteResult
{
    public string Target { get; }
    public string Type { get; }

    public GotoNoteResult(string target, string type)
    {
        Target = target;
        Type = type;
    }
}
