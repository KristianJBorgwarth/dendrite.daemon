using Dapper;
using Dendrite.Domain.Models;
using Dendrite.Infrastructure.Persistence.Rows;

namespace Dendrite.Infrastructure.Persistence.Repositories;

internal sealed class NoteRepository(Database db) : INoteRepository
{
    private const string Columns = "id AS Id, path AS Path, title AS Title, slug AS Slug, created_at AS CreatedAt, updated_at AS UpdatedAt";

    private readonly Database _db = db;

    public Note? GetBySlug(string slug) =>
        Load(_db.Connection.QueryFirstOrDefault<NoteRow>($"SELECT {Columns} FROM note WHERE slug = @slug", new { slug }, _db.Transaction));

    public Note? GetByPath(string path) =>
        Load(_db.Connection.QueryFirstOrDefault<NoteRow>($"SELECT {Columns} FROM note WHERE path = @path", new { path }, _db.Transaction));

    public void Add(Note note)
    {
        _db.Connection.Execute(
            """
            INSERT INTO note (id, title, path, slug, created_at, updated_at)
            VALUES (@Id, @Title, @Path, @Slug, @CreatedAt, @UpdatedAt)
            """,
            new { note.Id, note.Title, note.Path, note.Slug, note.CreatedAt, note.UpdatedAt },
            _db.Transaction);

        InsertChildren(note);
    }

    public void Save(Note note)
    {
        _db.Connection.Execute(
            "UPDATE note SET path = @Path, title = @Title, slug = @Slug, updated_at = @UpdatedAt WHERE id = @Id",
            new { note.Path, note.Title, note.Slug, note.UpdatedAt, note.Id },
            _db.Transaction);

        DeleteChildren(note.Id);
        InsertChildren(note);
    }

    public void Remove(Note note)
    {
        DeleteChildren(note.Id);
        _db.Connection.Execute("DELETE FROM note WHERE id = @Id", new { note.Id }, _db.Transaction);
    }

    private Note? Load(NoteRow? row)
    {
        if (row is null)
            return null;

        var noteId = row.Id;

        var tags = _db.Connection
            .Query<string>("SELECT tag_id FROM note_tag WHERE note_id = @noteId", new { noteId }, _db.Transaction)
            .AsList();

        var links = _db.Connection
            .Query<LinkRow>(
                """
                SELECT id AS Id, target_slug AS TargetSlug, raw AS Raw, display AS Display, line AS Line, col AS Col
                FROM link
                WHERE from_note_id = @noteId
                ORDER BY line, col
                """,
                new { noteId },
                _db.Transaction)
            .Select(l => new Link(l.Id, l.TargetSlug, l.Raw, l.Display ?? "", l.Line, l.Col))
            .ToList();

        var custom = _db.Connection
            .Query<CustomFrontMatterRow>(
                "SELECT key AS Key, value AS Value FROM custom_frontmatter WHERE note_id = @noteId",
                new { noteId },
                _db.Transaction)
            .Select(c => new CustomFrontMatter(c.Key, c.Value))
            .ToList();

        return Note.Rehydrate(row.Id, row.Path, row.Title, row.Slug, row.CreatedAt, row.UpdatedAt, tags, links, custom);
    }

    private void InsertChildren(Note note)
    {
        foreach (var tag in note.Tags)
        {
            _db.Connection.Execute("INSERT OR IGNORE INTO tag (name) VALUES (@tag)", new { tag }, _db.Transaction);
            _db.Connection.Execute(
                "INSERT OR IGNORE INTO note_tag (note_id, tag_id) VALUES (@NoteId, @tag)",
                new { NoteId = note.Id, tag },
                _db.Transaction);
        }

        foreach (var link in note.Links)
        {
            _db.Connection.Execute(
                """
                INSERT OR IGNORE INTO link (id, from_note_id, target_slug, raw, display, line, col)
                VALUES (@Id, @NoteId, @TargetSlug, @Raw, @Display, @Line, @Col)
                """,
                new { link.Id, NoteId = note.Id, link.TargetSlug, link.Raw, link.Display, link.Line, link.Col },
                _db.Transaction);
        }

        foreach (var entry in note.CustomFrontMatter)
        {
            _db.Connection.Execute(
                "INSERT OR IGNORE INTO custom_frontmatter (note_id, key, value) VALUES (@NoteId, @Key, @Value)",
                new { NoteId = note.Id, entry.Key, entry.Value },
                _db.Transaction);
        }
    }

    private void DeleteChildren(string noteId)
    {
        _db.Connection.Execute("DELETE FROM link WHERE from_note_id = @noteId", new { noteId }, _db.Transaction);
        _db.Connection.Execute("DELETE FROM note_tag WHERE note_id = @noteId", new { noteId }, _db.Transaction);
        _db.Connection.Execute("DELETE FROM custom_frontmatter WHERE note_id = @noteId", new { noteId }, _db.Transaction);
    }
}
