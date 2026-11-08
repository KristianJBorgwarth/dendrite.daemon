using Dendrite.Daemon.Dtos;
using Dendrite.Daemon.Persistence.Repositories;
using Mediator;

namespace Dendrite.Daemon.Features.Note;

public sealed record GetNotesByTag(string Tag) : IQuery<List<NoteDto>>;

public sealed class GetNotesByTagHandler(NoteRepository notes) : IQueryHandler<GetNotesByTag, List<NoteDto>>
{
    public ValueTask<List<NoteDto>> Handle(GetNotesByTag query, CancellationToken cancellationToken) =>
        ValueTask.FromResult(notes.GetByTag(query.Tag).Select(n => new NoteDto(n.Id, n.Path, n.Title, n.Slug)).ToList());
}

public sealed record GetNotesByCf(string Key, string Value) : IQuery<List<NoteDto>>;

public sealed class GetNotesByCfHandler(CfeRepository cfe) : IQueryHandler<GetNotesByCf, List<NoteDto>>
{
    public ValueTask<List<NoteDto>> Handle(GetNotesByCf query, CancellationToken cancellationToken) =>
        ValueTask.FromResult(cfe.Search(query.Key, query.Value).Select(n => new NoteDto(n.Id, n.Path, n.Title, n.Slug)).ToList());
}
