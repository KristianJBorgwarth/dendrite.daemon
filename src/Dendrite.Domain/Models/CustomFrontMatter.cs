namespace Dendrite.Domain.Models;

#pragma warning disable IDE0290 // Use primary constructor
public sealed class CustomFrontMatter
{
    public string NoteId { get; }
    public string Key { get; }
    public string Value { get; }

    public CustomFrontMatter(string noteId, string key, string value)
    {
        NoteId = noteId;
        Key = key;
        Value = value;
    }
}
