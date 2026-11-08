using Dendrite.Daemon.Dtos;
using Dendrite.Daemon.Persistence.Repositories;
using Mediator;

namespace Dendrite.Daemon.Features.Note;

public sealed record GetBacklinks(string Path) : IQuery<List<BacklinkDto>>;

public sealed class GetBacklinksHandler(NoteRepository notes, LinkRepository links)
    : IQueryHandler<GetBacklinks, List<BacklinkDto>>
{
    public ValueTask<List<BacklinkDto>> Handle(GetBacklinks query, CancellationToken cancellationToken)
    {
        var note = notes.GetByPath(query.Path);
        return ValueTask.FromResult(note is null ? [] : links.GetBacklinks(note.Slug));
    }
}
