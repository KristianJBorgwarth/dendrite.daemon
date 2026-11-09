using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record SaveNote : ICommand<Unit>
{
    public string Path { get; }

    public SaveNote(string path)
    {
        Path = path;
    }
}
