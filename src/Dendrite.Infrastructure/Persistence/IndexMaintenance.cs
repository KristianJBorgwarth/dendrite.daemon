using Dapper;
using Dendrite.Application.Abstractions;

namespace Dendrite.Infrastructure.Persistence;

internal sealed class IndexMaintenance : IIndexMaintenance
{
    private readonly Database _db;

    public IndexMaintenance(Database db)
    {
        _db = db;
    }

    public void Wipe()
    {
        _db.Connection.Execute("DELETE FROM note", null, _db.Transaction);
        _db.Connection.Execute("DELETE FROM tag", null, _db.Transaction);
        _db.Connection.Execute("DELETE FROM note_tag", null, _db.Transaction);
        _db.Connection.Execute("DELETE FROM link", null, _db.Transaction);
        _db.Connection.Execute("DELETE FROM custom_frontmatter", null, _db.Transaction);
    }
}
