using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record GetNotesByTag : IQuery<List<NoteDto>>
{
    public string Tag { get; }

    public GetNotesByTag(string tag)
    {
        Tag = tag;
    }
}
