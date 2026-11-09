using Dendrite.Domain.Models;

namespace Dendrite.Application.Abstractions;

public interface INoteRepository
{
    void Insert(Note note);
    void Update(string noteId, string path, string title, string slug);
    void Delete(string noteId);
    Note? GetBySlug(string slug);
    Note? GetByPath(string path);
    int GetNoteCount();
    List<Note> GetByTag(string tag);
    List<Note> Search(string query);
}

public interface ITagRepository
{
    void Insert(IEnumerable<Tag> tags);
    void InsertNoteTags(IEnumerable<NoteTag> noteTags);
    List<Tag> GetByNames(IEnumerable<string> names);
    List<Tag> Search(string query);
    void DeleteNoteTags(string noteId);
}

public interface ILinkRepository
{
    void Insert(IEnumerable<Link> links);
    void Delete(string fromNoteId);
    List<BacklinkDto> GetBacklinks(string slug);
    List<Link> GetBrokenLinks(string fromNoteId);
}

public interface ICfeRepository
{
    void InsertRange(IEnumerable<CustomFrontMatter> entries);
    void Delete(string noteId);
    List<Note> Search(string key, string value);
}

public interface IIndexRepository
{
    void Wipe();
}
