using Dendrite.Domain.FileHandling;

namespace Dendrite.Application.Services;

public sealed record Template
{
    public string Content { get; }
    public string Title { get; }
    public string Slug { get; }
    public string Filename { get; }
    public FrontMatter FrontMatter { get; }

    public Template(string content, string title, string slug, string filename, FrontMatter frontMatter)
    {
        Content = content;
        Title = title;
        Slug = slug;
        Filename = filename;
        FrontMatter = frontMatter;
    }
}
