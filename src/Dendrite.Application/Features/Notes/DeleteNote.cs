using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record DeleteNote : ICommand<Unit>
{
    public string Path { get; }

    public DeleteNote(string path)
    {
        Path = path;
    }
}
