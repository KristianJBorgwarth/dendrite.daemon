namespace Dendrite.Application.Abstractions;

/// <summary>Reads the vault's files on disk. Resolution is relative to the initialized vault.</summary>
public interface IVaultFiles
{
    string ReadText(string path);

    string ReadTemplate(string templateName);

    /// <summary>Full paths of every indexable markdown file, honouring the vault's excludes.</summary>
    IReadOnlyList<string> ListNotePaths();
}
