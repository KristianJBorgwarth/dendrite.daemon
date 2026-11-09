using Dapper;
using Dendrite.Application.Abstractions;
using Dendrite.Domain.Models;

[module: DapperAot]

namespace Dendrite.Infrastructure.Persistence.Repositories;

internal sealed record NoteRow(string Id, string Path, string? Title, string Slug, string CreatedAt, string UpdatedAt)
{
    public Note ToNote() => new(Id, Path, Title, Slug, CreatedAt, UpdatedAt);
}

internal sealed class NoteRepository(Database db) : INoteRepository
{
    private const string Columns =
        "n.id AS Id, n.path AS Path, n.title AS Title, n.slug AS Slug, n.created_at AS CreatedAt, n.updated_at AS UpdatedAt";

    public void Insert(Note note) =>
        db.Connection.Execute(
            """
            INSERT INTO note (id, title, path, slug, created_at, updated_at)
            VALUES (@Id, @Title, @Path, @Slug, datetime('now'), datetime('now'))
            ON CONFLICT DO NOTHING
            """,
            new { note.Id, note.Title, note.Path, note.Slug },
            db.Transaction);

    public void InsertRange(IEnumerable<Note> notes)
    {
        foreach (var note in notes)
            Insert(note);
    }

    public void Update(string noteId, string path, string title, string slug) =>
        db.Connection.Execute(
            "UPDATE note SET path = @path, title = @title, slug = @slug, updated_at = datetime('now') WHERE id = @noteId",
            new { path, title, slug, noteId },
            db.Transaction);

    public void Delete(string noteId) =>
        db.Connection.Execute("DELETE FROM note WHERE id = @noteId", new { noteId }, db.Transaction);

    public Note? GetBySlug(string slug) =>
        db.Connection
            .QueryFirstOrDefault<NoteRow>($"SELECT {Columns} FROM note n WHERE n.slug = @slug", new { slug }, db.Transaction)
            ?.ToNote();

    public Note? GetByPath(string path) =>
        db.Connection
            .QueryFirstOrDefault<NoteRow>($"SELECT {Columns} FROM note n WHERE n.path = @path", new { path }, db.Transaction)
            ?.ToNote();

    public int GetNoteCount() =>
        db.Connection.QueryFirst<int>("SELECT COUNT(*) FROM note", null, db.Transaction);

    public List<Note> GetByTag(string tag) =>
        db.Connection
            .Query<NoteRow>(
                $"""
                 SELECT DISTINCT {Columns}
                 FROM note n
                 JOIN note_tag nt ON n.id = nt.note_id
                 JOIN tag t ON t.name = nt.tag_id
                 WHERE t.name LIKE @pattern
                 """,
                new { pattern = $"%{tag}%" },
                db.Transaction)
            .Select(r => r.ToNote())
            .ToList();

    public List<Note> Search(string query) =>
        db.Connection
            .Query<NoteRow>(
                $"""
                 SELECT {Columns}
                 FROM note n
                 WHERE n.title LIKE @pattern OR n.path LIKE @pattern OR n.slug LIKE @pattern
                 ORDER BY n.created_at DESC
                 """,
                new { pattern = $"%{query}%" },
                db.Transaction)
            .Select(r => r.ToNote())
            .ToList();
}
