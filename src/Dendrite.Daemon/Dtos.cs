namespace Dendrite.Daemon.Dtos;

public sealed record NoteDto(string Id, string Path, string? Title, string Slug);

public sealed record BacklinkDto(string Id, string? Title, string Slug, string Path, string Raw, int Line, int Col);

public sealed record LinkDiagnosticDto(
    string NoteId,
    string NotePath,
    int Line,
    int Col,
    string Raw,
    string Target,
    string Message);

public sealed record GotoNoteResult(string Target, string Type);
