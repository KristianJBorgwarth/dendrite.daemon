namespace Dendrite.Domain.Models;

public interface INoteRepository
{
    Note? GetBySlug(string slug);

    Note? GetByPath(string path);

    void Add(Note note);

    void Save(Note note);

    void Remove(Note note);
}
