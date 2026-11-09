namespace Dendrite.Domain.Models;

/// <summary>Loads and persists whole <see cref="Note"/> aggregates (note, tags, links and custom frontmatter together).</summary>
public interface INoteRepository
{
    Note? GetBySlug(string slug);

    Note? GetByPath(string path);

    void Add(Note note);

    /// <summary>Persists the current state of an existing aggregate, replacing its tags, links and custom frontmatter.</summary>
    void Save(Note note);

    void Remove(Note note);
}
