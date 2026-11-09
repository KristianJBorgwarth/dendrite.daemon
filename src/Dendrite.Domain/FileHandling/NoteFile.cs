namespace Dendrite.Domain.FileHandling;

public sealed record NoteFile
{
    public string Path { get; }
    public string Title { get; }
    public string Slug { get; }
    public FrontMatter FrontMatter { get; }
    public IReadOnlyList<ExtractedLink> Links { get; }

    public NoteFile(string path, string title, string slug, FrontMatter frontMatter, IReadOnlyList<ExtractedLink> links)
    {
        Path = path;
        Title = title;
        Slug = slug;
        FrontMatter = frontMatter;
        Links = links;
    }
}
