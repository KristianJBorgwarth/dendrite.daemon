namespace Dendrite.Domain.Models;

#pragma warning disable IDE0290 // Use primary constructor
public sealed class NoteTag
{
    public string NoteId { get; }
    public string TagId { get; }

    public NoteTag(string noteId, string tagId)
    {
        NoteId = noteId;
        TagId = tagId;
    }
}
