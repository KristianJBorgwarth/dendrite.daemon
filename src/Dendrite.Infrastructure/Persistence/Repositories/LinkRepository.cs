using Dapper;
using Dendrite.Application;
using Dendrite.Application.Abstractions;
using Dendrite.Domain.Models;

namespace Dendrite.Infrastructure.Persistence.Repositories;

internal sealed record LinkRow(string Id, string FromNoteId, string TargetSlug, string Raw, string? Display, int Line, int Col);

internal sealed class LinkRepository(Database db) : ILinkRepository
{
    public void Insert(IEnumerable<Link> links)
    {
        foreach (var link in links)
        {
            db.Connection.Execute(
                """
                INSERT OR IGNORE INTO link (id, from_note_id, target_slug, raw, display, line, col)
                VALUES (@Id, @FromNoteId, @TargetSlug, @Raw, @Display, @Line, @Col)
                """,
                new { link.Id, link.FromNoteId, link.TargetSlug, link.Raw, link.Display, link.Line, link.Col },
                db.Transaction);
        }
    }

    public void Delete(string fromNoteId) =>
        db.Connection.Execute("DELETE FROM link WHERE from_note_id = @fromNoteId", new { fromNoteId }, db.Transaction);

    public List<BacklinkDto> GetBacklinks(string slug) =>
        db.Connection
            .Query<BacklinkDto>(
                """
                SELECT n.id AS Id, n.title AS Title, n.slug AS Slug, n.path AS Path,
                       l.raw AS Raw, l.line AS Line, l.col AS Col
                FROM link l
                JOIN note n ON l.from_note_id = n.id
                WHERE l.target_slug = @slug
                """,
                new { slug },
                db.Transaction)
            .AsList();

    public List<Link> GetBrokenLinks(string fromNoteId) =>
        db.Connection
            .Query<LinkRow>(
                """
                SELECT l.id AS Id, l.from_note_id AS FromNoteId, l.target_slug AS TargetSlug, l.raw AS Raw,
                       l.display AS Display, l.line AS Line, l.col AS Col
                FROM link l
                LEFT JOIN note n ON l.target_slug = n.slug
                WHERE l.from_note_id = @fromNoteId AND n.id IS NULL
                """,
                new { fromNoteId },
                db.Transaction)
            .Select(r => new Link(r.Id, r.FromNoteId, r.TargetSlug, r.Raw, r.Display ?? "", r.Line, r.Col))
            .ToList();
}
