using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record GetBacklinks : IQuery<List<BacklinkDto>>
{
    public string Path { get; }

    public GetBacklinks(string path)
    {
        Path = path;
    }
}
