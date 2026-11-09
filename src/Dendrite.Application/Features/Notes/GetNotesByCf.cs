using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record GetNotesByCf : IQuery<List<NoteDto>>
{
    public string Key { get; }
    public string Value { get; }

    public GetNotesByCf(string key, string value)
    {
        Key = key;
        Value = value;
    }
}
