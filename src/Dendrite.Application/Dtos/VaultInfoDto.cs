namespace Dendrite.Application;

public sealed record VaultInfoDto
{
    public string Name { get; }
    public string Path { get; }
    public long NoteCount { get; }

    public VaultInfoDto(string name, string path, long noteCount)
    {
        Name = name;
        Path = path;
        NoteCount = noteCount;
    }
}
