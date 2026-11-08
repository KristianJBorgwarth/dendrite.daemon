using Dapper;
using Dendrite.Daemon.Models;

namespace Dendrite.Daemon.Persistence.Repositories;

public sealed class TagRepository(Database db)
{
    public void Insert(IEnumerable<Tag> tags)
    {
        foreach (var tag in tags)
            db.Connection.Execute("INSERT OR IGNORE INTO tag (name) VALUES (@Name)", new { tag.Name }, db.Transaction);
    }

    public void InsertNoteTags(IEnumerable<NoteTag> noteTags)
    {
        foreach (var nt in noteTags)
        {
            db.Connection.Execute(
                "INSERT OR IGNORE INTO note_tag (note_id, tag_id) VALUES (@NoteId, @TagId)",
                new { nt.NoteId, nt.TagId },
                db.Transaction);
        }
    }

    public List<Tag> GetByNames(IEnumerable<string> names)
    {
        var tags = new List<Tag>();
        foreach (var name in names)
        {
            tags.AddRange(
                db.Connection
                    .Query<string>("SELECT name FROM tag WHERE name = @name", new { name }, db.Transaction)
                    .Select(n => new Tag(n)));
        }
        return tags;
    }

    public List<Tag> Search(string query) =>
        db.Connection
            .Query<string>("SELECT name FROM tag WHERE name LIKE @pattern", new { pattern = $"%{query}%" }, db.Transaction)
            .Select(n => new Tag(n))
            .ToList();

    public void DeleteNoteTags(string noteId) =>
        db.Connection.Execute("DELETE FROM note_tag WHERE note_id = @noteId", new { noteId }, db.Transaction);
}
