using Dendrite.Application.Abstractions;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record GetBacklinks(string Path) : IQuery<List<BacklinkDto>>;

public sealed class GetBacklinksHandler(INoteRepository notes, ILinkRepository links)
    : IQueryHandler<GetBacklinks, List<BacklinkDto>>
{
    public ValueTask<List<BacklinkDto>> Handle(GetBacklinks query, CancellationToken cancellationToken)
    {
        var note = notes.GetByPath(query.Path);
        return ValueTask.FromResult(note is null ? [] : links.GetBacklinks(note.Slug));
    }
}
