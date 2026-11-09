namespace Dendrite.Application;

public sealed record LinkDiagnosticDto
{
    public string NoteId { get; }
    public string NotePath { get; }
    public int Line { get; }
    public int Col { get; }
    public string Raw { get; }
    public string Target { get; }
    public string Message { get; }

    public LinkDiagnosticDto(string noteId, string notePath, int line, int col, string raw, string target, string message)
    {
        NoteId = noteId;
        NotePath = notePath;
        Line = line;
        Col = col;
        Raw = raw;
        Target = target;
        Message = message;
    }
}
