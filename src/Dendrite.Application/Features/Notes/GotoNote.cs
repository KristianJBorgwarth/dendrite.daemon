using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record GotoNote : IQuery<GotoNoteResult?>
{
    public string Link { get; }

    public GotoNote(string link)
    {
        Link = link;
    }
}
