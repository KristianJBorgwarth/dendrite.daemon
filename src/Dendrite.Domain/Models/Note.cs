namespace Dendrite.Domain.Models;

#pragma warning disable IDE0290 // Use primary constructor
public sealed class Note
{
    public string Id { get; }
    public string Path { get; }
    public string? Title { get; }
    public string Slug { get; }
    public string CreatedAt { get; }
    public string UpdatedAt { get; }

    public Note(string id, string path, string? title, string slug, string createdAt, string updatedAt)
    {
        Id = id;
        Path = path;
        Title = title;
        Slug = slug;
        CreatedAt = createdAt;
        UpdatedAt = updatedAt;
    }

    public static Note Create(string path, string title, string slug)
    {
        var now = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss");
        return new Note(Guid.CreateVersion7().ToString(), path, title, slug, now, now);
    }
}
