namespace Dendrite.Infrastructure.Persistence.Rows;

internal sealed record NoteRow
{
    public string Id { get; }
    public string Path { get; }
    public string? Title { get; }
    public string Slug { get; }
    public string CreatedAt { get; }
    public string UpdatedAt { get; }

    public NoteRow(string id, string path, string? title, string slug, string createdAt, string updatedAt)
    {
        Id = id;
        Path = path;
        Title = title;
        Slug = slug;
        CreatedAt = createdAt;
        UpdatedAt = updatedAt;
    }
}
