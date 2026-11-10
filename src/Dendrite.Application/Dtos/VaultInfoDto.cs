namespace Dendrite.Application;

public sealed record VaultInfoDto
{
    public string Name { get; }
    public string Path { get; }
    public long NoteCount { get; }
    public bool Active { get; }
    public bool Missing { get; }

    public VaultInfoDto(string name, string path, long noteCount, bool active, bool missing = false)
    {
        Name = name;
        Path = path;
        NoteCount = noteCount;
        Active = active;
        Missing = missing;
    }
}
