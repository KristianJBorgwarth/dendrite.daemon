namespace Dendrite.Application;

public sealed record BacklinkDto
{
    public string Id { get; }
    public string? Title { get; }
    public string Slug { get; }
    public string Path { get; }
    public string Raw { get; }
    public int Line { get; }
    public int Col { get; }

    public BacklinkDto(string id, string? title, string slug, string path, string raw, int line, int col)
    {
        Id = id;
        Title = title;
        Slug = slug;
        Path = path;
        Raw = raw;
        Line = line;
        Col = col;
    }
}
