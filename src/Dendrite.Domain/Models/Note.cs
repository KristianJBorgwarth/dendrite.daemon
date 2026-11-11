using Dendrite.Domain.FileHandling;

namespace Dendrite.Domain.Models;

public sealed class Note
{

    public string Id { get; }
    public string Path { get; private set; }
    public string? Title { get; private set; }
    public string Slug { get; private set; }
    public string CreatedAt { get; }
    public string UpdatedAt { get; private set; }

    private readonly List<string> _tags;
    private readonly List<Link> _links;
    private readonly List<CustomFrontMatter> _customFrontMatter;

    public IReadOnlyList<string> Tags => _tags;
    public IReadOnlyList<Link> Links => _links;
    public IReadOnlyList<CustomFrontMatter> CustomFrontMatter => _customFrontMatter;

    private Note(
        string id,
        string path,
        string? title,
        string slug,
        string createdAt,
        string updatedAt,
        List<string> tags,
        List<Link> links,
        List<CustomFrontMatter> customFrontMatter)
    {
        if (string.IsNullOrWhiteSpace(path))
            throw new ArgumentException("a note needs a path", nameof(path));
        if (string.IsNullOrWhiteSpace(slug))
            throw new ArgumentException("a note needs a slug", nameof(slug));

        Id = id;
        Path = path;
        Title = title;
        Slug = slug;
        CreatedAt = createdAt;
        UpdatedAt = updatedAt;
        _tags = tags;
        _links = links;
        _customFrontMatter = customFrontMatter;
    }

    public static Note Create(NoteFile file)
    {
        var now = Timestamp();
        return new Note(Guid.CreateVersion7().ToString(), file.Path, file.Title, file.Slug, now, now, [], [], [])
            .Apply(file);
    }

    public static Note Rehydrate(
        string id,
        string path,
        string? title,
        string slug,
        string createdAt,
        string updatedAt,
        IEnumerable<string> tags,
        IEnumerable<Link> links,
        IEnumerable<CustomFrontMatter> customFrontMatter) =>
        new(id, path, title, slug, createdAt, updatedAt, [.. tags], [.. links], [.. customFrontMatter]);

    public void Reindex(NoteFile file)
    {
        if (string.IsNullOrWhiteSpace(file.Path))
            throw new ArgumentException("a note needs a path", nameof(file));
        if (string.IsNullOrWhiteSpace(file.Slug))
            throw new ArgumentException("a note needs a slug", nameof(file));

        Path = file.Path;
        Slug = file.Slug;
        Title = file.Title;
        UpdatedAt = Timestamp();
        Apply(file);
    }

    private Note Apply(NoteFile file)
    {
        _tags.Clear();
        _tags.AddRange(file.FrontMatter.Tags.Where(t => !string.IsNullOrWhiteSpace(t)).Distinct());

        _customFrontMatter.Clear();
        _customFrontMatter.AddRange(
            file.FrontMatter.Custom.Select(e => new CustomFrontMatter(e.Key, e.Value)).Distinct());

        _links.Clear();
        _links.AddRange(file.Links.Select(l => Link.Create(l.TargetSlug, l.Raw, l.Display, l.Line, l.Col)));
        return this;
    }

    private static string Timestamp() => DateTime.UtcNow.ToString("yyyy-MM-dd HH:mm:ss");
}
