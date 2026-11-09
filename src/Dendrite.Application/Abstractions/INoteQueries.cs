namespace Dendrite.Application.Abstractions;

/// <summary>Read side of the index. Answers questions straight from the database without loading aggregates.</summary>
public interface INoteQueries
{
    int CountNotes();

    string? FindPathBySlug(string slug);

    List<NoteDto> SearchNotes(string query);

    List<NoteDto> GetNotesByTag(string tag);

    List<NoteDto> GetNotesByCustomFrontMatter(string key, string value);

    List<string> SearchTags(string query);

    List<BacklinkDto> GetBacklinks(string notePath);

    List<LinkDiagnosticDto> GetLinkDiagnostics(string notePath);
}
