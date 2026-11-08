using Dapper;
using Dendrite.Daemon.Models;

namespace Dendrite.Daemon.Persistence.Repositories;

public sealed class CfeRepository(Database db)
{
    public void InsertRange(IEnumerable<CustomFrontMatter> entries)
    {
        foreach (var cfe in entries)
        {
            db.Connection.Execute(
                "INSERT INTO custom_frontmatter (note_id, key, value) VALUES (@NoteId, @Key, @Value) ON CONFLICT DO NOTHING",
                new { cfe.NoteId, cfe.Key, cfe.Value },
                db.Transaction);
        }
    }

    public void Delete(string noteId) =>
        db.Connection.Execute("DELETE FROM custom_frontmatter WHERE note_id = @noteId", new { noteId }, db.Transaction);

    public List<Note> Search(string key, string value) =>
        db.Connection
            .Query<NoteRow>(
                """
                SELECT n.id AS Id, n.path AS Path, n.title AS Title, n.slug AS Slug,
                       n.created_at AS CreatedAt, n.updated_at AS UpdatedAt
                FROM note n
                INNER JOIN custom_frontmatter cfe ON n.id = cfe.note_id
                WHERE cfe.key = @key AND cfe.value LIKE @pattern
                """,
                new { key, pattern = $"%{value}%" },
                db.Transaction)
            .Select(r => r.ToNote())
            .ToList();
}

public sealed class IndexRepository(Database db)
{
    public void Wipe()
    {
        db.Connection.Execute("DELETE FROM note", null, db.Transaction);
        db.Connection.Execute("DELETE FROM tag", null, db.Transaction);
        db.Connection.Execute("DELETE FROM note_tag", null, db.Transaction);
        db.Connection.Execute("DELETE FROM link", null, db.Transaction);
        db.Connection.Execute("DELETE FROM custom_frontmatter", null, db.Transaction);
    }
}
