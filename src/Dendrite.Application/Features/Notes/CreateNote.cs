using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record CreateNote : ICommand<string>
{
    public string Title { get; }
    public string TemplateName { get; }
    public string Directory { get; }

    public CreateNote(string title, string templateName, string directory)
    {
        Title = title;
        TemplateName = templateName;
        Directory = directory;
    }
}
