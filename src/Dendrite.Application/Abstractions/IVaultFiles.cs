namespace Dendrite.Application.Abstractions;

public interface IVaultFiles
{
    string ReadText(string path);

    string ReadTemplate(string templateName);

    IReadOnlyList<string> ListNotePaths();
}
