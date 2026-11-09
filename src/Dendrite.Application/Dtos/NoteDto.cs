namespace Dendrite.Application;

public sealed record NoteDto
{
    public string Id { get; }
    public string Path { get; }
    public string? Title { get; }
    public string Slug { get; }

    public NoteDto(string id, string path, string? title, string slug)
    {
        Id = id;
        Path = path;
        Title = title;
        Slug = slug;
    }
}
