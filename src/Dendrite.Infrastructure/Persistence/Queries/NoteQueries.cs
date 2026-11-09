using Dapper;
using Dendrite.Application;
using Dendrite.Application.Abstractions;

namespace Dendrite.Infrastructure.Persistence.Queries;

internal sealed class NoteQueries(Database db) : INoteQueries
{
    private const string NoteColumns = "n.id AS Id, n.path AS Path, n.title AS Title, n.slug AS Slug";

    private readonly Database _db = db;

    public int CountNotes() =>
        _db.Connection.QueryFirst<int>("SELECT COUNT(*) FROM note", null, _db.Transaction);

    public string? FindPathBySlug(string slug) =>
        _db.Connection.QueryFirstOrDefault<string>("SELECT path FROM note WHERE slug = @slug", new { slug }, _db.Transaction);

    public List<NoteDto> SearchNotes(string query) =>
        _db.Connection
            .Query<NoteDto>(
                $"""
                 SELECT {NoteColumns}
                 FROM note n
                 WHERE n.title LIKE @pattern OR n.path LIKE @pattern OR n.slug LIKE @pattern
                 ORDER BY n.created_at DESC
                 """,
                new { pattern = $"%{query}%" },
                _db.Transaction)
            .AsList();

    public List<NoteDto> GetNotesByTag(string tag) =>
        _db.Connection
            .Query<NoteDto>(
                $"""
                 SELECT DISTINCT {NoteColumns}
                 FROM note n
                 JOIN note_tag nt ON n.id = nt.note_id
                 JOIN tag t ON t.name = nt.tag_id
                 WHERE t.name LIKE @pattern
                 """,
                new { pattern = $"%{tag}%" },
                _db.Transaction)
            .AsList();

    public List<NoteDto> GetNotesByCustomFrontMatter(string key, string value) =>
        _db.Connection
            .Query<NoteDto>(
                $"""
                 SELECT {NoteColumns}
                 FROM note n
                 INNER JOIN custom_frontmatter cfe ON n.id = cfe.note_id
                 WHERE cfe.key = @key AND cfe.value LIKE @pattern
                 """,
                new { key, pattern = $"%{value}%" },
                _db.Transaction)
            .AsList();

    public List<string> SearchTags(string query) =>
        _db.Connection
            .Query<string>("SELECT name FROM tag WHERE name LIKE @pattern", new { pattern = $"%{query}%" }, _db.Transaction)
            .AsList();

    public List<BacklinkDto> GetBacklinks(string notePath) =>
        _db.Connection
            .Query<BacklinkDto>(
                """
                SELECT n.id AS Id, n.title AS Title, n.slug AS Slug, n.path AS Path,
                       l.raw AS Raw, l.line AS Line, l.col AS Col
                FROM link l
                JOIN note n ON l.from_note_id = n.id
                WHERE l.target_slug = (SELECT slug FROM note WHERE path = @notePath)
                """,
                new { notePath },
                _db.Transaction)
            .AsList();

    public List<LinkDiagnosticDto> GetLinkDiagnostics(string notePath) =>
        _db.Connection
            .Query<LinkDiagnosticDto>(
                """
                SELECT l.id AS NoteId, src.path AS NotePath, l.line AS Line, l.col AS Col, l.raw AS Raw,
                       l.target_slug AS Target, 'Broken link: target note does not exist' AS Message
                FROM link l
                JOIN note src ON l.from_note_id = src.id
                LEFT JOIN note target ON l.target_slug = target.slug
                WHERE src.path = @notePath AND target.id IS NULL
                """,
                new { notePath },
                _db.Transaction)
            .AsList();
}
